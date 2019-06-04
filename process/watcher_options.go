package process

import (
	"time"

	"github.com/efritz/glock"
	"github.com/go-nacelle/nacelle/logging"
)

// watcherConfigFunc is a function used to configure an instance
// of a processWatcher.
type watcherConfigFunc func(*processWatcher)

// withWatcherLogger sets the logger used by the watcher.
func withWatcherLogger(logger logging.Logger) watcherConfigFunc {
	return func(w *processWatcher) { w.logger = logger }
}

// withWatcherClock sets the clock used by the watcher.
func withWatcherClock(clock glock.Clock) watcherConfigFunc {
	return func(w *processWatcher) { w.clock = clock }
}

// withWatcherShutdownTimeout sets the timeout duration between the first
// shutdown and the watcher exiting. If the watcher does not exit within
// this duration an abort is invoked.
func withWatcherShutdownTimeout(timeout time.Duration) watcherConfigFunc {
	return func(w *processWatcher) { w.shutdownTimeout = timeout }
}
