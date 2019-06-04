package process

import (
	"github.com/efritz/glock"
	"github.com/go-nacelle/nacelle/logging"
	"github.com/go-nacelle/nacelle/service"
)

// ParallelInitializerConfigFunc is a function used to configure an instance
// of a ParallelInitializer.
type ParallelInitializerConfigFunc func(*ParallelInitializer)

// WithParallelInitializerLogger sets the logger used by the runner.
func WithParallelInitializerLogger(logger logging.Logger) ParallelInitializerConfigFunc {
	return func(pi *ParallelInitializer) { pi.Logger = logger }
}

// WithParallelInitializerContainer sets the service container used by the runner.
func WithParallelInitializerContainer(container service.Container) ParallelInitializerConfigFunc {
	return func(pi *ParallelInitializer) { pi.Services = container }
}

// WithParallelInitializerClock sets the clock used by the runner.
func WithParallelInitializerClock(clock glock.Clock) ParallelInitializerConfigFunc {
	return func(pi *ParallelInitializer) { pi.clock = clock }
}
