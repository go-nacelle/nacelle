package process

import "time"

// ProcessConfigFunc is a function used to append additional metadata
// to an process during registration.
type ProcessConfigFunc func(*ProcessMeta)

// WithProcessName assigns a name to an process, visible in logs.
func WithProcessName(name string) ProcessConfigFunc {
	return func(meta *ProcessMeta) { meta.name = name }
}

// WithPriority assigns a priority to a process. A process with a lower-valued
// priority is initialized and started before a process with a higher-valued
// priority. Two processes with the same priority are started concurrently.
func WithPriority(priority int) ProcessConfigFunc {
	return func(meta *ProcessMeta) { meta.priority = priority }
}

// WithSilentExit allows a process to exit without causing the program to halt.
// The default is the opposite, where the completion of any registered process
// (even successful) causes a graceful shutdown of the other processes.
func WithSilentExit() ProcessConfigFunc {
	return func(meta *ProcessMeta) { meta.silentExit = true }
}

// WithProcessInitTimeout sets the time limit for the process's Init method.
func WithProcessInitTimeout(timeout time.Duration) ProcessConfigFunc {
	return func(meta *ProcessMeta) { meta.initTimeout = timeout }
}

// WithProcessShutdownTimeout sets the time limit for the process's Start method
// to yield after the Stop method has been called.
func WithProcessShutdownTimeout(timeout time.Duration) ProcessConfigFunc {
	return func(meta *ProcessMeta) { meta.shutdownTimeout = timeout }
}
