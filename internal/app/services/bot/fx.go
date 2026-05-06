package bot

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/clients/maxbot"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/repository/staterepo"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/repository/todorepo"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/analytic"
)

var Module = fx.Module(
	"bot",
	fx.Provide(
		New,
		Adapter,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("bot")
	}),
)

type (
	AdapterIn struct {
		fx.In

		Client   *maxbot.Client
		Store    *staterepo.Repository
		Repo     *todorepo.Repository
		Analytic *analytic.Service
	}

	AdapterOut struct {
		fx.Out

		Client   Client
		Store    Store
		Repo     Repository
		Analytic Analytic
	}
)

func Adapter(in AdapterIn) AdapterOut {
	return AdapterOut{
		Client:   in.Client,
		Store:    in.Store,
		Repo:     in.Repo,
		Analytic: in.Analytic,
	}
}
