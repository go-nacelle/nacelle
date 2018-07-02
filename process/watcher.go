package process

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/efritz/glock"
	"github.com/efritz/nacelle/logging"
)

// processWatcher coordinates goroutines to detect when an
// application should begin a shutdown or abort process. A
// graceful exit will happen on any of the following:
//   - the errChan closing
//   - the halt channel closing
//   - receiving a SIGINT or SIGKILL
//   - a non-nil error from a process
//   - a nil error from a process without silent exit
//
// An immediate abort will happen on any fo the following:
//   - receiving a second SIGINT or SIGKILL
//   - the timeout duration elapsing after shutdown begins
//
// The watcher will close the outChan after the errChan closes
// or after aborting.
type processWatcher struct {
	errChan         <-chan errMeta
	outChan         chan<- error
	halt            <-chan struct{}
	logger          logging.Logger
	clock           glock.Clock
	shutdownTimeout time.Duration
	abortSignal     chan struct{}
	shutdownSignal  chan struct{}
	done            chan struct{}
	abortOnce       *sync.Once
	shutdownOnce    *sync.Once
}

var shutdownSignals = []syscall.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
}

func newWatcher(
	errChan <-chan errMeta,
	outChan chan<- error,
	halt <-chan struct{},
	configs ...watcherConfigFunc,
) *processWatcher {
	watcher := &processWatcher{
		errChan:        errChan,
		outChan:        outChan,
		halt:           halt,
		logger:         logging.NewNilLogger(),
		clock:          glock.NewRealClock(),
		abortSignal:    make(chan struct{}),
		shutdownSignal: make(chan struct{}),
		done:           make(chan struct{}),
		abortOnce:      &sync.Once{},
		shutdownOnce:   &sync.Once{},
	}

	for _, f := range configs {
		f(watcher)
	}

	return watcher
}

// watch begins executing goroutines to watch for the shutdown
// and abort conditions described above.
func (w *processWatcher) watch() {
	go w.watchErrors()
	go w.watchSignal()
	go w.watchHaltChan()
	go w.watchShutdownTimeout()
}

// wait will unblock once the output channel has closed.
func (w *processWatcher) wait() {
	<-w.done
}

//
//

func (w *processWatcher) watchErrors() {
	defer close(w.done)
	defer close(w.outChan)
	defer w.shutdown()

	for {
		select {
		case <-w.abortSignal:
			// Just unblock the caller immediately
			return

		case err, ok := <-w.errChan:
			if !ok {
				// All processes have exited cleanly, so there
				// isn't anything left for us to do here. Yay!
				return
			}

			if err.err == nil {
				// Always log this
				w.logger.Info("%s has stopped cleanly", err.source.Name())

				// If we got a nil error but the process was marked
				// as something not necessarily long-running, stop
				// processing this value. Otherwise, let it fall all
				// the way through to the shutdown call below.

				if err.silentExit {
					continue
				}
			} else {
				// Always log this
				w.logger.Error("%s returned a fatal error (%s)", err.source.Name(), err.err.Error())

				// Inform the client of the watcher of this fatal error
				w.outChan <- err.err
			}

			w.shutdown()
		}
	}
}

func (w *processWatcher) watchSignal() {
	var (
		urgent  = false
		signals = make(chan os.Signal, 1)
	)

	for _, s := range shutdownSignals {
		signal.Notify(signals, s)
	}

	for {
		select {
		case <-w.done:
			return
		case <-signals:
		}

		if urgent {
			// Second signal, begin abort
			w.logger.Error("Received second signal, no longer waiting for graceful exit")
			w.abort()
			return
		}

		// First signal, begin shutdown
		w.logger.Info("Received signal")
		urgent = true
		w.shutdown()
	}
}

func (w *processWatcher) watchHaltChan() {
	select {
	case <-w.done:
		return
	case <-w.halt:
		// Wait for an external signal
	}

	w.logger.Info("Received external shutdown request")
	w.shutdown()
}

func (w *processWatcher) watchShutdownTimeout() {
	if w.shutdownTimeout == 0 {
		return
	}

	select {
	case <-w.done:
		return
	case <-w.shutdownSignal:
		// Wait for shutdown before starting the timer
	}

	select {
	case <-w.done:
	case <-w.clock.After(w.shutdownTimeout):
		// Shutdown has elapsed the shutdown timeout,
		// begin a forceful abort so the app isn't blocked
		// on shutdown indefinitely.
		w.abort()
	}
}

func (w *processWatcher) shutdown() {
	w.shutdownOnce.Do(func() {
		close(w.shutdownSignal)
		w.logger.Info("Starting graceful shutdown")
	})
}

func (w *processWatcher) abort() {
	w.abortOnce.Do(func() {
		close(w.abortSignal)
	})
}
