//nolint:revive
package maxbot

import "time"

type Config struct {
	TestApiUrl       string        `yaml:"test_api_url"`
	Token            string        `yaml:"token"`
	Url              string        `yaml:"url"`
	Path             string        `yaml:"path"`
	Secret           string        `yaml:"secret"` // nolint:gosec
	SubscriptionType string        `yaml:"subscription_type"`
	Timeout          time.Duration `yaml:"timeout"`
}

func (c Config) isTest() bool {
	return c.TestApiUrl != ""
}
