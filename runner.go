package nacelle

import (
	"errors"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
)

var ErrUrgentShutdown = errors.New("urgent shutdown requested")

type ProcessRunner struct {
	initializers []Initializer
	processes    map[int][]Process
	numProcesses int
}

func NewProcessRunner() *ProcessRunner {
	return &ProcessRunner{
		initializers: []Initializer{},
		processes:    map[int][]Process{},
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

func (pr *ProcessRunner) Run(config Config) <-chan error {
	errChan := make(chan error, pr.numProcesses+1)

	if err := pr.runInitializers(config); err != nil {
		errChan <- err
		defer close(errChan)
		return errChan
	}

	var (
		startErrors = make(chan error)
		priorities  = pr.getPriorities()
		wg          = &sync.WaitGroup{}
	)

	for i := range priorities {
		if err := pr.initAndStartProcesses(pr.processes[priorities[i]], config, wg, startErrors); err != nil {
			errChan <- err
			pr.stopProcesessBelowPriority(priorities, i, errChan)
			go closeAfterWait(wg, startErrors)

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

	go closeAfterWait(wg, startErrors)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	halt := make(chan struct{})

	go func() {
		defer close(errChan)
		defer close(halt)

		var (
			urgent  = false
			stopped = false
		)

		for {
			select {
			case <-sigChan:
				if urgent {
					errChan <- ErrUrgentShutdown
					return
				}

				urgent = true

			case err, ok := <-startErrors:
				if !ok {
					return
				}

				if err == nil {
					continue
				}

				errChan <- err
			}

			if !stopped {
				stopped = true
				pr.stopProcesessBelowPriority(priorities, len(priorities), errChan)
			}
		}
	}()

	return chainUntilHalt(errChan, halt)
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
		if err := initializer.Init(config); err != nil {
			return err
		}
	}

	return nil
}

func (pr *ProcessRunner) initAndStartProcesses(processes []Process, config Config, wg *sync.WaitGroup, errors chan<- error) error {
	for _, process := range processes {
		if err := process.Init(config); err != nil {
			return err
		}
	}

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

func (pr *ProcessRunner) stopProcesessBelowPriority(priorities []int, p int, errChan chan<- error) {
	for i := p - 1; i >= 0; i-- {
		pr.stopProcesses(pr.processes[priorities[i]], errChan)
	}
}

func (pr *ProcessRunner) stopProcesses(processes []Process, errChan chan<- error) {
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
