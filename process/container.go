package process

import "sort"

type (
	Container interface {
		RegisterInitializer(Initializer, ...InitializerConfigFunc)
		RegisterProcess(Process, ...ProcessConfigFunc)
		NumProcesses() int
		NumPriorities() int
		GetInitializers() []*InitializerMeta
		GetProcessesAtPriorityIndex(index int) []*ProcessMeta
	}

	container struct {
		initializers []*InitializerMeta
		processes    map[int][]*ProcessMeta
		priorities   []int
	}
)

func NewContainer() Container {
	return &container{
		initializers: []*InitializerMeta{},
		processes:    map[int][]*ProcessMeta{},
		priorities:   []int{},
	}
}

func (c *container) RegisterInitializer(
	initializer Initializer,
	initializerConfigs ...InitializerConfigFunc,
) {
	meta := newInitializerMeta(initializer)

	for _, f := range initializerConfigs {
		f(meta)
	}

	c.initializers = append(c.initializers, meta)
}

func (c *container) RegisterProcess(
	process Process,
	processConfigs ...ProcessConfigFunc,
) {
	meta := newProcessMeta(process)

	for _, f := range processConfigs {
		f(meta)
	}

	c.processes[meta.priority] = append(c.processes[meta.priority], meta)
	c.priorities = c.getPriorities()
}

func (c *container) NumProcesses() int {
	n := 0
	for _, ps := range c.processes {
		n += len(ps)
	}

	return n
}

func (c *container) NumPriorities() int {
	return len(c.priorities)
}

func (c *container) GetInitializers() []*InitializerMeta {
	return c.initializers
}

func (c *container) GetProcessesAtPriorityIndex(index int) []*ProcessMeta {
	return c.processes[c.priorities[index]]
}

func (c *container) getPriorities() []int {
	priorities := []int{}
	for priority := range c.processes {
		priorities = append(priorities, priority)
	}

	sort.Ints(priorities)
	return priorities
}
