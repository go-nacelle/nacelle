package process

import (
	"github.com/efritz/nacelle/config"
)

type (
	// Initializer is an object that is called once on app
	// startup.
	Initializer interface {
		// Init reads the given configuration and prepares
		// something for use by a process. This can be loading
		// files from disk, connecting to a remote service,
		// initializing shared data structures, and inserting
		// a service into a shared service container.
		Init(config config.Config) error
	}

	// InitializerFunc is a non-struct version of an initializer.
	InitializerFunc func(config config.Config) error
)

// Init calls the underlying function.
func (f InitializerFunc) Init(config config.Config) error {
	return f(config)
}
