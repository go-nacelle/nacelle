package main

import (
	"time"

	"github.com/go-nacelle/nacelle"
	"github.com/go-nacelle/nacelle/base/worker"
)

type spec struct {
	Logger nacelle.Logger `service:"logger"`
	halt   <-chan struct{}
	count  int
}

func (s *spec) Init(config nacelle.Config, worker *worker.Worker) error {
	s.halt = worker.HaltChan()
	return nil
}

func (s *spec) Tick() error {
	select {
	case <-s.halt:
		s.Logger.Warning("Aborting tick")
		return nil

		// Blocking operations should use the worker's halt channel
		// so that the worker can shutdown gracefully when the app
		// is attempting to drain.
	case <-time.After(time.Second):
	}

	s.count++
	s.Logger.Info("Tick #%d", s.count)
	return nil
}

//
//

func setup(processes nacelle.ProcessContainer, services nacelle.ServiceContainer) error {
	processes.RegisterProcess(worker.NewWorker(&spec{}))
	return nil
}

//
//

func main() {
	nacelle.NewBootstrapper("worker-example", setup).BootAndExit()
}
