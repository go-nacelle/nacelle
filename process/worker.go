package process

import (
	"errors"
	"sync"
	"time"

	"github.com/efritz/glock"

	"github.com/efritz/nacelle"
)

type (
	Worker struct {
		Container    *nacelle.ServiceContainer `service:"container"`
		spec         WorkerSpec
		clock        glock.Clock
		halt         chan struct{}
		once         *sync.Once
		tickInterval time.Duration
	}

	WorkerSpec interface {
		Init(nacelle.Config, *Worker) error
		Tick() error
	}
)

var ErrBadWorkerConfig = errors.New("worker config not registered properly")

func NewWorker(spec WorkerSpec) *Worker {
	return newWorker(spec, glock.NewRealClock())
}

func newWorker(spec WorkerSpec, clock glock.Clock) *Worker {
	return &Worker{
		spec:  spec,
		clock: clock,
		halt:  make(chan struct{}),
		once:  &sync.Once{},
	}
}

func (w *Worker) IsDone() bool {
	select {
	case <-w.HaltChan():
		return true
	default:
		return false
	}
}

func (w *Worker) HaltChan() <-chan struct{} {
	return w.halt
}

func (w *Worker) Init(config nacelle.Config) error {
	rawConfig, err := config.Get(WorkerConfigToken)
	if err != nil {
		return err
	}

	workerConfig, ok := rawConfig.(*WorkerConfig)
	if !ok {
		return ErrBadWorkerConfig
	}

	w.tickInterval = workerConfig.WorkerTickInterval

	if err := w.Container.Inject(w.spec); err != nil {
		return err
	}

	return w.spec.Init(config, w)
}

func (w *Worker) Start() error {
	defer w.Stop()

loop:
	for {
		select {
		case <-w.halt:
			break loop
		case <-w.clock.After(w.tickInterval):
		}

		if err := w.spec.Tick(); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) Stop() (err error) {
	w.once.Do(func() { close(w.halt) })
	return
}
