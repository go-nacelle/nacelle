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

// ParallelInitializer is a container for initializers that are initialized
// in parallel. This is useful when groups of initializers are independent
// and may contain some longer-running process (such as dialing a remote
// service).
type ParallelInitializer struct {
	Logger       logging.Logger    `service:"logger"`
	Services     service.Container `service:"container"`
	clock        glock.Clock
	initializers []*InitializerMeta
}

// NewParallelInitializer creates a new parallel initializer.
func NewParallelInitializer(initializerConfigs ...ParallelInitializerConfigFunc) *ParallelInitializer {
	pi := &ParallelInitializer{
		initializers: []*InitializerMeta{},
	}

	for _, f := range initializerConfigs {
		f(pi)
	}

	return pi
}

// RegisterInitializer adds an initializer to the initializer set
// with the given configuration.
func (i *ParallelInitializer) RegisterInitializer(
	initializer Initializer,
	initializerConfigs ...InitializerConfigFunc,
) {
	meta := newInitializerMeta(initializer)

	for _, f := range initializerConfigs {
		f(meta)
	}

	i.initializers = append(i.initializers, meta)
}

// Init runs Init on all registered initializers concurrently.
func (pi *ParallelInitializer) Init(config config.Config) error {
	for _, initializer := range pi.initializers {
		if err := pi.injectAll(initializer); err != nil {
			return errMetaSet{
				errMeta{err: err, source: initializer},
			}
		}
	}

	errMetas := errMetaSet{}
	initErrors := pi.initializeAll(config)

	for i, err := range initErrors {
		if err != nil {
			errMetas = append(errMetas, errMeta{err: err, source: pi.initializers[i]})
		}
	}

	if len(errMetas) > 0 {
		for i, err := range pi.finalizeAll(initErrors) {
			if err != nil {
				errMetas = append(errMetas, errMeta{err: err, source: pi.initializers[i]})
			}
		}

		return errMetas
	}

	return nil
}

// Finalize runs Finalize on all registered initializers concurrently.
func (pi *ParallelInitializer) Finalize() error {
	errMetas := errMetaSet{}
	for i, err := range pi.finalizeAll(make([]error, len(pi.initializers))) {
		if err != nil {
			errMetas = append(errMetas, errMeta{err: err, source: pi.initializers[i]})
		}
	}

	if len(errMetas) > 0 {
		return errMetas
	}

	return nil
}

func (pi *ParallelInitializer) injectAll(initializer namedFinalizer) error {
	pi.Logger.Info("Injecting services into %s", initializer.Name())

	if err := pi.Services.Inject(initializer.Wrapped()); err != nil {
		return fmt.Errorf(
			"failed to inject services into %s (%s)",
			initializer.Name(),
			err.Error(),
		)
	}

	return nil
}

func (pi *ParallelInitializer) initializeAll(config config.Config) []error {
	errors := make([]error, len(pi.initializers))
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := range pi.initializers {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			if err := pi.initWithTimeout(pi.initializers[i], config); err != nil {
				mutex.Lock()
				errors[i] = err
				mutex.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return errors
}

func (pi *ParallelInitializer) finalizeAll(initErrors []error) []error {
	errors := make([]error, len(pi.initializers))
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := range pi.initializers {
		if initErrors[i] != nil {
			continue
		}

		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			if err := pi.finalizeWithTimeout(pi.initializers[i]); err != nil {
				mutex.Lock()
				errors[i] = err
				mutex.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return errors
}

func (pi *ParallelInitializer) initWithTimeout(initializer namedInitializer, config config.Config) error {
	errChan := makeErrChan(func() error {
		return pi.init(initializer, config)
	})

	select {
	case err := <-errChan:
		return err

	case <-pi.makeTimeoutChan(initializer.InitTimeout()):
		return fmt.Errorf("%s did not initialize within timeout", initializer.Name())
	}
}

func (pi *ParallelInitializer) init(initializer namedInitializer, config config.Config) error {
	pi.Logger.Info("Initializing %s", initializer.Name())

	if err := initializer.Init(config); err != nil {
		return fmt.Errorf(
			"failed to initialize %s (%s)",
			initializer.Name(),
			err.Error(),
		)
	}

	pi.Logger.Info("Initialized %s", initializer.Name())
	return nil
}

func (pi *ParallelInitializer) finalizeWithTimeout(initializer namedFinalizer) error {
	errChan := makeErrChan(func() error {
		return pi.finalize(initializer)
	})

	select {
	case err := <-errChan:
		return err

	case <-pi.makeTimeoutChan(initializer.FinalizeTimeout()):
		return fmt.Errorf("%s did not finalize within timeout", initializer.Name())
	}
}

func (pi *ParallelInitializer) finalize(initializer namedFinalizer) error {
	finalizer, ok := initializer.Wrapped().(Finalizer)
	if !ok {
		return nil
	}

	pi.Logger.Info("Finalizing %s", initializer.Name())

	if err := finalizer.Finalize(); err != nil {
		return fmt.Errorf(
			"%s returned error from finalize (%s)",
			initializer.Name(),
			err.Error(),
		)
	}

	pi.Logger.Info("Finalized %s", initializer.Name())
	return nil
}

func (pi *ParallelInitializer) makeTimeoutChan(timeout time.Duration) <-chan time.Time {
	return makeTimeoutChan(pi.clock, timeout)
}
