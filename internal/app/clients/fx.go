package clients

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/clients/maxbot"
)

var Module = fx.Module(
	"clients",
	fx.Options(
		maxbot.Module,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("clients")
	}),
)
