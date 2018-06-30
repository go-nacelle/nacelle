package process

import (
	"github.com/efritz/nacelle/config"
)

type (
	Initializer interface {
		Init(config config.Config) error
	}

	InitializerFunc func(config config.Config) error
)

func (f InitializerFunc) Init(config config.Config) error {
	return f(config)
}
