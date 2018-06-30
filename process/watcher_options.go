package process

import (
	"time"

	"github.com/efritz/glock"
	"github.com/efritz/nacelle/logging"
)

type watcherConfigFunc func(*processWatcher)

func withWatcherLogger(logger logging.Logger) watcherConfigFunc {
	return func(w *processWatcher) { w.logger = logger }
}

func withWatcherClock(clock glock.Clock) watcherConfigFunc {
	return func(w *processWatcher) { w.clock = clock }
}

func withWatcherShutdownTimeout(timeout time.Duration) watcherConfigFunc {
	return func(w *processWatcher) { w.shutdownTimeout = timeout }
}
