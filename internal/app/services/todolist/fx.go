package todolist

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/repository/todorepo"
)

const serviceName = "todolist"

var Module = fx.Module(
	serviceName,
	fx.Provide(
		NewService,
		Adapter,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named(serviceName)
	}),
)

type (
	AdapterIn struct {
		fx.In

		MaterialsRepository *todorepo.Repository
	}

	AdapterOut struct {
		fx.Out

		MaterialsRepo Repository
	}
)

func Adapter(in AdapterIn) AdapterOut {
	return AdapterOut{

		MaterialsRepo: in.MaterialsRepository,
	}
}
