package process

import (
	"fmt"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/glock"
	"github.com/efritz/nacelle/config"
	"github.com/efritz/nacelle/service"
)

type RunnerSuite struct{}

func (s *RunnerSuite) TestRunOrder(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		finalize    = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedInitializer(init, "b")
		i3 = newTaggedFinalizer(init, finalize, "c")
		p1 = newTaggedProcess(init, start, stop, "d")
		p2 = newTaggedProcess(init, start, stop, "e")
		p3 = newTaggedProcess(init, start, stop, "f")
		p4 = newTaggedProcess(init, start, stop, "g")
		p5 = newTaggedProcess(init, start, stop, "h")
	)

	// Register things
	processes.RegisterInitializer(i1)
	processes.RegisterInitializer(i2)
	processes.RegisterInitializer(i3)
	processes.RegisterProcess(p1)
	processes.RegisterProcess(p2, WithPriority(5))
	processes.RegisterProcess(p3, WithPriority(5))
	processes.RegisterProcess(p4, WithPriority(3))
	processes.RegisterProcess(p5)

	var (
		n1, n2, n3, n4, n5 string
		errChan            = make(chan error)
		shutdownChan       = make(chan error)
	)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Initializers
	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(init).Should(Receive(Equal("c")))

	// Priority index 0
	Eventually(init).Should(Receive(Equal("d")))
	Eventually(init).Should(Receive(Equal("h")))

	// May start in either order
	Eventually(start).Should(Receive(&n1))
	Eventually(start).Should(Receive(&n2))
	Expect([]string{n1, n2}).To(ConsistOf("d", "h"))

	// Priority index 1
	Eventually(init).Should(Receive(Equal("g")))
	Eventually(start).Should(Receive(Equal("g")))

	// Priority index 2
	Eventually(init).Should(Receive(Equal("e")))
	Eventually(init).Should(Receive(Equal("f")))
	Eventually(start).Should(Receive(&n1))
	Eventually(start).Should(Receive(&n2))
	Expect([]string{n1, n2}).To(ConsistOf("e", "f"))

	go func() {
		defer close(shutdownChan)
		shutdownChan <- runner.Shutdown(time.Minute)
	}()

	// May stop in any order
	Eventually(stop).Should(Receive(&n1))
	Eventually(stop).Should(Receive(&n2))
	Eventually(stop).Should(Receive(&n3))
	Eventually(stop).Should(Receive(&n4))
	Eventually(stop).Should(Receive(&n5))
	Expect([]string{n1, n2, n3, n4, n5}).To(ConsistOf("d", "e", "f", "g", "h"))

	// Finalizers
	Eventually(finalize).Should(Receive(&n1))
	Eventually(finalize).Should(Receive(&n2))
	Expect([]string{n1, n2}).To(Equal([]string{"c", "a"}))

	// Ensure unblocked
	Eventually(shutdownChan).Should(Receive(BeNil()))
	Eventually(shutdownChan).Should(BeClosed())
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestEarlyExit(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		p1 = newTaggedProcess(init, start, stop, "a")
		p2 = newTaggedProcess(init, start, stop, "b")
	)

	// Register things
	processes.RegisterProcess(p1)
	processes.RegisterProcess(p2)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(start).Should(Receive())
	Eventually(start).Should(Receive())

	go p2.Stop()

	// Stopping one process should shutdown the rest
	Eventually(stop).Should(Receive(Equal("b")))

	var n1, n2 string
	Eventually(stop).Should(Receive(&n1))
	Eventually(stop).Should(Receive(&n2))
	Expect([]string{n1, n2}).To(ConsistOf("a", "b"))
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestSilentExit(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
	)

	var (
		p1 = newTaggedProcess(init, start, stop, "a")
		p2 = newTaggedProcess(init, start, stop, "b")
	)

	// Register things
	processes.RegisterProcess(p1)
	processes.RegisterProcess(p2, WithSilentExit())

	go runner.Run(nil)

	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(start).Should(Receive())
	Eventually(start).Should(Receive())

	go p2.Stop()

	var n1 string
	Eventually(stop).Should(Receive(&n1))
	Expect(n1).To(Equal("b"))
	Consistently(stop).ShouldNot(Receive())
}

