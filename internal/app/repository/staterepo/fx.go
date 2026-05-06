package staterepo

import (
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/rediscli"
)

var Module = fx.Module(
	"state",
	fx.Provide(
		NewRepository,
		Adapter,
	),
)

type (
	AdapterIn struct {
		fx.In

		Store *rediscli.Pool
	}

	AdapterOut struct {
		fx.Out

		Store Store
	}
)

func Adapter(in AdapterIn) (AdapterOut, error) {
	client, err := in.Store.GetPool("main")
	if err != nil {
		return AdapterOut{}, err
	}

	return AdapterOut{
		Store: client,
	}, nil
}
