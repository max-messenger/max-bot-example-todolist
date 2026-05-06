package domain

import (
	maxbotcli "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Message struct {
	UserID    int64
	ChatID    int64
	Text      string
	Keyboards []*maxbotcli.Keyboard
	Files     []*schemes.UploadedInfo
}

func NewMessage(u *schemes.User, text string) *Message {
	return &Message{
		UserID: u.UserId,
		Text:   text,
	}
}

func (m *Message) AddKeyboard(kbd *maxbotcli.Keyboard) *Message {
	m.Keyboards = append(m.Keyboards, kbd)

	return m
}

func (m *Message) AddFile(file *schemes.UploadedInfo) *Message {
	m.Files = append(m.Files, file)

	return m
}