func (s *RunnerSuite) TestShutdownTimeout(t sweet.T) {
	var (
		services, _  = service.NewContainer()
		processes    = NewContainer()
		health       = NewHealth()
		clock        = glock.NewMockClock()
		runner       = NewRunner(processes, services, health, WithClock(clock))
		sync         = make(chan struct{})
		process      = newBlockingProcess(sync)
		errChan      = make(chan error)
		shutdownChan = make(chan error)
	)

	// Register things
	processes.RegisterProcess(process)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	Eventually(sync).Should(BeClosed())

	go func() {
		defer close(shutdownChan)
		shutdownChan <- runner.Shutdown(time.Minute)
	}()

	clock.BlockingAdvance(time.Minute)
	Eventually(shutdownChan).Should(Receive(MatchError("process runner did not shutdown within timeout")))
}

func (s *RunnerSuite) TestProcessStartTimeout(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		clock       = glock.NewMockClock()
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
		runner      = NewRunner(
			processes,
			services,
			health,
			WithClock(clock),
			WithStartTimeout(time.Minute),
		)
	)

	// Stop the process from going healthy
	health.AddReason("utoh1")
	health.AddReason("utoh2")

	var (
		p1 = newTaggedProcess(init, start, stop, "a")
		p2 = newTaggedProcess(init, start, stop, "b")
	)

	processes.RegisterProcess(
		p1,
		WithProcessName("a"),
		WithProcessStartTimeout(time.Second*30),
	)

	processes.RegisterProcess(
		p2,
		WithProcessName("b"),
		WithProcessStartTimeout(time.Second*45),
	)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Don't block startup
	Eventually(init).Should(Receive())
	Eventually(init).Should(Receive())
	Eventually(start).Should(Receive())
	Eventually(start).Should(Receive())

	// Ensure timeout is respected
	Consistently(errChan).ShouldNot(Receive())
	clock.Advance(time.Second * 30)

	// Watcher should shut down
	Eventually(stop).Should(Receive())
	Eventually(stop).Should(Receive())

	var err error
	Eventually(errChan).Should(Receive(&err))
	Eventually(errChan).Should(BeClosed())

	// Check error message
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("process did not become healthy within timeout"))
	Expect(err.Error()).To(ContainSubstring("utoh1"))
	Expect(err.Error()).To(ContainSubstring("utoh2"))
}

