package analytic

type Config struct {
	Enabled     bool           `yaml:"enabled"`
	Topic       string         `yaml:"topic"`
	ExtraFields map[string]any `yaml:"extra_fields"`
}
