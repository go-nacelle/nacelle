package process

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/efritz/backoff"
	"github.com/efritz/glock"
	"github.com/efritz/watchdog"

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
		processes          Container
		services           service.Container
		health             Health
		watcher            *processWatcher
		errChan            chan errMeta
		outChan            chan error
		wg                 *sync.WaitGroup
		logger             logging.Logger
		clock              glock.Clock
		startupTimeout     time.Duration
		shutdownTimeout    time.Duration
		healthCheckBackoff backoff.Backoff
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

	namedFinalizer interface {
		Initializer
		Name() string
		FinalizeTimeout() time.Duration
		Wrapped() interface{}
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
	// per started process plus one finalizer error per initializer.
	// Make the output channel buffer as large as the maximum number
	// of errors.
	maxErrs := processes.NumInitializers() + processes.NumProcesses()*2 + 1

	errChan := make(chan errMeta)
	outChan := make(chan error, maxErrs)

	r := &runner{
		processes:          processes,
		services:           services,
		health:             health,
		errChan:            errChan,
		outChan:            outChan,
		wg:                 &sync.WaitGroup{},
		logger:             logging.NewNilLogger(),
		clock:              glock.NewRealClock(),
		healthCheckBackoff: backoff.NewConstantBackoff(time.Second),
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

	for i, initializer := range r.processes.GetInitializers() {
		if err := r.inject(initializer); err != nil {
			_ = r.unwindInitializers(i)
			r.errChan <- errMeta{err: err, source: initializer}
			close(r.errChan)
			return false
		}

		if err := r.initWithTimeout(initializer, config); err != nil {
			_ = r.unwindInitializers(i)
			r.errChan <- errMeta{err: err, source: initializer}
			close(r.errChan)
			return false
		}
	}

	return true
}

func (r *runner) runFinalizers() bool {
	return r.unwindInitializers(r.processes.NumInitializers())
}

func (r *runner) unwindInitializers(beforeIndex int) bool {
	r.logger.Info("Running finalizers")

	success := true
	initializers := r.processes.GetInitializers()

	for i := beforeIndex - 1; i >= 0; i-- {
		if err := r.finalizeWithTimeout(initializers[i]); err != nil {
			r.errChan <- errMeta{err: err, source: initializers[i]}
			success = false

		}
	}

	return success
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

		if !r.startProcessesAtPriorityIndex(index) {
			success = false
			break
		}
	}

	// Wait for all booted processes to exit any Start/Stop methods,
	// then run all the initializers with finalize methods in their
	// reverse startup order. After all possible writes to the error
	// channel have occurred, close it to signal to the watcher to
	// do its own cleanup (and close the output channel).

	go func() {
		r.wg.Wait()
		_ = r.runFinalizers()
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
				r.errChan <- errMeta{err: err, source: process}
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
			r.errChan <- errMeta{err: err, source: process}
			return false
		}
	}

	return true
}