func (s *RunnerSuite) TestProcessShutdownTimeout(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		clock       = glock.NewMockClock()
		runner      = NewRunner(processes, services, health, WithClock(clock))
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
	)

	processes.RegisterProcess(
		newTaggedProcess(init, start, stop, "a"),
		WithProcessName("a"),
		WithProcessShutdownTimeout(time.Second*10),
	)

	var (
		errChan      = make(chan error)
		shutdownChan = make(chan error)
	)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	go func() {
		// Stupid flaky goroutine scheduling
		<-time.After(time.Millisecond * 100)

		defer close(shutdownChan)
		shutdownChan <- runner.Shutdown(time.Minute)
	}()

	Eventually(init).Should(Receive())
	Eventually(stop).Should(Receive())

	// Blocked on process start method
	Consistently(shutdownChan).ShouldNot(Receive())
	clock.Advance(time.Second * 5)
	Consistently(shutdownChan).ShouldNot(Receive())
	clock.Advance(time.Second * 5)

	// Unblock after timeout
	Eventually(shutdownChan).Should(Receive(BeNil()))
	Eventually(shutdownChan).Should(BeClosed())
	Eventually(errChan).Should(Receive(MatchError("a did not shutdown within timeout")))
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestInitializerInjectionError(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedInitializer(init, "a")
		i2 = newInitializerWithService()
		i3 = newTaggedInitializer(init, "c")
		p1 = newTaggedProcess(init, start, stop, "d")
	)

	// Register things
	processes.RegisterInitializer(i1)
	processes.RegisterInitializer(i2, WithInitializerName("b"))
	processes.RegisterInitializer(i3)
	processes.RegisterProcess(p1)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Ensure error is encountered
	Eventually(init).Should(Receive(Equal("a")))
	Consistently(init).ShouldNot(Receive())

	var err error
	Eventually(errChan).Should(Receive(&err))
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("failed to inject services into b"))

	// Nothing else called
	Consistently(init).ShouldNot(Receive())
	Consistently(start).ShouldNot(Receive())
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestProcessInjectionError(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedInitializer(init, "a")
		i2 = newTaggedInitializer(init, "b")
		i3 = newTaggedInitializer(init, "c")
		p1 = newTaggedProcess(init, start, stop, "d")
		p2 = newTaggedProcess(init, start, stop, "e")
		p3 = newProcessWithService()
		p4 = newTaggedProcess(init, start, stop, "g")
		p5 = newTaggedProcess(init, start, stop, "h")
	)

	// Register things
	processes.RegisterInitializer(i1)
	processes.RegisterInitializer(i2)
	processes.RegisterInitializer(i3)
	processes.RegisterProcess(p1)
	processes.RegisterProcess(p2, WithPriority(2))
	processes.RegisterProcess(p3, WithPriority(2), WithProcessName("f"))
	processes.RegisterProcess(p4, WithPriority(2))
	processes.RegisterProcess(p5, WithPriority(3))

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Initializers
	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(init).Should(Receive(Equal("c")))

	// All processes are injected before any are initialized
	Consistently(init).ShouldNot(Receive())

	var err error
	Eventually(errChan).Should(Receive(&err))
	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("failed to inject services into f"))

	// Nothing else called
	Consistently(init).ShouldNot(Receive())
	Consistently(start).ShouldNot(Receive())
	Consistently(stop).ShouldNot(Receive())
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestInitializerInitTimeout(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		clock       = glock.NewMockClock()
		runner      = NewRunner(processes, services, health, WithClock(clock))
		init        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedInitializer(init, "a")
		i2 = newTaggedInitializer(init, "b")
	)

	// Register things
	processes.RegisterInitializer(i1, WithInitializerName("a"))
	processes.RegisterInitializer(i2, WithInitializerName("b"), WithInitializerTimeout(time.Minute))

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Don't read second value - this blocks i2.Init
	Eventually(init).Should(Receive(Equal("a")))

	// Ensure error / unblocked
	clock.BlockingAdvance(time.Minute)
	Eventually(errChan).Should(Receive(MatchError("b did not initialize within timeout")))
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestFinalizerFinalizeTimeout(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		clock       = glock.NewMockClock()
		runner      = NewRunner(processes, services, health, WithClock(clock))
		init        = make(chan string)
		finalize    = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedFinalizer(init, finalize, "b")
		p1 = newTaggedProcess(init, start, stop, "c")
	)

	// Register things
	processes.RegisterInitializer(i1, WithInitializerName("a"), WithFinalizerTimeout(time.Minute))
	processes.RegisterInitializer(i2, WithInitializerName("b"))
	processes.RegisterProcess(p1, WithProcessName("c"))

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(init).Should(Receive(Equal("c")))
	Eventually(start).Should(Receive(Equal("c")))

	// Shutdown
	go runner.Shutdown(0)
	Eventually(stop).Should(Receive())

	// Finalize first initializer
	Eventually(finalize).Should(Receive(Equal("b")))
	Consistently(errChan).ShouldNot(Receive())

	// Timeout second finalizer
	clock.BlockingAdvance(time.Minute)
	Eventually(errChan).Should(Receive(MatchError("a did not finalize within timeout")))
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestFinalizerError(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		finalize    = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedFinalizer(init, finalize, "b")
		i3 = newTaggedFinalizer(init, finalize, "c")
	)

	// Register things
	processes.RegisterInitializer(i1, WithInitializerName("a"))
	processes.RegisterInitializer(i2, WithInitializerName("b"))
	processes.RegisterInitializer(i3, WithInitializerName("c"))

	i1.finalizeErr = fmt.Errorf("utoh x")
	i2.finalizeErr = fmt.Errorf("utoh y")
	i3.finalizeErr = fmt.Errorf("utoh z")

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	for i := 0; i < 3; i++ {
		Eventually(init).Should(Receive())
	}

	for i := 0; i < 3; i++ {
		Eventually(finalize).Should(Receive())
	}

	// Stop should emit errors but continue running
	// the remaining finalizers.

	var err1, err2, err3 error
	Eventually(errChan).Should(Receive(&err1))
	Eventually(errChan).Should(Receive(&err2))
	Eventually(errChan).Should(Receive(&err3))
	Eventually(errChan).Should(BeClosed())

	Expect([]string{err1.Error(), err2.Error(), err3.Error()}).To(ConsistOf(
		"c returned error from finalize (utoh z)",
		"b returned error from finalize (utoh y)",
		"a returned error from finalize (utoh x)",
	))
}

