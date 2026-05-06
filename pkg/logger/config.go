package logger

type Config struct {
	// Level is max level for write, required, one of (lower or uppercase):
	//   debug,info, warn, error, dpanic, panic, fatal
	Level             string `yaml:"level"`
	Format            string `yaml:"format"`
	DisableStacktrace bool   `yaml:"disable_stacktrace"`
}

func DefaultConfig() Config {
	return Config{
		Level:             "info",
		Format:            "json",
		DisableStacktrace: false,
	}
}
