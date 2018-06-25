package nacelle

import (
	"os"
	"os/signal"
	"syscall"
)

type (
	processWatcher struct {
		shutdownCallback func()
		logger           Logger
		startErrors      <-chan errMeta
		errChan          chan<- error
		halt             <-chan struct{}
		directives       chan watcherDirective
	}

	watcherDirective int
)

const (
	directiveAbort watcherDirective = iota
	directiveShutdown
)

func newWatcher(
	shutdownCallback func(),
	logger Logger,
	startErrors <-chan errMeta,
	errChan chan<- error,
	halt <-chan struct{},
) *processWatcher {
	return &processWatcher{
		shutdownCallback: shutdownCallback,
		logger:           logger,
		startErrors:      startErrors,
		errChan:          errChan,
		halt:             halt,
		directives:       make(chan watcherDirective),
	}
}

func (w *processWatcher) watch() {
	go w.watchSignal()
	go w.watchErrors()
	go w.watchHaltChan()

	stopped := false
	for directive := range w.directives {
		if directive == directiveAbort {
			return
		}

		if !stopped {
			stopped = true
			w.logger.Info("Starting graceful shutdown")
			w.shutdownCallback()
		}
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
			w.directives <- directiveAbort
			return
		}

		urgent = true
		w.logger.Info("Received signal, starting graceful shutdown")
		w.directives <- directiveShutdown
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
			w.logger.Error("%s returned a fatal error", err.process.Name())
			w.errChan <- err.err
		}

		w.directives <- directiveShutdown
	}

	w.directives <- directiveShutdown
}

func (w *processWatcher) watchHaltChan() {
	<-w.halt
	w.logger.Info("Received external shutdown request")
	w.directives <- directiveShutdown
}
