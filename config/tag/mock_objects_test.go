package tag

type (
	BasicConfig struct {
		X string `env:"a" default:"q"`
		Y string
	}

	ParentConfig struct {
		BasicConfig
	}
)
