package nacelle

import (
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
		// TODO - throw a once here you bozo!
	}
)

func (m *initializerMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}

func (m *processMeta) Name() string {
	if m.name == "" {
		return "<unnamed>"
	}

	return m.name
}
