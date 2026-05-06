package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	maxbotcli "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/domain"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/analytic"
)

type msgHandler func(context.Context, schemes.MessageCreatedUpdate) (domain.Message, error)

const (
	commandHelp   = "/help"
	commandInfo   = "/info"
	commandCreate = "/create"
	commandList   = "/list"
	commandDelete = "/delete"

	stateCreate = "state_create"
	infoMessage = `/info - info
/create - new record
/list - all record
/delete - delete record`
)

type Repository interface {
	Add(ctx context.Context, d domain.Todo) (int64, error)
	List(ctx context.Context, userID int64) ([]domain.Todo, error)
	Delete(ctx context.Context, d domain.Todo) error
}

type Store interface {
	GetState(ctx context.Context, userID int64) (string, error)
	SetState(ctx context.Context, userID int64, value string) error
	ClearState(ctx context.Context, userID int64) error
}

type Client interface {
	Updates() <-chan schemes.UpdateInterface
	SendMessage(ctx context.Context, msg domain.Message) (string, error)
	EditMessage(ctx context.Context, mid string, msg domain.Message) error
	DeleteMessage(ctx context.Context, mid string) error
}

type Analytic interface {
	Send(ctx context.Context, userID int64, event string, payload analytic.EventPayload)
}

type MaxBot struct {
	log *zap.Logger

	client   Client
	store    Store
	repo     Repository
	analytic Analytic
}

func New(
	l *zap.Logger,
	c Client,
	s Store,
	r Repository,
	a Analytic,
) *MaxBot {
	return &MaxBot{
		log:      l,
		client:   c,
		store:    s,
		repo:     r,
		analytic: a,
	}
}

func (m *MaxBot) ProcessUpdate(ctx context.Context, upd schemes.UpdateInterface) {
	m.log.Debug(
		"process update",
		zap.Any("user", upd.GetUserID()),
		zap.Any("update_type", upd.GetUpdateType()),
	)

	m.analytic.Send(
		ctx,
		upd.GetUserID(),
		"new_event",
		analytic.EventPayload{
			"update_type": upd.GetUpdateType(),
		},
	)

	switch update := upd.(type) {
	case *schemes.BotStartedUpdate:
		m.botStartedUpdate(ctx, *update)
	case *schemes.MessageCreatedUpdate:
		m.messageCreatedUpdate(ctx, *update)
	case *schemes.MessageCallbackUpdate:
		m.messageCallbackUpdate(ctx, *update)
	default:
		m.unknownMessageHandler(ctx, update)
	}
}

func (m *MaxBot) messageCallbackUpdate(ctx context.Context, upd schemes.MessageCallbackUpdate) {
	m.log.Info("Message callback update")

	if upd.Callback.Payload == "" {
		return
	}

	pd := domain.CallBack{}
	err := pd.Value(upd.Callback.Payload)
	if err != nil {
		m.log.Error("Failed to unmarshal callback payload", zap.Error(err))

		return
	}

	h := m.callbackDefaultHandler
	if pd.Key == commandDelete {
		h = m.callbackDeleteHandler
	}

	err = h(ctx, upd)
	if err != nil {
		m.log.Error("Failed to handle callback", zap.Error(err))
	}
}

func (m *MaxBot) messageCreatedUpdate(ctx context.Context, upd schemes.MessageCreatedUpdate) {
	m.log.Info("Message created update")

	stage, err := m.getState(ctx, upd)
	if err != nil {
		m.log.Error("", zap.Error(err))

		return
	}

	var h msgHandler
	switch stage {
	case commandHelp, commandInfo:
		h = m.helpStateHandler
	case commandCreate:
		h = m.commandCreateHandler
	case commandList:
		h = m.commandListHandler
	case stateCreate:
		h = m.stateCreateHandler
	case commandDelete:
		h = m.commandDeleteHandler
	default:
		h = m.defaultStateHandler
	}

	msg, err := h(ctx, upd)
	if err != nil {
		m.log.Warn("messageCreatedUpdate", zap.Error(err))
	}

	mid, err := m.client.SendMessage(ctx, msg)
	if err != nil {
		m.log.Warn("messageCreatedUpdate error", zap.Error(err))

		return
	}

	m.log.Debug("send message",
		zap.String("mid", mid),
		zap.Int64("chat_id", upd.GetChatID()),
		zap.Int64("user_id", upd.GetUserID()),
	)
}

func (m *MaxBot) unknownMessageHandler(_ context.Context, update schemes.UpdateInterface) {
	m.log.Info("unknownMessageHandler", zap.Any("update", update))
}

func (m *MaxBot) botStartedUpdate(ctx context.Context, upd schemes.BotStartedUpdate) {
	msg := domain.Message{
		UserID:    upd.GetUserID(),
		ChatID:    upd.GetChatID(),
		Text:      infoMessage,
		Keyboards: nil,
		Files:     nil,
	}

	_, err := m.client.SendMessage(ctx, msg)
	if err != nil {
		m.log.Warn("botStartedUpdate error", zap.Error(err))

		return
	}
}

func (m *MaxBot) getState(ctx context.Context, upd schemes.MessageCreatedUpdate) (string, error) {
	cmd := upd.GetCommand()
	if cmd != schemes.CommandUndefined {
		return cmd, nil
	}

	state, err := m.store.GetState(ctx, upd.GetUserID())
	if err != nil {
		return state, fmt.Errorf("getState: %w", err)
	}

	return state, nil
}

