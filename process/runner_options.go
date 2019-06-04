package process

import (
	"time"

	"github.com/efritz/backoff"
	"github.com/efritz/glock"
	"github.com/go-nacelle/nacelle/logging"
)

// RunnerConfigFunc is a function used to configure an instance
// of a ProcessRunner.
type RunnerConfigFunc func(*runner)

// WithLogger sets the logger used by the runner.
func WithLogger(logger logging.Logger) RunnerConfigFunc {
	return func(r *runner) { r.logger = logger }
}

// WithClock sets the clock used by the runner.
func WithClock(clock glock.Clock) RunnerConfigFunc {
	return func(r *runner) { r.clock = clock }
}

// WithStartTimeout sets the time it will wait for a process to become
// healthy after startup.
func WithStartTimeout(timeout time.Duration) RunnerConfigFunc {
	return func(r *runner) { r.startupTimeout = timeout }
}

// WithHealthCheckBackoff sets the backoff to use when waiting for processes
// to become healthy after startup.
func WithHealthCheckBackoff(backoff backoff.Backoff) RunnerConfigFunc {
	return func(r *runner) { r.healthCheckBackoff = backoff }
}

// WithShutdownTimeout sets the maximum time it will wait for a process to
// exit during a graceful shutdown.
func WithShutdownTimeout(timeout time.Duration) RunnerConfigFunc {
	return func(r *runner) { r.shutdownTimeout = timeout }
}
