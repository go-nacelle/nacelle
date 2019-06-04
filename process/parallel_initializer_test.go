package process

import (
	"fmt"

	"github.com/aphistic/sweet"
	"github.com/go-nacelle/nacelle/logging"
	"github.com/go-nacelle/nacelle/service"
	. "github.com/onsi/gomega"
)

type ParallelInitializerSuite struct{}

func (s *ParallelInitializerSuite) TestInitialize(t sweet.T) {
	var (
		container, _ = service.NewContainer()
		init         = make(chan string, 3)
		finalize     = make(chan string, 3)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedInitializer(init, "b")
		i3 = newTaggedFinalizer(init, finalize, "c")
	)

	pi := NewParallelInitializer(
		WithParallelInitializerContainer(container),
		WithParallelInitializerLogger(logging.NewNilLogger()),
	)

	// Register things
	pi.RegisterInitializer(i1)
	pi.RegisterInitializer(i2)
	pi.RegisterInitializer(i3)

	err := pi.Init(nil)
	Expect(err).To(BeNil())

	// May initialize in any order
	var n1, n2, n3 string
	Eventually(init).Should(Receive(&n1))
	Eventually(init).Should(Receive(&n2))
	Eventually(init).Should(Receive(&n3))
	Expect([]string{n1, n2, n3}).To(ConsistOf("a", "b", "c"))
}

func (s *ParallelInitializerSuite) TestInitError(t sweet.T) {
	var (
		container, _ = service.NewContainer()
		init         = make(chan string, 4)
		finalize     = make(chan string, 4)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedFinalizer(init, finalize, "b")
		i3 = newTaggedFinalizer(init, finalize, "c")
		i4 = newTaggedFinalizer(init, finalize, "d")
		m1 = newInitializerMeta(i1)
		m2 = newInitializerMeta(i2)
		m3 = newInitializerMeta(i3)
		m4 = newInitializerMeta(i4)
	)

	pi := NewParallelInitializer(
		WithParallelInitializerContainer(container),
		WithParallelInitializerLogger(logging.NewNilLogger()),
	)

	// Register things
	pi.RegisterInitializer(i1, WithInitializerName("a"))
	pi.RegisterInitializer(i2, WithInitializerName("b"))
	pi.RegisterInitializer(i3, WithInitializerName("c"))
	pi.RegisterInitializer(i4, WithInitializerName("d"))

	i2.initErr = fmt.Errorf("utoh y")
	i3.initErr = fmt.Errorf("utoh z")
	i4.finalizeErr = fmt.Errorf("utoh w")

	WithInitializerName("a")(m1)
	WithInitializerName("b")(m2)
	WithInitializerName("c")(m3)
	WithInitializerName("d")(m4)

	err := pi.Init(nil)
	Expect(err).To(ConsistOf(
		errMeta{err: fmt.Errorf("failed to initialize b (utoh y)"), source: m2},
		errMeta{err: fmt.Errorf("failed to initialize c (utoh z)"), source: m3},
		errMeta{err: fmt.Errorf("d returned error from finalize (utoh w)"), source: m4},
	))

	var n1, n2, n3, n4 string
	Eventually(init).Should(Receive(&n1))
	Eventually(init).Should(Receive(&n2))
	Eventually(init).Should(Receive(&n3))
	Eventually(init).Should(Receive(&n4))
	Expect([]string{n1, n2, n3, n4}).To(ConsistOf("a", "b", "c", "d"))

	var n5, n6 string
	Eventually(finalize).Should(Receive(&n5))
	Eventually(finalize).Should(Receive(&n6))
	Expect([]string{n5, n6}).To(ConsistOf("a", "d"))
}

func (s *ParallelInitializerSuite) TestFinalize(t sweet.T) {
	var (
		container, _ = service.NewContainer()
		init         = make(chan string, 3)
		finalize     = make(chan string, 3)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedInitializer(init, "b")
		i3 = newTaggedFinalizer(init, finalize, "c")
	)

	pi := NewParallelInitializer(
		WithParallelInitializerContainer(container),
		WithParallelInitializerLogger(logging.NewNilLogger()),
	)

	// Register things
	pi.RegisterInitializer(i1)
	pi.RegisterInitializer(i2)
	pi.RegisterInitializer(i3)

	err := pi.Finalize()
	Expect(err).To(BeNil())

	// Should finalize in any order
	var n1, n2 string
	Eventually(finalize).Should(Receive(&n1))
	Eventually(finalize).Should(Receive(&n2))
	Expect([]string{n1, n2}).To(ConsistOf("a", "c"))
}

func (s *ParallelInitializerSuite) TestFinalizeError(t sweet.T) {
	var (
		container, _ = service.NewContainer()
		init         = make(chan string, 3)
		finalize     = make(chan string, 3)
	)

	var (
		i1 = newTaggedFinalizer(init, finalize, "a")
		i2 = newTaggedFinalizer(init, finalize, "b")
		i3 = newTaggedFinalizer(init, finalize, "c")
		m1 = newInitializerMeta(i1)
		m2 = newInitializerMeta(i2)
		m3 = newInitializerMeta(i3)
	)

	pi := NewParallelInitializer(
		WithParallelInitializerContainer(container),
		WithParallelInitializerLogger(logging.NewNilLogger()),
	)

	// Register things
	pi.RegisterInitializer(i1, WithInitializerName("a"))
	pi.RegisterInitializer(i2, WithInitializerName("b"))
	pi.RegisterInitializer(i3, WithInitializerName("c"))

	i1.finalizeErr = fmt.Errorf("utoh x")
	i2.finalizeErr = fmt.Errorf("utoh y")
	i3.finalizeErr = fmt.Errorf("utoh z")

	WithInitializerName("a")(m1)
	WithInitializerName("b")(m2)
	WithInitializerName("c")(m3)

	err := pi.Finalize()
	Expect(err).To(ConsistOf(
		errMeta{err: fmt.Errorf("a returned error from finalize (utoh x)"), source: m1},
		errMeta{err: fmt.Errorf("b returned error from finalize (utoh y)"), source: m2},
		errMeta{err: fmt.Errorf("c returned error from finalize (utoh z)"), source: m3},
	))

	var n1, n2, n3 string
	Eventually(finalize).Should(Receive(&n1))
	Eventually(finalize).Should(Receive(&n2))
	Eventually(finalize).Should(Receive(&n3))
	Expect([]string{n1, n2, n3}).To(ConsistOf("a", "b", "c"))
}
