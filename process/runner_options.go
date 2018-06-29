package process

import "time"

type (
	// RunnerConfigFunc is a function used to configure an instance
	// of a ProcessRunner.
	RunnerConfigFunc func(*runner)
)

// WithShutdownTimeout sets the maximum time it will wait for a process to
// exit during a graceful shutdown.
func WithShutdownTimeout(timeout time.Duration) RunnerConfigFunc {
	return func(c *runner) { c.shutdownTimeout = timeout }
}
