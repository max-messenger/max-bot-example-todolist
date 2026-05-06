package kafka

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Module = fx.Module(
		"kafka",
		fx.Provide(
			NewPool,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("kafka")
		}),
	)
)
