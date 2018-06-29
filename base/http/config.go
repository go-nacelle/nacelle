package http

import (
	"fmt"
	"time"
)

type (
	Config struct {
		HTTPPort           int    `env:"http_port" default:"5000"`
		HTTPCertFile       string `env:"http_cert_file"`
		HTTPKeyFile        string `env:"http_key_file"`
		RawShutdownTimeout int    `env:"http_shutdown_timeout" default:"5"`

		ShutdownTimeout time.Duration
	}

	configToken string
)

var (
	ConfigToken      = NewConfigToken("default")
	ErrBadCertConfig = fmt.Errorf("cert file and key file must both be supplied or both be omitted")
)

func NewConfigToken(name string) interface{} {
	return configToken(fmt.Sprintf("nacelle-base-http-%s", name))
}

func (c *Config) PostLoad() error {
	if (c.HTTPCertFile == "") != (c.HTTPKeyFile == "") {
		return ErrBadCertConfig
	}

	c.ShutdownTimeout = time.Duration(c.RawShutdownTimeout) * time.Second
	return nil
}
