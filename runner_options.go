package nacelle

import "time"

type (
	// ProcessRunnerConfigFunc is a function used to configure an instance
	// of a ProcessRunner.
	ProcessRunnerConfigFunc func(*ProcessRunner)
)

// WithShutdownTimeout sets the maximum time it will wait for a process to
// exit during a graceful shutdown.
func WithShutdownTimeout(timeout time.Duration) ProcessRunnerConfigFunc {
	return func(c *ProcessRunner) { c.shutdownTimeout = timeout }
}
