package process

import (
	"time"

	"github.com/efritz/glock"
	"github.com/efritz/nacelle/logging"
)

// RunnerConfigFunc is a function used to configure an instance
// of a ProcessRunner.
type RunnerConfigFunc func(*runner)

// WithLogger sets the logger used by the runner.
func WithLogger(logger logging.Logger) RunnerConfigFunc {
	return func(c *runner) { c.logger = logger }
}

//WithClock sets the clock used by the runner.
func WithClock(clock glock.Clock) RunnerConfigFunc {
	return func(c *runner) { c.clock = clock }
}

// WithShutdownTimeout sets the maximum time it will wait for a process to
// exit during a graceful shutdown.
func WithShutdownTimeout(timeout time.Duration) RunnerConfigFunc {
	return func(c *runner) { c.shutdownTimeout = timeout }
}
