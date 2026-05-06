package maxbot

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/pkg/grace"
	"github.com/max-messenger/max-bot-example-todolist/pkg/http/client"
)

var Module = fx.Module(
	"maxbot",
	fx.Provide(
		New,
		adapter,
		graceAdapter,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("maxbot")
	}),
)

func adapter(log *zap.Logger) HTTPClient {
	opts := []client.Option{
		client.WithLogger(log.Named("maxbot")),
	}

	return client.New("maxbot", opts...)
}

type graceOut struct {
	fx.Out

	Service grace.Service `group:"grace"`
}

func graceAdapter(sub *Client) graceOut {
	return graceOut{
		Service: sub,
	}
}
