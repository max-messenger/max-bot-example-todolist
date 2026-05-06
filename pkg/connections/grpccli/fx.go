package grpccli

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"grpc-conn",
	fx.Provide(
		NewPool,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("grpc-conn")
	}),
	fx.Invoke(func(lc fx.Lifecycle, pool *Pool) {
		lc.Append(fx.Hook{
			OnStop: pool.Stop,
		})
	}),
)
