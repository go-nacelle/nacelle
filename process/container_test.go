package process

import (
	"fmt"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/nacelle/config"
)

type ContainerSuite struct{}

func (s *ContainerSuite) TestInitializers(t sweet.T) {
	i1 := InitializerFunc(func(config.Config) error { return fmt.Errorf("a") })
	i2 := InitializerFunc(func(config.Config) error { return fmt.Errorf("b") })
	i3 := InitializerFunc(func(config.Config) error { return fmt.Errorf("c") })

	c := NewContainer()
	c.RegisterInitializer(i1)
	c.RegisterInitializer(i2, WithInitializerName("b"))
	c.RegisterInitializer(i3, WithInitializerName("c"), WithInitializerTimeout(time.Minute*2))

	initializers := c.GetInitializers()
	Expect(initializers).To(HaveLen(3))

	// Test names
	Expect(initializers[0].Name()).To(Equal("<unnamed>"))
	Expect(initializers[1].Name()).To(Equal("b"))
	Expect(initializers[2].Name()).To(Equal("c"))

	// Test timeout
	Expect(initializers[0].InitTimeout()).To(Equal(time.Second * 0))
	Expect(initializers[1].InitTimeout()).To(Equal(time.Second * 0))
	Expect(initializers[2].InitTimeout()).To(Equal(time.Minute * 2))

	// Test inner function
	Expect(initializers[0].Initializer.Init(nil)).Should(MatchError("a"))
	Expect(initializers[1].Initializer.Init(nil)).Should(MatchError("b"))
	Expect(initializers[2].Initializer.Init(nil)).Should(MatchError("c"))
}

func (s *ContainerSuite) TestProcesses(t sweet.T) {
	c := NewContainer()
	c.RegisterProcess(newInitFailProcess("a"))
	c.RegisterProcess(newInitFailProcess("b"), WithProcessName("b"), WithPriority(5))
	c.RegisterProcess(newInitFailProcess("c"), WithProcessName("c"), WithPriority(2))
	c.RegisterProcess(newInitFailProcess("d"), WithProcessName("d"), WithPriority(3))
	c.RegisterProcess(newInitFailProcess("e"), WithProcessName("e"), WithPriority(2))
	c.RegisterProcess(newInitFailProcess("f"), WithProcessName("f"))

	Expect(c.NumProcesses()).To(Equal(6))
	Expect(c.NumPriorities()).To(Equal(4))

	p1 := c.GetProcessesAtPriorityIndex(0)
	p2 := c.GetProcessesAtPriorityIndex(1)
	p3 := c.GetProcessesAtPriorityIndex(2)
	p4 := c.GetProcessesAtPriorityIndex(3)

	Expect(p1).To(HaveLen(2))
	Expect(p2).To(HaveLen(2))
	Expect(p3).To(HaveLen(1))
	Expect(p4).To(HaveLen(1))

	// Test priorities
	Expect(p1[0].priority).To(Equal(0))
	Expect(p2[0].priority).To(Equal(2))
	Expect(p3[0].priority).To(Equal(3))
	Expect(p4[0].priority).To(Equal(5))

	// Test names + order
	Expect(p1[0].Name()).To(Equal("<unnamed>"))
	Expect(p1[1].Name()).To(Equal("f"))
	Expect(p2[0].Name()).To(Equal("c"))
	Expect(p2[1].Name()).To(Equal("e"))
	Expect(p3[0].Name()).To(Equal("d"))
	Expect(p4[0].Name()).To(Equal("b"))

	// Test inner function
	Expect(p1[0].Process.Init(nil)).To(MatchError("a"))
	Expect(p1[1].Process.Init(nil)).To(MatchError("f"))
	Expect(p2[0].Process.Init(nil)).To(MatchError("c"))
	Expect(p2[1].Process.Init(nil)).To(MatchError("e"))
	Expect(p3[0].Process.Init(nil)).To(MatchError("d"))
	Expect(p4[0].Process.Init(nil)).To(MatchError("b"))
}

//
//

type initFailProcess struct {
	name string
}

func newInitFailProcess(name string) Process {
	return &initFailProcess{name: name}
}

func (p *initFailProcess) Init(config config.Config) error { return fmt.Errorf("%s", p.name) }
func (p *initFailProcess) Start() error                    { return nil }
func (p *initFailProcess) Stop() error                     { return nil }
