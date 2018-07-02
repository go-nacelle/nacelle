package process

import (
	"sync"
	"time"
)

type (
	// ProcessMeta wraps a process with some package private
	// fields.
	ProcessMeta struct {
		Process
		name        string
		priority    int
		silentExit  bool
		initTimeout time.Duration
		once        *sync.Once
	}
)

func newProcessMeta(process Process) *ProcessMeta {
	return &ProcessMeta{
		Process: process,
		once:    &sync.Once{},
	}
}

// Name returns the name of the process.
func (m *ProcessMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}

// InitTimeout returns the maximum timeout allowed for a call to
// the Init function. A zero value indicates no timeout.
func (m *ProcessMeta) InitTimeout() time.Duration {
	return m.initTimeout
}

// Stop wraps the underlying process's Stop method with a Once
// value in order to guarantee that the Stop method will not
// take effect multiple times.
func (m *ProcessMeta) Stop() (err error) {
	m.once.Do(func() {
		err = m.Process.Stop()
	})

	return
}

// Wrapped returns the underlying process.
func (m *ProcessMeta) Wrapped() interface{} {
	return m.Process
}