func (s *RunnerSuite) TestProcessInitTimeout(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		clock       = glock.NewMockClock()
		runner      = NewRunner(processes, services, health, WithClock(clock))
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		p1 = newTaggedProcess(init, start, stop, "a")
		p2 = newTaggedProcess(init, start, stop, "b")
	)

	// Register things
	processes.RegisterProcess(p1, WithProcessName("a"))
	processes.RegisterProcess(p2, WithProcessName("b"), WithProcessInitTimeout(time.Minute))

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Don't read second value - this blocks i2.Init
	Eventually(init).Should(Receive(Equal("a")))

	// Ensure error / unblocked
	clock.BlockingAdvance(time.Minute)
	Eventually(errChan).Should(Receive(MatchError("b did not initialize within timeout")))
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestInitializerError(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		finalize    = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedFinalizer(init, finalize, "b")
		i3 = newTaggedInitializer(init, "c")
		p1 = newTaggedProcess(init, start, stop, "d")
	)

	i2.initErr = fmt.Errorf("utoh")

	// Register things
	processes.RegisterInitializer(i1)
	processes.RegisterInitializer(i2, WithInitializerName("b"))
	processes.RegisterInitializer(i3)
	processes.RegisterProcess(p1)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Check run order
	var n1 string
	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(finalize).Should(Receive(&n1))
	Expect(n1).To(Equal("a"))

	// Ensure error is encountered
	Eventually(errChan).Should(Receive(MatchError("failed to initialize b (utoh)")))

	// Nothing else called
	Consistently(init).ShouldNot(Receive())
	Consistently(start).ShouldNot(Receive())
	Consistently(finalize).ShouldNot(Receive())
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestProcessInitError(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedInitializer(init, "a")
		i2 = newTaggedInitializer(init, "b")
		i3 = newTaggedInitializer(init, "c")
		p1 = newTaggedProcess(init, start, stop, "d")
		p2 = newTaggedProcess(init, start, stop, "e")
		p3 = newTaggedProcess(init, start, stop, "f")
		p4 = newTaggedProcess(init, start, stop, "g")
		p5 = newTaggedProcess(init, start, stop, "h")
	)

	// Register things
	processes.RegisterInitializer(i1)
	processes.RegisterInitializer(i2)
	processes.RegisterInitializer(i3)
	processes.RegisterProcess(p1)
	processes.RegisterProcess(p2, WithPriority(2))
	processes.RegisterProcess(p3, WithPriority(2), WithProcessName("f"))
	processes.RegisterProcess(p4, WithPriority(2))
	processes.RegisterProcess(p5, WithPriority(3))

	p3.initErr = fmt.Errorf("utoh")

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Initializers
	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(init).Should(Receive(Equal("c")))

	// Lower-priority process
	Eventually(init).Should(Receive(Equal("d")))
	Eventually(start).Should(Receive(Equal("d")))

	// Ensure error is encountered
	Eventually(init).Should(Receive(Equal("e")))
	Eventually(init).Should(Receive(Equal("f")))

	var err error
	Eventually(errChan).Should(Receive(&err))
	Expect(err).To(MatchError("failed to initialize f (utoh)"))

	Consistently(init).ShouldNot(Receive())

	// Shutdown only things that started
	Eventually(stop).Should(Receive(Equal("d")))
	Consistently(stop).ShouldNot(Receive())

	// Nothing else called
	Consistently(init).ShouldNot(Receive())
	Consistently(start).ShouldNot(Receive())
	Eventually(errChan).Should(BeClosed())
}

