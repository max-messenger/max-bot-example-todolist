package services

import (
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/analytic"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/bot"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/todolist"
)

var Module = fx.Module(
	"services",
	fx.Options(
		analytic.Module,
		bot.Module,
		todolist.Module,
	),
)
