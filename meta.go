package nacelle

import (
	"fmt"
	"time"
)

type (
	initializerMeta struct {
		Initializer
		name    string
		timeout time.Duration
	}

	processMeta struct {
		Process
		name        string
		priority    int
		silentExit  bool
		initTimeout time.Duration
	}

	// InitializerConfigFunc is a function used to append additional
	// metadata to an initializer during registration.
	InitializerConfigFunc func(*initializerMeta)

	// ProcessConfigFunc is a function used to append additional metadata
	// to an process during registration.
	ProcessConfigFunc func(*processMeta)
)

func (m *initializerMeta) Name() string {
	if m.name == "" {
		return "unnamed initializer"
	}

	return fmt.Sprintf("initializer %s", m.name)
}

func (m *processMeta) Name() string {
	if m.name == "" {
		return "unnamed process"
	}

	return fmt.Sprintf("process %s", m.name)
}

//
// Configuration Functions

// WithInitializerName assigns a name to an initializer, visible in logs.
func WithInitializerName(name string) InitializerConfigFunc {
	return func(meta *initializerMeta) { meta.name = name }
}

// WithProcessName assigns a name to an process, visible in logs.
func WithProcessName(name string) ProcessConfigFunc {
	return func(meta *processMeta) { meta.name = name }
}

// WithPriority assigns a priority to a process. A process with a lower-valued
// priority is initialized and started before a process with a higher-valued
// priority. Two processes with the same priority are started concurrently.
func WithPriority(priority int) ProcessConfigFunc {
	return func(meta *processMeta) { meta.priority = priority }
}

// WithSilentExit allows a process to exit without causing the progrma to halt.
// The default is the opposite, where the completion of any registered process
// (even successful) causes a graceful shutdown of the other processes.
func WithSilentExit() ProcessConfigFunc {
	return func(meta *processMeta) { meta.silentExit = true }
}

// WithInitializerTimeout sets the time limit for the initializer.
func WithInitializerTimeout(timeout time.Duration) InitializerConfigFunc {
	return func(meta *initializerMeta) { meta.timeout = timeout }
}

// WithProcessInitTimeout sets the time limit for the process's Init method.
func WithProcessInitTimeout(timeout time.Duration) ProcessConfigFunc {
	return func(meta *processMeta) { meta.initTimeout = timeout }
}
