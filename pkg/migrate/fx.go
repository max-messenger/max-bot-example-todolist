package migrate

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"migrate",
	fx.Provide(
		New,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("migrate")
	}),
	fx.Invoke(
		func(lc fx.Lifecycle, mg *Migrations) {
			lc.Append(fx.Hook{
				OnStart: mg.Start,
			})
		},
	),
)
