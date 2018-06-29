package process

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/efritz/nacelle/logging"
)

type (
	processWatcher struct {
		shutdownCallback func()
		logger           logging.Logger
		shutdownTimeout  time.Duration
		startErrors      <-chan errMeta
		errChan          chan<- error
		halt             <-chan struct{}
		directives       chan watcherDirective
		draining         chan struct{}
	}

	watcherDirective int
)

const (
	directiveAbort watcherDirective = iota
	directiveShutdown
)

func newWatcher(
	shutdownCallback func(),
	logger logging.Logger,
	shutdownTimeout time.Duration,
	startErrors <-chan errMeta,
	errChan chan<- error,
	halt <-chan struct{},
) *processWatcher {
	return &processWatcher{
		shutdownCallback: shutdownCallback,
		logger:           logger,
		shutdownTimeout:  shutdownTimeout,
		startErrors:      startErrors,
		errChan:          errChan,
		halt:             halt,
		directives:       make(chan watcherDirective),
		draining:         make(chan struct{}),
	}
}

func (w *processWatcher) watch() {
	go w.watchSignal()
	go w.watchErrors()
	go w.watchHaltChan()
	go w.watchShutdownTimeout()

	for directive := range w.directives {
		if directive == directiveAbort {
			return
		}

		select {
		case <-w.draining:
			continue
		default:
		}

		close(w.draining)
		w.logger.Info("Starting graceful shutdown")
		w.shutdownCallback()
	}
}

func (w *processWatcher) watchSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	urgent := false
	for range sigChan {
		if urgent {
			w.logger.Error("Received second signal, no longer waiting for graceful exit")
			w.abort()
			return
		}

		urgent = true
		w.logger.Info("Received signal, starting graceful shutdown")
		w.shutdown()
	}
}

func (w *processWatcher) watchErrors() {
	defer close(w.errChan)

	for err := range w.startErrors {
		if err.err == nil {
			w.logger.Info("%s has stopped cleanly", err.process.Name())

			if err.process.silentExit {
				continue
			}
		} else {
			w.logger.Error("%s returned a fatal error (%s)", err.process.Name(), err.err.Error())
			w.errChan <- err.err
		}

		w.shutdown()
	}

	w.shutdown()
}

func (w *processWatcher) watchHaltChan() {
	<-w.halt
	w.logger.Info("Received external shutdown request")
	w.shutdown()
}

func (w *processWatcher) watchShutdownTimeout() {
	<-w.draining
	<-time.After(w.shutdownTimeout)
	w.abort()
}

func (w *processWatcher) abort() {
	w.directives <- directiveAbort
}

func (w *processWatcher) shutdown() {
	w.directives <- directiveShutdown
}
