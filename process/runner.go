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
	// Runner wraps a process container. Given a loaded configuration object,
	// it can run the registered initializers and processes and wait for them
	// to exit (cleanly or via shutdown request).
	Runner interface {
		// Run takes a loaded configuration object, then starts and monitors
		// the registered items in the process container. This method returns
		// a channel of errors. Each error from an initializer or a process will
		// be sent on this channel (nil errors are ignored). This channel will
		// close once all processes have exited (or, alternatively, when the
		// shutdown timeout has elapsed).
		Run(config.Config) <-chan error

		// Shutdown will begin a graceful exit of all processes. This method
		// will block until the runner has exited (the channel from the Run
		// method has closed) or the given duration has elapsed. In the later
		// case a non-nil error is returned.
		Shutdown(time.Duration) error
	}

	runner struct {
		processes       Container
		services        service.Container
		watcher         *processWatcher
		errChan         chan errMeta
		outChan         chan error
		wg              *sync.WaitGroup
		shutdownTimeout time.Duration
		logger          logging.Logger
		clock           glock.Clock
	}

	namedInjectable interface {
		Name() string
		Wrapped() interface{}
	}

	namedInitializer interface {
		Initializer
		Name() string
		InitTimeout() time.Duration
	}

	errMeta struct {
		err        error
		source     namedInitializer
		silentExit bool
	}
)

// NewRunner creates a process runner from the given process and service
// containers.
func NewRunner(
	processes Container,
	services service.Container,
	health Health,
	runnerConfigs ...RunnerConfigFunc,
) Runner {
	// There can be one init error plus one start and one stop error
	// per started process. Make the output channel buffer as large
	// as the maximum number of errors.
	maxErrs := processes.NumProcesses()*2 + 1

	errChan := make(chan errMeta)
	outChan := make(chan error, maxErrs)

	r := &runner{
		processes: processes,
		services:  services,
		errChan:   errChan,
		outChan:   outChan,
		wg:        &sync.WaitGroup{},
		logger:    logging.NewNilLogger(),
		clock:     glock.NewRealClock(),
	}

	for _, f := range runnerConfigs {
		f(r)
	}

	// Create a watcher around the meta error channel (written to by
	// the runner) and the output channel (read by the boot process).
	// Pass our own logger and clock instances and the requested
	// shutdown timeout.

	r.watcher = newWatcher(
		errChan,
		outChan,
		withWatcherLogger(r.logger),
		withWatcherClock(r.clock),
		withWatcherShutdownTimeout(r.shutdownTimeout),
	)

	return r
}

func (r *runner) Run(config config.Config) <-chan error {
	// Start watching things before running anything. This ensures that
	// we start listening for shutdown requests and intercepted signals
	// as soon as anything starts being initialized.

	r.watcher.watch()

	// Run the initializers in sequence. If there were no errors, begin
	// initializing and running processes in priority/registration order.

	_ = r.runInitializers(config) && r.runProcesses(config)
	return r.outChan
}

func (r *runner) Shutdown(timeout time.Duration) error {
	r.watcher.halt()

	select {
	case <-r.clock.After(timeout):
		return fmt.Errorf("process runner did not shutdown within timeout")
	case <-r.watcher.done:
		return nil
	}
}

//
// Running and Watching

func (r *runner) runInitializers(config config.Config) bool {
	r.logger.Info("Running initializers")

	for _, initializer := range r.processes.GetInitializers() {
		if err := r.inject(initializer); err != nil {
			r.errChan <- errMeta{err, initializer, false}
			close(r.errChan)
			return false
		}

		if err := r.initWithTimeout(initializer, config); err != nil {
			r.errChan <- errMeta{err, initializer, false}
			close(r.errChan)
			return false
		}
	}

	return true
}

func (r *runner) runProcesses(config config.Config) bool {
	r.logger.Info("Running processes")

	if !r.injectProcesses() {
		return false
	}

	// For each priority index, attempt to initialize the processes
	// in sequence. Then, start all processes in a goroutine. If there
	// is any synchronous error occurs (either due to an Init call
	// returning a non-nil error, or the watcher has begun shutdown),
	// stop booting up processes adn simply wait for them to spin down.

	success := true
	for index := 0; index < r.processes.NumPriorities(); index++ {
		if !r.initProcessesAtPriorityIndex(config, index) {
			success = false
			break
		}

		r.startProcessesAtPriorityIndex(index)
	}

	// Wait for all booted processes to exit any Start/Stop methods,
	// then close the channel that receives their errors. This will
	// signal to the watcher to do its own cleanup (and close the
	// output channel).

	go func() {
		r.wg.Wait()
		close(r.errChan)
	}()

	if !success {
		return false
	}

	r.logger.Info("All processes have started")
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

func (r *runner) initProcessesAtPriorityIndex(config config.Config, index int) bool {
	r.logger.Info("Initializing processes at priority index %d", index)

	for _, process := range r.processes.GetProcessesAtPriorityIndex(index) {
		if err := r.initWithTimeout(process, config); err != nil {
			r.errChan <- errMeta{err, process, false}
			return false
		}
	}

	return true
}

func (r *runner) initWithTimeout(initializer namedInitializer, config config.Config) error {
	// Run the initializer in a goroutine. Write its return value
	// to a buffered channel so that it does not block if we happen
	// to abandon the read (on timeout or during shutdown).

	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		ch <- r.init(initializer, config)
	}()

	select {
	case err := <-ch:
		// Init completed, return its value
		return err

	case <-r.watcher.shutdownSignal:
		// Watcher is shutting down, ignore the return value of this call
		return fmt.Errorf("aborting initialization of %s", initializer.Name())

	case <-r.makeTimeoutChan(initializer.InitTimeout()):
		// Initialization took too long, return an error
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

func (r *runner) startProcessesAtPriorityIndex(index int) {
	r.logger.Info("Starting processes at priority index %d", index)

	// For each process group, we create a goroutine that will shutdown
	// all processes once the watcher begins shutting down. We add one to
	// the wait group to "bridge the gap" between the exit of the start
	// methods and the call to a stop method -- this situation is likely
	// rare, but would cause a panic.

	r.wg.Add(1)

	go func() {
		defer r.wg.Done()
		<-r.watcher.shutdownSignal
		r.stopProcessesAtPriorityIndex(index)
	}()

	for _, process := range r.processes.GetProcessesAtPriorityIndex(index) {
		r.wg.Add(1)

		go func(p *ProcessMeta) {
			defer r.wg.Done()
			r.startProcess(p)
		}(process)
	}
}

func (r *runner) startProcess(process *ProcessMeta) {
	r.logger.Info("Starting %s", process.Name())

	if err := process.Start(); err != nil {
		wrappedErr := fmt.Errorf(
			"%s returned a fatal error (%s)",
			process.Name(),
			err.Error(),
		)

		r.errChan <- errMeta{wrappedErr, process, false}
	} else {
		r.errChan <- errMeta{nil, process, process.silentExit}
	}
}

//
// Process Stopping

func (r *runner) stopProcessesAtPriorityIndex(index int) {
	r.logger.Info("Stopping processes at priority index %d", index)

	// Call stop on all processes at this priority index in parallel. We
	// add one to the wait group for each routine to ensure that we do
	// not close the err channel until all possible error producers have
	// exited.

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

//
// Helpers

func (r *runner) makeTimeoutChan(timeout time.Duration) <-chan time.Time {
	if timeout == 0 {
		return nil
	}

	return r.clock.After(timeout)
}
