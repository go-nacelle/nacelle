package nacelle

type (
	// Process is a monitored worker. It is meant for long-running
	// parts of a program (servers and background workers) which can
	// be initialized at program startup and be left to themselves.
	// See the process package for more concrete examples.
	Process interface {
		// Init configures the process so that it can readily begin
		// work. This may stash config values, open ports, or attempt
		// to connect to a remote service.
		Init(config Config) error

		// Start begins doing work. This method should block until the
		// useful part of the process has completed. This method should
		// be written so that a call to Stop interrupts the work in this
		// method in a meaningful way. On success, a nil error should be
		// returned.
		Start() error

		// Stop should interrupt the routine running the Start method.
		// Generally this is done by closing a listener or a channel which
		// the frequently read by the Start method.
		Stop() error
	}

	// Initializer is the init-only portion of a Process. This is meant
	// to do things like setting up global services (e.g. remote connections)
	// which can be used by processes.
	Initializer interface {
		Init(config Config) error
	}

	// InitializerFunc is a function which implements Initializer.
	InitializerFunc func(config Config) error
)

// Init calls the underlying InitializerFunc.
func (f InitializerFunc) Init(config Config) error {
	return f(config)
}
