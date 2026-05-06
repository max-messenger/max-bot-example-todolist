package app

import (
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist"
)

func CreateApp(cfgName string) *fx.App {
	box := todolist.NewBox(
		Modules,
		todolist.WithAppName("todolist"),
		todolist.WithConfigFile(cfgName),
	)

	return box.CreateApp()
}
