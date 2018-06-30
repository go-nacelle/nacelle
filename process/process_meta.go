package process

import (
	"sync"
	"time"
)

type (
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

func (m *ProcessMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}

func (m *ProcessMeta) InitTimeout() time.Duration {
	return m.initTimeout
}

func (m *ProcessMeta) Stop() (err error) {
	m.once.Do(func() {
		err = m.Process.Stop()
	})

	return
}

func (m *ProcessMeta) Wrapped() interface{} {
	return m.Process
}
