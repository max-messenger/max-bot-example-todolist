package ratelimiter

import (
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/rediscli"
)

var Module = fx.Module(
	"rate_limiter",
	fx.Provide(
		adapterLimitGet,
	),
)

type adapterLimitGetIn struct {
	fx.In

	LimitGet []CustomLimitFunc `group:"rate_limit_custom"`
}

func adapterLimitGet(in adapterLimitGetIn) CustomLimitFuncs {
	return in.LimitGet
}

var ModuleLocal = fx.Module(
	"local_rate_limiter",
	fx.Provide(
		NewLocalLimiter,
	),
)

var ModuleShared = fx.Module(
	"shared_rate_limiter",
	fx.Provide(
		NewSharedLimiter,
		sharedRateLimiterRedisAdapter,
	),
)

func sharedRateLimiterRedisAdapter(cfg Config, rediscliPool *rediscli.Pool) (RedisClient, error) {
	sharedCfg := cfg.Shared

	return rediscliPool.GetPool(sharedCfg.PoolName)
}
