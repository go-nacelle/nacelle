package process

import "time"

type (
	InitializerMeta struct {
		Initializer
		name        string
		initTimeout time.Duration
	}
)

func newInitializerMeta(initializer Initializer) *InitializerMeta {
	return &InitializerMeta{
		Initializer: initializer,
	}
}

func (m *InitializerMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}

func (m *InitializerMeta) InitTimeout() time.Duration {
	return m.initTimeout
}

func (m *InitializerMeta) Wrapped() interface{} {
	return m.Initializer
}
