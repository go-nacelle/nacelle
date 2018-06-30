package process

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/efritz/glock"
	"github.com/efritz/nacelle/logging"
)

type (
	processWatcher struct {
		errChan         <-chan errMeta
		outChan         chan<- error
		halt            <-chan struct{}
		logger          logging.Logger
		clock           glock.Clock
		shutdownTimeout time.Duration
		directives      chan watcherDirective
		draining        chan struct{}
		done            chan struct{}
	}

	watcherDirective int
)

const (
	directiveAbort watcherDirective = iota
	directiveShutdown
)

var shutdownSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
}

func newWatcher(
	errChan <-chan errMeta,
	outChan chan<- error,
	halt <-chan struct{},
	configs ...watcherConfigFunc,
) *processWatcher {
	watcher := &processWatcher{
		errChan:    errChan,
		outChan:    outChan,
		halt:       halt,
		logger:     logging.NewNilLogger(),
		clock:      glock.NewRealClock(),
		directives: make(chan watcherDirective),
		draining:   make(chan struct{}),
		done:       make(chan struct{}),
	}

	for _, f := range configs {
		f(watcher)
	}

	return watcher
}

func (w *processWatcher) watch() {
	go w.control()
	go w.watchSignal()
	go w.watchErrors()
	go w.watchHaltChan()
	go w.watchShutdownTimeout()
}

func (w *processWatcher) wait() {
	<-w.done
	close(w.directives)
}

func (w *processWatcher) control() {
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
	defer close(w.done)
	defer close(w.outChan)

	for err := range w.errChan {
		if err.err == nil {
			w.logger.Info("%s has stopped cleanly", err.source.Name())

			if err.silentExit {
				continue
			}
		} else {
			w.logger.Error("%s returned a fatal error (%s)", err.source.Name(), err.err.Error())
			w.outChan <- err.err
		}

		w.shutdown()
	}

	w.shutdown()
}

func (w *processWatcher) watchHaltChan() {
	select {
	case <-w.done:
		return
	case <-w.halt:
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
	case <-w.draining:
	}

	select {
	case <-w.done:
	case <-w.clock.After(w.shutdownTimeout):
		w.abort()
	}
}

func (w *processWatcher) abort() {
	w.directives <- directiveAbort
}

func (w *processWatcher) shutdown() {
	w.directives <- directiveShutdown
}
