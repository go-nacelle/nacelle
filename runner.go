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

	initializerMeta struct {
		Initializer
	}

	processMeta struct {
		Process
		priority   int
		silentExit bool
	}

	errMeta struct {
		err     error
		process *processMeta
	}

	InitializerConfigFunc func(*initializerMeta)
	ProcessConfigFunc     func(*processMeta)
)

func WithPriority(priority int) ProcessConfigFunc {
	return func(meta *processMeta) { meta.priority = priority }
}

func WithSilentExit() ProcessConfigFunc {
	return func(meta *processMeta) { meta.silentExit = true }
}

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
	meta := &initializerMeta{
		Initializer: initializer,
	}

	for _, f := range initializerConfigs {
		f(meta)
	}

	pr.initializers = append(pr.initializers, meta)
}

func (pr *ProcessRunner) RegisterProcess(process Process, processConfigs ...ProcessConfigFunc) {
	meta := &processMeta{
		Process: process,
	}

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
	errChan := make(chan error, pr.numProcesses+1)

	logger.Info("Running initializers")

	if err := pr.runInitializers(config); err != nil {
		defer close(errChan)
		errChan <- err
		return errChan
	}

	var (
		startErrors = make(chan errMeta)
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
					if err.err == nil {
						continue
					}

					errChan <- err.err
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

				if err.err == nil {
					if err.process.silentExit {
						continue
					}

					logger.Info("Process has stopped cleanly, starting graceful shutdown")
				} else {
					logger.Info("Encountered error, starting graceful shutdown")
					errChan <- err.err
				}

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

func (pr *ProcessRunner) initAndStartProcesses(
	processes []*processMeta,
	priority int,
	config Config,
	logger Logger,
	wg *sync.WaitGroup,
	errors chan<- errMeta,
) error {
	logger.Info("Initializing processes at priority %d", priority)

	for _, process := range processes {
		if err := process.Init(config); err != nil {
			return err
		}
	}

	logger.Info("Starting processes at priority %d", priority)

	for _, process := range processes {
		wg.Add(1)

		go func(process *processMeta) {
			defer wg.Done()
			err := process.Start()
			errors <- errMeta{err, process}
		}(process)
	}

	return nil
}

func (pr *ProcessRunner) stopProcesessBelowPriority(priorities []int, p int, logger Logger, errChan chan<- error) {
	for i := p - 1; i >= 0; i-- {
		pr.stopProcesses(pr.processes[priorities[i]], priorities[i], logger, errChan)
	}
}

func (pr *ProcessRunner) stopProcesses(processes []*processMeta, priority int, logger Logger, errChan chan<- error) {
	logger.Info("Stopping processes at priority %d", priority)

	for _, process := range processes {
		errChan <- process.Stop()
	}
}

func closeAfterWait(wg *sync.WaitGroup, ch chan errMeta) {
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
