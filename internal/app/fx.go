package app

import (
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist/docs"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/clients"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/repository"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/router"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/router/todolistctrl"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/router/webhookctrl"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/subscriber"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/kafka"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/postgres"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/rediscli"
	"github.com/max-messenger/max-bot-example-todolist/pkg/migrate"
)

// Modules application modules.
var Modules = fx.Options(
	// controllers
	todolistctrl.Module,
	webhookctrl.Module,

	// subscriber
	subscriber.Module,

	// clients
	clients.Module,

	// repositories
	repository.Module,

	// connections
	postgres.Module,
	rediscli.Module,
	kafka.Module,

	// migrations
	migrate.Module,

	// services
	services.Module,

	router.Module,
	docs.Module,

	fx.Provide(NewConfig),
)
