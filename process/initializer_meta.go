package process

import "time"

type (
	// InitializerMeta wraps an initializer with some package
	// private fields.
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

// Name returns the name of the initializer.
func (m *InitializerMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}

// InitTimeout returns the maximum timeout allowed for a call to
// the Init function. A zero value indicates no timeout.
func (m *InitializerMeta) InitTimeout() time.Duration {
	return m.initTimeout
}

// Wrapped returns the underlying initializer.
func (m *InitializerMeta) Wrapped() interface{} {
	return m.Initializer
}