func (s *RunnerSuite) TestProcessStartError(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
	)

	var (
		i1 = newTaggedInitializer(init, "a")
		i2 = newTaggedInitializer(init, "b")
		i3 = newTaggedInitializer(init, "c")
		p1 = newTaggedProcess(init, start, stop, "d")
		p2 = newTaggedProcess(init, start, stop, "e")
		p3 = newTaggedProcess(init, start, stop, "f")
		p4 = newTaggedProcess(init, start, stop, "g")
		p5 = newTaggedProcess(init, start, stop, "h")
	)

	// Register things
	processes.RegisterInitializer(i1)
	processes.RegisterInitializer(i2)
	processes.RegisterInitializer(i3)
	processes.RegisterProcess(p1)
	processes.RegisterProcess(p2, WithPriority(2))
	processes.RegisterProcess(p3, WithPriority(2), WithProcessName("f"))
	processes.RegisterProcess(p4, WithPriority(2))
	processes.RegisterProcess(p5, WithPriority(3), WithProcessName("h"))

	p3.startErr = fmt.Errorf("utoh")

	var (
		n1, n2, n3, n4 string
		errChan        = make(chan error)
	)

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	// Initializers
	Eventually(init).Should(Receive(Equal("a")))
	Eventually(init).Should(Receive(Equal("b")))
	Eventually(init).Should(Receive(Equal("c")))

	// Lower-priority process
	Eventually(init).Should(Receive(Equal("d")))
	Eventually(start).Should(Receive(Equal("d")))

	// Ensure error is encountered
	Eventually(init).Should(Receive(Equal("e")))
	Eventually(init).Should(Receive(Equal("f")))
	Eventually(init).Should(Receive(Equal("g")))

	Eventually(start).Should(Receive(&n1))
	Eventually(start).Should(Receive(&n2))
	Eventually(start).Should(Receive(&n3))
	Expect([]string{n1, n2, n3}).To(ConsistOf("e", "f", "g"))

	// Shutdown everything that's started
	Eventually(stop).Should(Receive(&n1))
	Eventually(stop).Should(Receive(&n2))
	Eventually(stop).Should(Receive(&n3))
	Eventually(stop).Should(Receive(&n4))
	Expect([]string{n1, n2, n3, n4}).To(ConsistOf("d", "e", "f", "g"))
	Consistently(stop).ShouldNot(Receive())

	// We get a start error from a goroutine, which means that
	// the next priority may be initializing. Since we're blocked
	// there, we should get a message saying that we're ignoring
	// the value of that process as well as the error from the
	// failing process.

	var err1, err2 error
	Eventually(errChan).Should(Receive(&err1))
	Eventually(errChan).Should(Receive(&err2))
	Eventually(errChan).Should(BeClosed())

	Expect([]string{err1.Error(), err2.Error()}).To(ConsistOf(
		"aborting initialization of h",
		"f returned a fatal error (utoh)",
	))
}

