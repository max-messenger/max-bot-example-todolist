package rediscli

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"rediscli",
	fx.Provide(
		NewPool,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("rediscli")
	}),
	fx.Invoke(
		func(lc fx.Lifecycle, pp *Pool) {
			lc.Append(fx.Hook{
				OnStop: pp.Stop,
			})
		},
	),
)
