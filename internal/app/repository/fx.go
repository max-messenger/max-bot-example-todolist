package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/repository/staterepo"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/repository/todorepo"
)

var Module = fx.Module(
	"repository",
	fx.Options(
		staterepo.Module,
		todorepo.Module,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("repository")
	}),
)
