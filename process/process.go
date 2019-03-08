package process

type (
	// Process is an interface that continually performs a behavior
	// during the life of a program. Generally, one process should
	// do a single thing. Multiple processes can be registered to
	// a process container and those processes can coordinate and
	// communicate through shared services.
	Process interface {
		Initializer

		// Start begins performing the core action of the process.
		// For servers, this will begin accepting clients ona  port.
		// For workers, this may begin reading from a remote work
		// queue and processing messages. This method should block
		// until a fatal error occurs, or until the Stop method is
		// called (at which point a nil error should be returned).
		// If this method is non-blocking, then the process should
		// be registered with the WithSilentExit option.
		Start() error

		// Stop informs the work being performed by the Start
		// method to begin a graceful shutdown. This method is
		// not expected to block until shutdown completes.
		Stop() error
	}
)
