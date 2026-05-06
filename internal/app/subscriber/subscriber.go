package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/kafka"
)

const (
	consumerGroupNameAnalyticsEvents = "todolist_analytics-events"
)

type (
	Messages interface {
		Updates() <-chan schemes.UpdateInterface
	}

	Bot interface {
		ProcessUpdate(context.Context, schemes.UpdateInterface)
	}
)

type Subscriber struct {
	logger *zap.Logger
	config Config

	botMessages Messages

	cgAnalyticEvents *kafka.ConsumerGroup

	bot Bot

	quit chan struct{}
}

func New(
	config Config,
	logger *zap.Logger,
	kf *kafka.Kafka,
	botMessages Messages,
	bot Bot,
) (*Subscriber, error) {

	var cgAnalyticEvents *kafka.ConsumerGroup
	if config.AnalyticEvents.Enabled {
		var err error
		cgAnalyticEvents, err = kafka.NewConsumerGroup(kf, consumerGroupNameAnalyticsEvents)
		if err != nil {
			return nil, fmt.Errorf("new consumer: %w", err)
		}
	}

	s := &Subscriber{
		logger:      logger,
		config:      config,
		botMessages: botMessages,
		bot:         bot,

		cgAnalyticEvents: cgAnalyticEvents,

		quit: make(chan struct{}),
	}

	if s.config.AnalyticEvents.Enabled {
		s.addAnalyticInternalEventsConsumer()
	}

	return s, nil
}

func (s *Subscriber) Start(_ context.Context) error {
	go s.messagesConsumer()
	if s.config.AnalyticEvents.Enabled {
		s.cgAnalyticEvents.Consume(s.config.AnalyticEvents.PollRecords)
	}

	return nil
}

func (s *Subscriber) Stop(_ context.Context) error {
	close(s.quit)

	if s.config.AnalyticEvents.Enabled {
		s.cgAnalyticEvents.Close()
	}

	return nil
}

func (s *Subscriber) messagesConsumer() {
	ch := s.botMessages.Updates()
	for {
		select {
		case <-s.quit:
			return
		case upd, ok := <-ch:
			if !ok {
				return // channel closed
			}
			// крутим всё горутинах, чтобы не блокироваться
			go func() {
				ctx, span := otel.Tracer("subscriber").Start(context.Background(), "bot-message")
				defer span.End()

				defer func(t time.Time) {
					metricRequestsTotal.WithLabelValues("bot-message").Inc()
					metricRequestsDuration.WithLabelValues("bot-message").Observe(time.Since(t).Seconds())
				}(time.Now())

				s.bot.ProcessUpdate(ctx, upd)
			}()
		}
	}

}

func (s *Subscriber) addAnalyticInternalEventsConsumer() {
	callback := func(_ context.Context, payload kafka.Payload) error {
		s.logger.Info("new analytic event", zap.String("payload", string(payload)))

		return nil
	}

	for _, topic := range s.config.AnalyticEvents.Topics {
		s.cgAnalyticEvents.Subscribe(topic, callback)
	}
}
