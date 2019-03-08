package process

import "sort"

type (
	// Container is a collection of initializers and processes.
	Container interface {
		// RegisterInitializer adds an initializer to the container
		// with the given configuration.
		RegisterInitializer(Initializer, ...InitializerConfigFunc)

		// RegisterProcess adds a process to the container with the
		// given configuration.
		RegisterProcess(Process, ...ProcessConfigFunc)

		// NumInitializers returns the number of registered initializers.
		NumInitializers() int

		// NumProcesses returns the number of registered processes.
		NumProcesses() int

		// NumPriorities returns the number of distinct registered
		// process priorities.
		NumPriorities() int

		// GetInitializers returns a slice of meta objects wrapping
		// all registered initializers.
		GetInitializers() []*InitializerMeta

		// GetProcessesAtPriorityIndex returns  aslice of meta objects
		// wrapping all processes registered to this priority index,
		// where zero denotes the lowest priority, one the second
		// lowest, and so on. The index parameter is not checked for
		// validity before indexing an internal slice - caller beware.
		GetProcessesAtPriorityIndex(index int) []*ProcessMeta
	}

	container struct {
		initializers []*InitializerMeta
		processes    map[int][]*ProcessMeta
		priorities   []int
	}
)

// NewContainer creates an empty process container.
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

func (c *container) NumInitializers() int {
	return len(c.initializers)
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
