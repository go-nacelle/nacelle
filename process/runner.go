package process

import (
	"fmt"
	"sync"
	"time"

	"github.com/efritz/glock"
	"github.com/efritz/nacelle/config"
	"github.com/efritz/nacelle/logging"
	"github.com/efritz/nacelle/service"
)

type (
	Runner interface {
		Run(config.Config) <-chan error
		Shutdown(time.Duration) error
	}

	runner struct {
		processes       Container
		services        service.Container
		logger          logging.Logger
		clock           glock.Clock
		done            chan struct{}
		halt            chan struct{}
		once            *sync.Once
		errChan         chan errMeta
		wg              *sync.WaitGroup
		shutdownTimeout time.Duration
	}
)

func NewRunner(
	processes Container,
	services service.Container,
	runnerConfigs ...RunnerConfigFunc,
) Runner {
	runner := &runner{
		processes: processes,
		services:  services,
		logger:    logging.NewNilLogger(),
		clock:     glock.NewRealClock(),
		done:      make(chan struct{}),
		halt:      make(chan struct{}),
		once:      &sync.Once{},
		errChan:   make(chan errMeta),
		wg:        &sync.WaitGroup{},
	}

	for _, f := range runnerConfigs {
		f(runner)
	}

	return runner
}

func (r *runner) Run(config config.Config) <-chan error {
	outChan := make(chan error, r.processes.NumProcesses()*2+1)

	watcher := newWatcher(
		r.errChan,
		outChan,
		r.halt,
		withWatcherLogger(r.logger),
		withWatcherClock(r.clock),
		withWatcherShutdownTimeout(r.shutdownTimeout),
	)

	watcher.watch()

	if r.runInitializers(config, watcher) {
		r.runProcesses(config, watcher)
	}

	go watcher.wait()
	return outChan
}

func (r *runner) Shutdown(timeout time.Duration) error {
	r.once.Do(func() {
		close(r.halt)
	})

	select {
	case <-time.After(timeout):
		return fmt.Errorf("process runner did not shutdown within timeout")
	case <-r.done:
		return nil
	}
}

//
// Running and Watching

func (r *runner) runInitializers(config config.Config, watcher *processWatcher) bool {
	r.logger.Info("Running initializers")

	for _, initializer := range r.processes.GetInitializers() {
		if err := r.inject(initializer); err != nil {
			r.errChan <- errMeta{err, initializer, false}
			close(r.errChan)
			return false
		}

		if err := r.initWithTimeout(initializer, config, watcher); err != nil {
			r.errChan <- errMeta{err, initializer, false}
			close(r.errChan)
			return false
		}
	}

	return true
}

func (r *runner) runProcesses(config config.Config, watcher *processWatcher) bool {
	r.logger.Info("Running processes")

	if !r.injectProcesses() {
		return false
	}

	var (
		success = true
	)

	for index := 0; index < r.processes.NumPriorities(); index++ {
		if !r.initProcessesAtPriorityIndex(config, watcher, index) {
			success = false
			break
		}

		r.startProcessesAtPriorityIndex(
			watcher,
			index,
		)
	}

	go func() {
		r.wg.Wait()
		close(r.errChan)
	}()

	if success {
		r.logger.Info("All processes have started")
	}

	return true
}

//
// Injection

func (r *runner) injectProcesses() bool {
	for i := 0; i < r.processes.NumPriorities(); i++ {
		for _, process := range r.processes.GetProcessesAtPriorityIndex(i) {
			if err := r.inject(process); err != nil {
				r.errChan <- errMeta{err, process, false}
				close(r.errChan)
				return false
			}
		}
	}

	return true
}

func (r *runner) inject(v namedInjectable) error {
	r.logger.Info("Injecting services into %s", v.Name())

	if err := r.services.Inject(v.Wrapped()); err != nil {
		return fmt.Errorf(
			"failed to inject services into %s (%s)",
			v.Name(),
			err.Error(),
		)
	}

	return nil
}

//
// Initialization

func (r *runner) initProcessesAtPriorityIndex(
	config config.Config,
	watcher *processWatcher,
	index int,
) bool {
	r.logger.Info("Initializing processes at priority index %d", index)

	for _, process := range r.processes.GetProcessesAtPriorityIndex(index) {
		if err := r.initWithTimeout(process, config, watcher); err != nil {
			r.errChan <- errMeta{err, process, false}
			return false
		}
	}

	return true
}

func (r *runner) initWithTimeout(
	initializer namedInitializer,
	config config.Config,
	watcher *processWatcher,
) error {
	var timeoutCh <-chan time.Time
	if timeout := initializer.InitTimeout(); timeout > 0 {
		timeoutCh = r.clock.After(timeout)
	}

	ch := make(chan error)

	go func() {
		defer close(ch)
		ch <- r.init(initializer, config)
	}()

	select {
	case err := <-ch:
		return err

	case <-watcher.draining:
		return fmt.Errorf("aborting initialization of %s", initializer.Name())

	case <-timeoutCh:
		return fmt.Errorf("%s did not initialize within timeout", initializer.Name())
	}
}

func (r *runner) init(initializer namedInitializer, config config.Config) error {
	r.logger.Info("Initializing %s", initializer.Name())

	if err := initializer.Init(config); err != nil {
		return fmt.Errorf(
			"failed to initialize %s (%s)",
			initializer.Name(),
			err.Error(),
		)
	}

	r.logger.Info("Initialized %s", initializer.Name())
	return nil
}

//
// Process Starting

func (r *runner) startProcessesAtPriorityIndex(watcher *processWatcher, index int) {
	r.logger.Info("Starting processes at priority index %d", index)

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		<-watcher.draining

		r.stopProcessesAtPriorityIndex(
			index,
		)
	}()

	for _, process := range r.processes.GetProcessesAtPriorityIndex(index) {
		r.wg.Add(1)

		go func(process *ProcessMeta) {
			defer r.wg.Done()

			r.logger.Info("Starting %s", process.Name())
			err := process.Start()

			if err != nil {
				err = fmt.Errorf(
					"%s returned a fatal error (%s)",
					process.Name(),
					err.Error(),
				)
			}

			r.errChan <- errMeta{err, process, process.silentExit}
		}(process)
	}
}

//
// Process Stopping

func (r *runner) stopProcessesAtPriorityIndex(index int) {
	r.logger.Info("Stopping processes at priority index %d", index)

	for _, process := range r.processes.GetProcessesAtPriorityIndex(index) {
		r.wg.Add(1)
		go func(process *ProcessMeta) {
			defer r.wg.Done()

			if err := r.stop(process); err != nil {
				r.errChan <- errMeta{err, process, false}
			}
		}(process)
	}
}

func (r *runner) stop(process *ProcessMeta) error {
	r.logger.Info("Stopping %s", process.Name())

	if err := process.Stop(); err != nil {
		return fmt.Errorf(
			"%s returned error from stop (%s)",
			process.Name(),
			err.Error(),
		)
	}

	return nil
}
