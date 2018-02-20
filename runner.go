package nacelle

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

var (
	ErrCleanShutdown  = errors.New("process stopped cleanly")
	ErrUrgentShutdown = errors.New("urgent shutdown requested")
)

type (
	ProcessRunner struct {
		container    *ServiceContainer
		initializers []*initializerMeta
		processes    map[int][]*processMeta
		numProcesses int
		done         chan struct{}
		halt         chan struct{}
		once         *sync.Once
	}

	errMeta struct {
		err     error
		process *processMeta
	}
)

func NewProcessRunner(container *ServiceContainer) *ProcessRunner {
	return &ProcessRunner{
		container:    container,
		initializers: []*initializerMeta{},
		processes:    map[int][]*processMeta{},
		done:         make(chan struct{}),
		halt:         make(chan struct{}),
		once:         &sync.Once{},
	}
}

func (pr *ProcessRunner) RegisterInitializer(initializer Initializer, initializerConfigs ...InitializerConfigFunc) {
	meta := &initializerMeta{Initializer: initializer}

	for _, f := range initializerConfigs {
		f(meta)
	}

	pr.initializers = append(pr.initializers, meta)
}

func (pr *ProcessRunner) RegisterProcess(process Process, processConfigs ...ProcessConfigFunc) {
	meta := &processMeta{Process: process}

	for _, f := range processConfigs {
		f(meta)
	}

	if _, ok := pr.processes[meta.priority]; !ok {
		pr.processes[meta.priority] = []*processMeta{}
	}

	pr.numProcesses++
	pr.processes[meta.priority] = append(pr.processes[meta.priority], meta)
}

func (pr *ProcessRunner) Run(config Config, logger Logger) <-chan error {
	errChan := make(chan error, pr.numProcesses*2+1)

	if err := pr.runInitializers(config, logger); err != nil {
		defer close(errChan)
		errChan <- err
		return errChan
	}

	var (
		startErrors = make(chan errMeta)
		priorities  = pr.getPriorities()
		wg          = &sync.WaitGroup{}
	)

	if !pr.runProcesses(priorities, config, logger, startErrors, errChan, wg) {
		return errChan
	}

	logger.Info("All processes running")

	go pr.watch(priorities, logger, startErrors, errChan)
	go closeAfterWait(wg, startErrors)

	return chainUntilHalt(errChan, pr.done)
}

func (pr *ProcessRunner) getPriorities() []int {
	priorities := []int{}
	for priority := range pr.processes {
		priorities = append(priorities, priority)
	}

	sort.Ints(priorities)
	return priorities
}

func (pr *ProcessRunner) runInitializers(config Config, logger Logger) error {
	logger.Info("Running initializers")

	for _, initializer := range pr.initializers {
		logger.Debug("Injecting services into %s", initializer.Name())

		if err := pr.container.Inject(initializer.Initializer); err != nil {
			return fmt.Errorf(
				"failed to inject services into %s (%s)",
				initializer.Name(),
				err.Error(),
			)
		}

		logger.Debug("Initializing %s", initializer.Name())

		if err := initializer.Init(config); err != nil {
			return fmt.Errorf(
				"failed to initialize %s (%s)",
				initializer.Name(),
				err.Error(),
			)
		}

		logger.Debug("Initialized %s", initializer.Name())
	}

	return nil
}

func (pr *ProcessRunner) runProcesses(
	priorities []int,
	config Config,
	logger Logger,
	startErrors chan errMeta,
	errChan chan error,
	wg *sync.WaitGroup,
) bool {
	logger.Debug("Injecting services into process instances")

	for i := range priorities {
		for _, process := range pr.processes[priorities[i]] {
			if err := pr.container.Inject(process.Process); err != nil {
				defer close(errChan)

				errChan <- fmt.Errorf(
					"failed to inject services into %s (%s)",
					process.Name(),
					err.Error(),
				)

				return false
			}
		}
	}

	logger.Info("Initializing and starting processes")

	for i := range priorities {
		err := pr.initAndStartProcesses(
			pr.processes[priorities[i]],
			priorities[i],
			config,
			logger,
			wg,
			startErrors,
		)

		if err != nil {
			errChan <- err
			pr.stopProcesessBelowPriority(priorities, i, logger, errChan)
			go closeAfterWait(wg, startErrors)

			go func() {
				defer close(errChan)

				for err := range startErrors {
					if err.err != nil {
						errChan <- err.err
					}
				}
			}()

			return false
		}
	}

	return true
}

func (pr *ProcessRunner) initAndStartProcesses(
	processes []*processMeta,
	priority int,
	config Config,
	logger Logger,
	wg *sync.WaitGroup,
	startErrors chan<- errMeta,
) error {
	logger.Debug("Initializing processes at priority %d", priority)

	for _, process := range processes {
		logger.Debug("Initializing %s", process.Name())

		if err := process.Init(config); err != nil {
			return fmt.Errorf("failed to initialize %s (%s)", process.Name(), err.Error())
		}

		logger.Debug("Initialized %s", process.Name())
	}

	logger.Debug("Starting processes at priority %d", priority)

	for _, process := range processes {
		wg.Add(1)

		go func(process *processMeta) {
			defer wg.Done()

			logger.Debug("Starting %s", process.Name())

			err := process.Start()
			if err != nil {
				err = fmt.Errorf("%s returned a fatal error (%s)", process.Name(), err.Error())
			}

			startErrors <- errMeta{err, process}
		}(process)
	}

	return nil
}

func (pr *ProcessRunner) watch(
	priorities []int,
	logger Logger,
	startErrors <-chan errMeta,
	errChan chan<- error,
) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

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
				logger.Info("Received second signal, no longer waiting for graceful exit")
				return
			}

			logger.Info("Received signal, starting graceful shutdown")
			urgent = true

		case err, ok := <-startErrors:
			if !ok {
				return
			}

			if err.err == nil {
				if err.process.silentExit {
					continue
				}

				logger.Info(
					"%s has stopped cleanly, starting graceful shutdown",
					err.process.Name(),
				)
			} else {
				logger.Error(
					"%s returned a fatal error, starting graceful shutdown",
					err.process.Name(),
				)

				errChan <- err.err
			}

		case <-pr.halt:
			logger.Info("Received external shutdown request")
		}

		if !stopped {
			stopped = true
			pr.stopProcesessBelowPriority(priorities, len(priorities), logger, errChan)
		}
	}
}

func (pr *ProcessRunner) Shutdown(timeout time.Duration) error {
	pr.once.Do(func() {
		close(pr.halt)
	})

	select {
	case <-time.After(timeout):
		return errors.New("process failed to stop in timeout")
	case <-pr.done:
		return nil
	}
}

func (pr *ProcessRunner) stopProcesessBelowPriority(priorities []int, p int, logger Logger, errChan chan<- error) {
	for i := p - 1; i >= 0; i-- {
		pr.stopProcesses(pr.processes[priorities[i]], priorities[i], logger, errChan)
	}
}

func (pr *ProcessRunner) stopProcesses(processes []*processMeta, priority int, logger Logger, errChan chan<- error) {
	logger.Debug("Stopping processes at priority %d", priority)

	for _, process := range processes {
		logger.Debug("Stopping %s", process.Name())

		if err := process.Stop(); err != nil {
			errChan <- fmt.Errorf("%s returned error from stop (%s)", process.Name(), err.Error())
		}
	}
}

//
// Helpers

func closeAfterWait(wg *sync.WaitGroup, startErrors chan errMeta) {
	wg.Wait()
	close(startErrors)
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
