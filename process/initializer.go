package process

import "github.com/efritz/nacelle/config"

type (
	// Initializer is an interface that is called once on app
	// startup.
	Initializer interface {
		// Init reads the given configuration and prepares
		// something for use by a process. This can be loading
		// files from disk, connecting to a remote service,
		// initializing shared data structures, and inserting
		// a service into a shared service container.
		Init(config config.Config) error
	}

	// Finalizer is an optional extension of an Initializer that
	// supports finalization. This is useful for initializers
	// that need to tear down a background process before the
	// process exits, but needs to be started early in the boot
	// process (such as flushing logs or metrics).
	Finalizer interface {
		// Finalize is called after the application has stopped
		// all running processes.
		Finalize() error
	}

	// InitializerFunc is a non-struct version of an initializer.
	InitializerFunc func(config config.Config) error
)

// Init calls the underlying function.
func (f InitializerFunc) Init(config config.Config) error {
	return f(config)
}
