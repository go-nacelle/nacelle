package worker

import (
	"fmt"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	"github.com/efritz/nacelle"
	"github.com/efritz/nacelle/process"
	"github.com/efritz/nacelle/service"
	. "github.com/onsi/gomega"
)

type WorkerSuite struct{}

func (s *WorkerSuite) TestRunAndStop(t sweet.T) {
	var (
		spec     = newMockWorkerSpec()
		clock    = glock.NewMockClock()
		worker   = makeWorker(spec, clock)
		tickChan = make(chan struct{})
		errChan  = make(chan error)
	)

	defer close(tickChan)

	spec.tick = func() error {
		tickChan <- struct{}{}
		return nil
	}

	err := worker.Init(makeConfig(&Config{RawWorkerTickInterval: 5}))
	Expect(err).To(BeNil())

	go func() {
		errChan <- worker.Start()
	}()

	Eventually(tickChan).Should(Receive())
	Consistently(tickChan).ShouldNot(Receive())
	clock.BlockingAdvance(time.Second * 5)
	Eventually(tickChan).Should(Receive())
	Consistently(tickChan).ShouldNot(Receive())
	clock.BlockingAdvance(time.Second * 5)
	Eventually(tickChan).Should(Receive())

	Expect(worker.IsDone()).To(BeFalse())
	worker.Stop()
	Expect(worker.IsDone()).To(BeTrue())
	Eventually(errChan).Should(Receive(BeNil()))
}

func (s *WorkerSuite) TestBadInject(t sweet.T) {
	worker := NewWorker(&badInjectWorkerSpec{})
	worker.Services = makeBadContainer()
	worker.Health = process.NewHealth()

	err := worker.Init(makeConfig(&Config{RawWorkerTickInterval: 5}))
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("ServiceA"))
}

func (s *WorkerSuite) TestInitError(t sweet.T) {
	var (
		spec   = newMockWorkerSpec()
		worker = makeWorker(spec, glock.NewRealClock())
	)

	spec.init = func(config nacelle.Config, worker *Worker) error {
		return fmt.Errorf("utoh")
	}

	err := worker.Init(makeConfig(&Config{RawWorkerTickInterval: 5}))
	Expect(err).To(MatchError("utoh"))
}

func (s *WorkerSuite) TestTickError(t sweet.T) {
	var (
		spec    = newMockWorkerSpec()
		clock   = glock.NewMockClock()
		worker  = makeWorker(spec, clock)
		errChan = make(chan error)
	)

	spec.tick = func() error {
		return fmt.Errorf("utoh")
	}

	err := worker.Init(makeConfig(&Config{RawWorkerTickInterval: 5}))
	Expect(err).To(BeNil())

	go func() {
		errChan <- worker.Start()
	}()

	Eventually(errChan).Should(Receive(MatchError("utoh")))
	Expect(worker.IsDone()).To(BeTrue())
}

func makeWorker(spec WorkerSpec, clock glock.Clock) *Worker {
	worker := newWorker(spec, clock)
	worker.Services, _ = service.NewContainer()
	worker.Health = process.NewHealth()
	return worker
}

//
// Mocks

type mockSpec struct {
	init func(nacelle.Config, *Worker) error
	tick func() error
}

func newMockWorkerSpec() *mockSpec {
	return &mockSpec{
		init: func(nacelle.Config, *Worker) error { return nil },
		tick: func() error { return nil },
	}
}

func (s *mockSpec) Init(c nacelle.Config, w *Worker) error { return s.init(c, w) }
func (s *mockSpec) Tick() error                            { return s.tick() }

//
// Bad Injection

type badInjectWorkerSpec struct {
	ServiceA *A `service:"A"`
}

func (s *badInjectWorkerSpec) Init(c nacelle.Config, w *Worker) error { return nil }
func (s *badInjectWorkerSpec) Tick() error                            { return nil }