func (s *RunnerSuite) TestProcessStopError(t sweet.T) {
	var (
		services, _ = service.NewContainer()
		processes   = NewContainer()
		health      = NewHealth()
		runner      = NewRunner(processes, services, health)
		init        = make(chan string)
		start       = make(chan string)
		stop        = make(chan string)
		errChan     = make(chan error)
	)

	var (
		i1 = newTaggedInitializer(init, "a")
		i2 = newTaggedInitializer(init, "b")
		i3 = newTaggedInitializer(init, "c")
		p1 = newTaggedProcess(init, start, stop, "d")
		p2 = newTaggedProcess(init, start, stop, "e")
		p3 = newTaggedProcess(init, start, stop, "f")
		p4 = newTaggedProcess(init, start, stop, "g")
		p5 = newTaggedProcess(init, start, stop, "h")
	)

	// Register things
	processes.RegisterInitializer(i1)
	processes.RegisterInitializer(i2)
	processes.RegisterInitializer(i3)
	processes.RegisterProcess(p1, WithProcessName("d"))
	processes.RegisterProcess(p2, WithProcessName("e"), WithPriority(5))
	processes.RegisterProcess(p3, WithProcessName("f"), WithPriority(5))
	processes.RegisterProcess(p4, WithProcessName("g"), WithPriority(3))
	processes.RegisterProcess(p5, WithProcessName("h"))

	p1.stopErr = fmt.Errorf("utoh x")
	p3.stopErr = fmt.Errorf("utoh y")
	p5.stopErr = fmt.Errorf("utoh z")

	go func() {
		defer close(errChan)

		for err := range runner.Run(nil) {
			errChan <- err
		}
	}()

	for i := 0; i < 8; i++ {
		Eventually(init).Should(Receive())
	}

	for i := 0; i < 5; i++ {
		Eventually(start).Should(Receive())
	}

	// Shutdown
	go runner.Shutdown(time.Minute)

	for i := 0; i < 5; i++ {
		Eventually(stop).Should(Receive())
	}

	// Stop should emit errors but not block the progress
	// of the runner in any significant way (stop is only
	// called on shutdown, so cannot be double-fatal).

	var err1, err2, err3 error
	Eventually(errChan).Should(Receive(&err1))
	Eventually(errChan).Should(Receive(&err2))
	Eventually(errChan).Should(Receive(&err3))
	Eventually(errChan).Should(BeClosed())

	Expect([]string{err1.Error(), err2.Error(), err3.Error()}).To(ConsistOf(
		"d returned error from stop (utoh x)",
		"f returned error from stop (utoh y)",
		"h returned error from stop (utoh z)",
	))
}

//
//

type taggedInitializer struct {
	name    string
	init    chan<- string
	initErr error
}

func newTaggedInitializer(init chan<- string, name string) *taggedInitializer {
	return &taggedInitializer{
		name: name,
		init: init,
	}
}

func (i *taggedInitializer) Init(c config.Config) error {
	i.init <- i.name
	return i.initErr
}

//
//

type taggedFinalizer struct {
	taggedInitializer
	finalize    chan<- string
	finalizeErr error
}

func newTaggedFinalizer(init chan<- string, finalize chan<- string, name string) *taggedFinalizer {
	return &taggedFinalizer{
		taggedInitializer: taggedInitializer{
			name: name,
			init: init,
		},
		finalize: finalize,
	}
}

func (i *taggedFinalizer) Finalize() error {
	i.finalize <- i.name
	return i.finalizeErr
}

//
//

type taggedProcess struct {
	name  string
	init  chan<- string
	start chan<- string
	stop  chan<- string
	wait  chan struct{}

	initErr  error
	startErr error
	stopErr  error
}

func newTaggedProcess(init, start, stop chan<- string, name string) *taggedProcess {
	return &taggedProcess{
		name:  name,
		init:  init,
		start: start,
		stop:  stop,
		wait:  make(chan struct{}, 1), // Make this safe to close twice w/o blocking
	}
}

func (p *taggedProcess) Init(c config.Config) error {
	p.init <- p.name
	return p.initErr
}

func (p *taggedProcess) Start() error {
	p.start <- p.name

	if p.startErr != nil {
		return p.startErr
	}

	<-p.wait
	return nil
}

func (p *taggedProcess) Stop() error {
	p.stop <- p.name
	p.wait <- struct{}{}
	return p.stopErr
}

type blockingProcess struct {
	sync chan struct{}
	wait chan struct{}
}

func newBlockingProcess(sync chan struct{}) *blockingProcess {
	return &blockingProcess{
		sync: sync,
	}
}

func (p *blockingProcess) Init(c config.Config) error { return nil }
func (p *blockingProcess) Start() error               { close(p.sync); <-p.wait; return nil }
func (p *blockingProcess) Stop() error                { return nil }

//
//

type initializerWithService struct {
	X struct{} `service:"notset"`
}

type processWithService struct {
	X struct{} `service:"notset"`
}

func newInitializerWithService() *initializerWithService { return &initializerWithService{} }
func newProcessWithService() *processWithService         { return &processWithService{} }

func (i *initializerWithService) Init(c config.Config) error { return nil }
func (p *processWithService) Init(c config.Config) error     { return nil }
func (p *processWithService) Start() error                   { return nil }
func (p *processWithService) Stop() error                    { return nil }
