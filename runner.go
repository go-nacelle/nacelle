package nacelle

import (
	"errors"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

var ErrUrgentShutdown = errors.New("urgent shutdown requested")

type ProcessRunner struct {
	container    *ServiceContainer
	initializers []Initializer
	processes    map[int][]Process
	numProcesses int
	done         chan struct{}
	halt         chan struct{}
}

func NewProcessRunner(container *ServiceContainer) *ProcessRunner {
	return &ProcessRunner{
		container:    container,
		initializers: []Initializer{},
		processes:    map[int][]Process{},
		done:         make(chan struct{}),
		halt:         make(chan struct{}),
	}
}

func (pr *ProcessRunner) RegisterInitializer(initializer Initializer) {
	pr.initializers = append(pr.initializers, initializer)
}

func (pr *ProcessRunner) RegisterProcess(process Process, priority int) {
	if _, ok := pr.processes[priority]; !ok {
		pr.processes[priority] = []Process{}
	}

	pr.numProcesses++
	pr.processes[priority] = append(pr.processes[priority], process)
}

func (pr *ProcessRunner) Run(config Config, logger Logger) <-chan error {
	errChan := make(chan error, pr.numProcesses+1)

	logger.Info("Running initializers")

	if err := pr.runInitializers(config); err != nil {
		defer close(errChan)
		errChan <- err
		return errChan
	}

	var (
		startErrors = make(chan error)
		priorities  = pr.getPriorities()
		wg          = sync.WaitGroup{}
	)

	logger.Info("Injecting services to process instances")

	for i := range priorities {
		for _, process := range pr.processes[priorities[i]] {
			if err := pr.container.Inject(process); err != nil {
				defer close(errChan)
				errChan <- err
				return errChan
			}
		}
	}

	for i := range priorities {
		if err := pr.initAndStartProcesses(pr.processes[priorities[i]], priorities[i], config, logger, &wg, startErrors); err != nil {
			logger.Error("Encountered error starting process at priority %d", priorities[i])
			errChan <- err
			pr.stopProcesessBelowPriority(priorities, i, logger, errChan)
			go closeAfterWait(&wg, startErrors)

			go func() {
				defer close(errChan)

				for err := range startErrors {
					if err == nil {
						continue
					}

					errChan <- err
				}
			}()

			return errChan
		}
	}

	go closeAfterWait(&wg, startErrors)

	logger.Info("All processes have started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	go func() {
		defer close(errChan)
		defer close(pr.done)

		var (
			urgent  = false
			stopped = false
		)

		for {
			select {
			case <-sigChan:
				if urgent {
					logger.Info("Received second signal, shutting down NOW")
					return
				}

				logger.Info("Received signal, starting graceful shutdown")
				urgent = true

			case err, ok := <-startErrors:
				if !ok {
					return
				}

				if err == nil {
					continue
				}

				logger.Info("Encountered error, starting graceful shutdown")
				errChan <- err

			case <-pr.halt:
				logger.Info("Process requested shutdown")
			}

			if !stopped {
				stopped = true
				pr.stopProcesessBelowPriority(priorities, len(priorities), logger, errChan)
			}
		}
	}()

	return chainUntilHalt(errChan, pr.done)
}

func (pr *ProcessRunner) Shutdown(timeout time.Duration) error {
	// TODO - make idempotent
	close(pr.halt)

	select {
	case <-time.After(timeout):
		return errors.New("process failed to stop in timeout")
	case <-pr.done:
		return nil
	}
}

func (pr *ProcessRunner) getPriorities() []int {
	priorities := []int{}
	for priority := range pr.processes {
		priorities = append(priorities, priority)
	}

	sort.Ints(priorities)
	return priorities
}

func (pr *ProcessRunner) runInitializers(config Config) error {
	for _, initializer := range pr.initializers {
		if err := pr.container.Inject(initializer); err != nil {
			return err
		}

		if err := initializer.Init(config); err != nil {
			return err
		}
	}

	return nil
}

func (pr *ProcessRunner) initAndStartProcesses(processes []Process, priority int, config Config, logger Logger, wg *sync.WaitGroup, errors chan<- error) error {
	logger.Info("Initializing processes at priority %d", priority)

	for _, process := range processes {
		if err := process.Init(config); err != nil {
			return err
		}
	}

	logger.Info("Starting processes at priority %d", priority)

	for _, process := range processes {
		wg.Add(1)

		go func(process Process) {
			defer wg.Done()

			if err := process.Start(); err != nil {
				errors <- err
			}
		}(process)
	}

	return nil
}

func (pr *ProcessRunner) stopProcesessBelowPriority(priorities []int, p int, logger Logger, errChan chan<- error) {
	for i := p - 1; i >= 0; i-- {
		pr.stopProcesses(pr.processes[priorities[i]], priorities[i], logger, errChan)
	}
}

func (pr *ProcessRunner) stopProcesses(processes []Process, priority int, logger Logger, errChan chan<- error) {
	logger.Info("Stopping processes at priority %d", priority)

	for _, process := range processes {
		if err := process.Stop(); err != nil {
			errChan <- err
		}
	}
}

func closeAfterWait(wg *sync.WaitGroup, ch chan error) {
	wg.Wait()
	close(ch)
}

func chainUntilHalt(src <-chan error, halt <-chan struct{}) <-chan error {
	out := make(chan error)

	go func() {
	loop:
		for {
			select {
			case err, ok := <-src:
				if !ok {
					break loop
				}

				out <- err

			case <-halt:
				break loop
			}
		}

		close(out)

		for range src {
		}
	}()

	return out
}
