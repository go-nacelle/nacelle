package process

import (
	"errors"
	"fmt"
	"time"
)

type (
	HTTPConfig struct {
		HTTPPort           int    `env:"http_port" default:"5000"`
		HTTPCertFile       string `env:"http_cert_file"`
		HTTPKeyFile        string `env:"http_key_file"`
		RawShutdownTimeout int    `env:"http_shutdown_timeout" default:"5"`

		ShutdownTimeout time.Duration
	}

	httpConfigToken string
)

var (
	HTTPConfigToken  = MakeHTTPConfigToken("default")
	ErrBadCertConfig = errors.New("cert file and key file must both be supplied or both be omitted")
)

func MakeHTTPConfigToken(name string) interface{} {
	return httpConfigToken(fmt.Sprintf("nacelle-process-http-%s", name))
}

func (c *HTTPConfig) PostLoad() error {
	if (c.HTTPCertFile == "") != (c.HTTPKeyFile == "") {
		return ErrBadCertConfig
	}

	c.ShutdownTimeout = time.Duration(c.RawShutdownTimeout) * time.Second
	return nil
}
