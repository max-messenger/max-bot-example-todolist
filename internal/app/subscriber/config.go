package subscriber

import "time"

const (
	defaultProcessTimeout = time.Second * 2
)

type Config struct {
	AnalyticEvents ConsumerConfig `yaml:"analytic_events"`
}

type ConsumerConfig struct {
	Enabled        bool          `yaml:"enabled"`
	ProcessTimeout time.Duration `yaml:"process_timeout"`
	PollRecords    int           `yaml:"poll_records"`
	Topics         []string      `yaml:"topics"`
}

func NewConfig() Config {
	return Config{
		AnalyticEvents: ConsumerConfig{
			ProcessTimeout: defaultProcessTimeout,
		},
	}
}
