package analytic

import (
	"context"
	"maps"

	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/kafka"
	"github.com/max-messenger/max-bot-example-todolist/pkg/marshaler"
)

type EventPayload map[string]any

type (
	Producer interface {
		SendAsync(ctx context.Context, topic string, payload kafka.Payload)
	}
)

type Service struct {
	logger   *zap.Logger
	config   Config
	producer Producer
}

func New(
	logger *zap.Logger,
	config Config,
	producer Producer,
) *Service {
	return &Service{
		logger:   logger,
		config:   config,
		producer: producer,
	}
}

func (s *Service) Send(
	ctx context.Context,
	userID int64,
	event string,
	payload EventPayload,
) {
	if !s.config.Enabled {
		return
	}
	data := make(map[string]any, len(payload)+len(s.config.ExtraFields)+2)
	maps.Copy(data, payload)
	maps.Copy(data, s.config.ExtraFields)
	data["user_id"] = userID
	data["event"] = event

	s.logger.Debug("analytic send", zap.Any("data", data))

	rawData, err := marshaler.MarshalJSON(data)
	if err != nil {
		s.logger.Warn("analytic send", zap.Error(err))

		return
	}

	s.producer.SendAsync(
		context.WithoutCancel(ctx),
		s.config.Topic,
		rawData,
	)
}
