//nolint:ireturn

package postgres

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"postgres_pool",
	fx.Provide(
		NewPostgresPool,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("postgres_pool")
	}),
	fx.Invoke(
		func(lc fx.Lifecycle, pp *Pool) {
			lc.Append(fx.Hook{
				OnStop: pp.Stop,
			})
		},
	),
)

type PoolParams struct {
	fx.In

	Cfgs map[string]*Config

	Logger *zap.Logger

	BeforeConnectCallbacks []BeforeConnectCallback `group:"pg_callbacks"`
	AfterConnectCallbacks  []AfterConnectCallback  `group:"pg_callbacks"`
}