func (m *MaxBot) defaultStateHandler(_ context.Context, upd schemes.MessageCreatedUpdate) (domain.Message, error) {
	// /create
	out := "команда : " + upd.GetCommand()
	msg := domain.Message{
		UserID:    upd.GetUserID(),
		ChatID:    upd.GetChatID(),
		Text:      out,
		Keyboards: nil,
		Files:     nil,
	}

	return msg, nil
}

func (m *MaxBot) helpStateHandler(_ context.Context, upd schemes.MessageCreatedUpdate) (domain.Message, error) {
	msg := domain.Message{
		UserID:    upd.GetUserID(),
		ChatID:    upd.GetChatID(),
		Text:      infoMessage,
		Keyboards: nil,
		Files:     nil,
	}

	return msg, nil
}

func (m *MaxBot) commandCreateHandler(ctx context.Context, upd schemes.MessageCreatedUpdate) (domain.Message, error) {
	msg := domain.Message{
		UserID: upd.GetUserID(),
		ChatID: upd.GetChatID(),
		Text:   "Введите текст новой записи",
	}

	err := m.store.SetState(ctx, upd.GetUserID(), stateCreate)
	if err != nil {
		return msg, fmt.Errorf("setState: %w", err)
	}

	return msg, nil
}

func (m *MaxBot) commandListHandler(ctx context.Context, upd schemes.MessageCreatedUpdate) (domain.Message, error) {
	msg := domain.Message{
		UserID: upd.GetUserID(),
		ChatID: upd.GetChatID(),
		Text:   "У вас нет записей",
	}

	records, err := m.repo.List(ctx, msg.UserID)
	if err != nil {
		return msg, fmt.Errorf("list: %w", err)
	}

	fields := make([]string, 0, len(records))
	for _, record := range records {
		fields = append(fields, fmt.Sprintf("%s %s", record.Created.Format(time.DateTime), record.Message))
	}

	if len(fields) > 0 {
		msg.Text = strings.Join(fields, "\n")
	}

	return msg, nil
}

func (m *MaxBot) commandDeleteHandler(ctx context.Context, upd schemes.MessageCreatedUpdate) (domain.Message, error) {
	msg := domain.Message{
		UserID:    upd.GetUserID(),
		ChatID:    upd.GetChatID(),
		Keyboards: make([]*maxbotcli.Keyboard, 0),
		Text:      "У вас нет записей",
	}

	records, err := m.repo.List(ctx, msg.UserID)
	if err != nil {
		return msg, fmt.Errorf("list: %w", err)
	}

	if len(records) == 0 {
		return msg, nil
	}

	keyboard := &maxbotcli.Keyboard{}
	for _, record := range records {
		txt := fmt.Sprintf("%s %s", record.Created.Format(time.DateTime), record.Message)
		pd := domain.CallBack{
			InternalID: record.ID,
			Key:        commandDelete,
		}
		keyboard.AddRow().AddCallback(cutLimit(txt, 130), schemes.POSITIVE, pd.String())
	}

	msg.Keyboards = append(msg.Keyboards, keyboard)
	msg.Text = "Выберите запись для удаления:"

	return msg, nil
}

func (m *MaxBot) stateCreateHandler(ctx context.Context, upd schemes.MessageCreatedUpdate) (domain.Message, error) {
	msg := domain.Message{
		UserID: upd.GetUserID(),
		Text:   "Запись создана",
	}

	_, err := m.repo.Add(ctx, domain.Todo{
		UserID:  upd.GetUserID(),
		Message: upd.Message.Body.Text,
	})
	if err != nil {
		return msg, fmt.Errorf("add: %w", err)
	}

	return msg, nil
}

func (m *MaxBot) callbackDeleteHandler(ctx context.Context, upd schemes.MessageCallbackUpdate) error {
	pd := domain.CallBack{}
	err := pd.Value(upd.Callback.Payload)
	if err != nil {
		return fmt.Errorf("callback decode: %w", err)
	}

	td := domain.Todo{
		ID:     pd.InternalID,
		UserID: upd.GetUserID(),
	}
	err = m.repo.Delete(ctx, td)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	msg := domain.Message{
		UserID:    upd.GetUserID(),
		ChatID:    upd.GetChatID(),
		Keyboards: make([]*maxbotcli.Keyboard, 0),
	}

	records, err := m.repo.List(ctx, msg.UserID)
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}

	keyboard := &maxbotcli.Keyboard{}
	for _, record := range records {
		txt := fmt.Sprintf("%s %s", record.Created.Format(time.DateTime), record.Message)
		pd = domain.CallBack{
			InternalID: record.ID,
			Key:        commandDelete,
		}
		keyboard.AddRow().AddCallback(txt, schemes.POSITIVE, pd.String())
	}

	if len(records) == 0 {
		return m.client.DeleteMessage(ctx, upd.Message.Body.Mid)
	}

	msg.Keyboards = append(msg.Keyboards, keyboard)
	msg.Text = "Выберите запись для удаления:"

	err = m.client.EditMessage(ctx, upd.Message.Body.Mid, msg)
	if err != nil {
		return fmt.Errorf("editMessage: %w", err)
	}

	return nil
}

func (m *MaxBot) callbackDefaultHandler(_ context.Context, _ schemes.MessageCallbackUpdate) error {
	return nil
}

func cutLimit(s string, limit int) string {
	if len(s) <= limit {
		return s
	}

	return s[:limit]
}
