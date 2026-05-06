package subscriber

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/clients/maxbot"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/bot"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/kafka"
	"github.com/max-messenger/max-bot-example-todolist/pkg/grace"
)

var Module = fx.Module(
	"subscriber",
	fx.Provide(
		New,
		Adapter,
		kafkaAdapter,
		graceAdapter,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("subscriber")
	}),
)

type AdapterIn struct {
	fx.In

	Messages *maxbot.Client
	Bot      *bot.MaxBot
}

type AdapterOut struct {
	fx.Out

	Messages Messages
	Bot      Bot
}

func Adapter(in AdapterIn) AdapterOut {
	return AdapterOut{
		Messages: in.Messages,
		Bot:      in.Bot,
	}
}

func kafkaAdapter(kPool *kafka.Pool) (*kafka.Kafka, error) {
	kf, err := kPool.GetPool("analytic")
	if err != nil {
		return nil, err
	}

	return kf, nil
}

type graceOut struct {
	fx.Out

	Service grace.Service `group:"grace"`
}

func graceAdapter(sub *Subscriber) graceOut {
	return graceOut{
		Service: sub,
	}
}
