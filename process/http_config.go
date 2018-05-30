package process

import (
	"errors"
	"fmt"
	"time"

	"github.com/efritz/nacelle"
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

// RegisterHTTPConfigs adds the required configs for an HTTP server to the given map. If any tag
// modifiers are supplied, they are run over each of the required configs (this may require
// some knowledge about package internals).
func RegisterHTTPConfigs(m map[interface{}]interface{}, modifiers ...nacelle.TagModifier) error {
	c, err := nacelle.ApplyTagModifiers(&HTTPConfig{}, modifiers...)
	if err != nil {
		return err
	}

	m[HTTPConfigToken] = c
	return nil
}
