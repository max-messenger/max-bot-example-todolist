package config

import (
	"os"

	"go.uber.org/config"
	"go.uber.org/fx"

	"github.com/max-messenger/max-bot-example-todolist/pkg/bgtasker"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/grpccli"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/kafka"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/postgres"
	"github.com/max-messenger/max-bot-example-todolist/pkg/connections/rediscli"
	"github.com/max-messenger/max-bot-example-todolist/pkg/health"
	"github.com/max-messenger/max-bot-example-todolist/pkg/info"
	"github.com/max-messenger/max-bot-example-todolist/pkg/logger"
	"github.com/max-messenger/max-bot-example-todolist/pkg/migrate"
	"github.com/max-messenger/max-bot-example-todolist/pkg/ratelimiter"
	"github.com/max-messenger/max-bot-example-todolist/pkg/server"
	"github.com/max-messenger/max-bot-example-todolist/pkg/telemetry"
)

type Config struct {
	fx.Out

	App       AppConfig
	Logger    logger.Config
	Health    health.Config
	Info      info.Config
	Telemetry telemetry.Config

	SystemServer server.SystemConfig
	HTTPServer   *server.HTTPConfig
	GRPCServer   *server.GRPCConfig

	RateLimiter ratelimiter.Config

	Background *bgtasker.Config

	// Connncections
	RedisCli map[string]*rediscli.Config
	Postgres map[string]*postgres.Config
	GRPCClis map[string]*grpccli.ConnConfig
	Kafka    map[string]*kafka.Config

	// migrations
	Migrate migrate.Config

	Provider config.Provider `name:"config_provider"`
}

type AppConfig struct {
	Hostname string `yaml:"hostname"`
}

type internalConfig struct {
	App       AppConfig        `yaml:"app"`
	Logger    logger.Config    `yaml:"logger"`
	Servers   server.Config    `yaml:"servers"`
	Telemetry telemetry.Config `yaml:"telemetry"`

	RateLimiter ratelimiter.Config `yaml:"rate_limiter"`

	Background *bgtasker.Config `yaml:"background"`

	// connections
	RedisCli map[string]*rediscli.Config    `yaml:"redis"`        //nolint:tagliatelle
	Postgres map[string]*postgres.Config    `yaml:"postgres"`     //nolint:tagliatelle
	GRPCClis map[string]*grpccli.ConnConfig `yaml:"grpc_clients"` //nolint:tagliatelle
	Kafka    map[string]*kafka.Config       `yaml:"kafka"`        //nolint:tagliatelle

	Migrate migrate.Config `yaml:"migrate"`

	Provider config.Provider
}

func NewConfig(appName, cfgFile string) func() (Config, error) {
	return func() (Config, error) {
		inCfg, err := newInternalConfig(cfgFile)
		if err != nil {
			return Config{}, err
		}

		infoCfg := info.Config{
			AppName:  appName,
			Hostname: inCfg.App.Hostname,
		}

		healthCfg := health.Config{
			Hostname: inCfg.App.Hostname,
			Version:  infoCfg.BuildVersion(),
		}

		inCfg.Telemetry.ServiceName = appName
		inCfg.Telemetry.Hostname = inCfg.App.Hostname
		inCfg.Telemetry.Version = infoCfg.BuildVersion()

		cfg := Config{
			App:    inCfg.App,
			Logger: inCfg.Logger,
			Health: healthCfg,
			Info:   infoCfg,

			Telemetry: inCfg.Telemetry,

			RateLimiter: inCfg.RateLimiter,

			Background: inCfg.Background,

			SystemServer: inCfg.Servers.System,
			HTTPServer:   inCfg.Servers.HTTP,
			GRPCServer:   inCfg.Servers.GRPC,

			// connections
			RedisCli: inCfg.RedisCli,
			Postgres: inCfg.Postgres,
			GRPCClis: inCfg.GRPCClis,
			Kafka:    inCfg.Kafka,

			Migrate: inCfg.Migrate,

			// provider
			Provider: inCfg.Provider,
		}

		if inCfg.Background != nil {
			cfg.Background = inCfg.Background.Prepare()
		} else {
			cfg.Background = bgtasker.NewConfig()
		}

		for connName, connCfg := range cfg.RedisCli {
			cfg.RedisCli[connName] = connCfg.Prepare()
		}

		for connName, connCfg := range cfg.Kafka {
			cfg.Kafka[connName] = connCfg.Prepare()
		}

		for connName, connCfg := range cfg.GRPCClis {
			cfg.GRPCClis[connName] = connCfg.Prepare()
		}

		return cfg, nil
	}
}

func newInternalConfig(fileName string) (internalConfig, error) {
	provider, err := config.NewYAML(
		config.Expand(os.LookupEnv),
		config.File(fileName),
		config.Permissive(),
	)

	if err != nil {
		return internalConfig{}, err
	}

	c := internalConfig{
		Logger:   logger.DefaultConfig(),
		Provider: provider,
	}

	err = provider.Get("").Populate(&c)
	if err != nil {
		return internalConfig{}, err
	}

	return c, nil
}
