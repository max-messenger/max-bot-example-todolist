package app

import (
	"go.uber.org/config"
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist/internal/app/clients/maxbot"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/services/analytic"
	"github.com/max-messenger/max-bot-example-todolist/internal/app/subscriber"
)

type (
	// Config is custom app configuration.
	Config struct {
		fx.Out

		MaxBot     maxbot.Config     `yaml:"max_bot"`
		Analytic   analytic.Config   `yaml:"analytic"`
		Subscriber subscriber.Config `yaml:"subscriber"`
	}

	// NewConfigParams contains raw config.
	ConfigProvider struct {
		fx.In

		Provider config.Provider `name:"config_provider"`
	}
)

// NewConfig return new config instance.
func NewConfig(cp ConfigProvider) (Config, error) {
	maxbotCfg := maxbot.Config{}

	err := cp.Provider.Get("max_bot").Populate(&maxbotCfg)
	if err != nil {
		return Config{}, err
	}

	analyticsCfg := analytic.Config{}
	err = cp.Provider.Get("analytic").Populate(&analyticsCfg)
	if err != nil {
		return Config{}, err
	}

	subscriberCfg := subscriber.NewConfig()
	err = cp.Provider.Get("subscriber").Populate(&subscriberCfg)
	if err != nil {
		return Config{}, err
	}

	return Config{
		MaxBot:     maxbotCfg,
		Analytic:   analyticsCfg,
		Subscriber: subscriberCfg,
	}, nil
}