func (r *runner) initWithTimeout(initializer namedInitializer, config config.Config) error {
	// Run the initializer in a goroutine. We don't want to block
	// on this in case we want to abandon reading from this channel
	// (timeout or shutdown). This is only true for initializer
	// methods (will not be true for process Start methods).

	errChan := makeErrChan(func() error {
		return r.init(initializer, config)
	})

	// Construct a timeout chan for the init (if timeout is set to
	// zero, this chan is nil and will never yield a value).

	initTimeoutChan := r.makeTimeoutChan(initializer.InitTimeout())

	// Now, wait for one of three results:
	//   - Init completed, return its value
	//   - Initialization took too long, return an error
	//   - Watcher is shutting down, ignore the return value

	select {
	case err := <-errChan:
		return err

	case <-initTimeoutChan:
		return fmt.Errorf("%s did not initialize within timeout", initializer.Name())

	case <-r.watcher.shutdownSignal:
		return fmt.Errorf("aborting initialization of %s", initializer.Name())
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
// Finalization

func (r *runner) finalizeWithTimeout(initializer namedFinalizer) error {
	// Similar to initWithTimeout, run the finalizer in a goroutine
	// and either return the error result or return an error value
	// if the finalizer took too long.

	errChan := makeErrChan(func() error {
		return r.finalize(initializer)
	})

	finalizeTimeoutChan := r.makeTimeoutChan(initializer.FinalizeTimeout())

	select {
	case err := <-errChan:
		return err

	case <-finalizeTimeoutChan:
		return fmt.Errorf("%s did not finalize within timeout", initializer.Name())
	}
}

func (r *runner) finalize(initializer namedFinalizer) error {
	// Finalizer is an optional interface on Initializer. Skip
	// this method if this initializer doesn't have the proper
	// method.

	finalizer, ok := initializer.Wrapped().(Finalizer)
	if !ok {
		return nil
	}

	r.logger.Info("Finalizing %s", initializer.Name())

	if err := finalizer.Finalize(); err != nil {
		return fmt.Errorf(
			"%s returned error from finalize (%s)",
			initializer.Name(),
			err.Error(),
		)
	}

	r.logger.Info("Finalized %s", initializer.Name())
	return nil
}

//
// Process Starting

func (r *runner) startProcessesAtPriorityIndex(index int) bool {
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

	// Create a context object whose cancellation marks either all processes
	// of this priority index going healthy or the startup timeout elapsing.
	// Create an abandon channel that closes at the same time to signal the
	// routine invoking the Start method of a process to ignore its return
	// value -- we really should not allow a timed-out start method to block
	// the entire process.
	//
	// The timeout is calculated from the minimum timeout values of the runner
	// and each process at this priority index. If no such values are set, then
	// the channel is nil and will never yield a value.

	var (
		ctx, cancel        = context.WithCancel(context.Background())
		abandonSignal      = make(chan struct{})
		minTimeout         = r.startupTimeoutForPriorityIndex(index)
		startupTimeoutChan = r.makeTimeoutChan(minTimeout)
	)

	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-startupTimeoutChan:
			r.logger.Error(
				"processes at priority index %d did not start within timeout",
				index,
			)

			cancel()
			close(abandonSignal)
		}
	}()

	// Actually start each process. Each call to startProcess blocks until
	// the process exits, so we perform each one in a goroutine guarded by
	// the runner's wait group.

	for _, process := range r.processes.GetProcessesAtPriorityIndex(index) {
		r.wg.Add(1)

		go func(p *ProcessMeta) {
			defer r.wg.Done()
			r.startProcess(p, abandonSignal)
		}(process)
	}

	// The health check function will read the outstanding unhealthy reasons
	// from the health reporter. We will block until the reason list is empty.

	retry := watchdog.RetryFunc(func() bool {
		if descriptions := r.getHealthDescriptions(); len(descriptions) > 0 {
			r.logger.Warning(
				"Process is not yet healthy - outstanding reasons: %s",
				strings.Join(descriptions, ", "),
			)

			return false
		}

		return true
	})

	// Perform the health check function with the given backoff. Will return
	// false if the context has been canceled before the health function has
	// returned true.

	if !watchdog.BlockUntilSuccess(ctx, retry, r.healthCheckBackoff) {
		// Do one last check to ensure that there are still unhealthy reasons.
		// This stops a weird race condition where we stop the runner with a
		// series of process errors with no outstanding reasons.

		if descriptions := r.getHealthDescriptions(); len(descriptions) > 0 {
			err := fmt.Errorf(
				"process did not become healthy within timeout - outstanding reasons: %s",
				strings.Join(descriptions, ", "),
			)

			r.errChan <- errMeta{err: err}
			return false
		}
	}

	r.logger.Info(
		"All processes at priority index %d have reported healthy",
		index,
	)

	return true
}

func (r *runner) startupTimeoutForPriorityIndex(index int) time.Duration {
	timeout := r.startupTimeout

	for _, process := range r.processes.GetProcessesAtPriorityIndex(index) {
		if process.startTimeout != 0 && (timeout == 0 || process.startTimeout < timeout) {
			timeout = process.startTimeout
		}
	}

	return timeout
}

func (r *runner) getHealthDescriptions() []string {
	descriptions := []string{}
	for _, reason := range r.health.Reasons() {
		descriptions = append(descriptions, fmt.Sprintf("%s", reason.Key))
	}

	return descriptions
}

func (r *runner) startProcess(process *ProcessMeta, abandonSignal <-chan struct{}) {
	r.logger.Info("Starting %s", process.Name())

	// Run the start method in a goroutine. We need to do
	// this as we assume all processes are long-running
	// and need to read from other sources for shutdown
	// and timeout behavior.

	errChan := makeErrChan(func() error {
		return process.Start()
	})

	// Create a channel for the shutdown timeout. This
	// channel will close only after the timeout duration
	// elapses AFTER the stop method of the process is
	// called. If the shutdown timeout is set to zero, this
	// channel will remain nil and will never yield.

	var shutdownTimeout chan (struct{})

	if process.shutdownTimeout > 0 {
		shutdownTimeout = make(chan struct{})

		go func() {
			<-process.stopped
			<-r.clock.After(process.shutdownTimeout)
			close(shutdownTimeout)
		}()
	}

	// Now, wait for the Start method to yield, in which case the error value
	// is passed to the watcher, or for either the abandon channel or timeout
	// channel to signal, in which case we abandon the reading of the return
	// value from the Start method.

	select {
	case <-abandonSignal:
		r.logger.Error("Abandoning result of %s", process.Name())
		return

	case err := <-errChan:
		if err != nil {
			wrappedErr := fmt.Errorf(
				"%s returned a fatal error (%s)",
				process.Name(),
				err.Error(),
			)

			r.errChan <- errMeta{err: wrappedErr, source: process}
		} else {
			r.errChan <- errMeta{err: nil, source: process, silentExit: process.silentExit}
		}

	case <-shutdownTimeout:
		wrappedErr := fmt.Errorf(
			"%s did not shutdown within timeout",
			process.Name(),
		)

		r.errChan <- errMeta{err: wrappedErr, source: process}
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
				r.errChan <- errMeta{err: err, source: process}
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

func makeErrChan(f func() error) <-chan error {
	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		ch <- f()
	}()

	return ch
}
