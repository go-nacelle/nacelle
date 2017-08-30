package http

type (
	Config struct {
		HTTPPort int `env:"HTTP_PORT" default:"5000"`
	}

	configToken struct{}
)

var ConfigToken = configToken{}
