package process

import "time"

type (
	named interface {
		Name() string
	}

	namedInjectable interface {
		named
		Wrapped() interface{}
	}

	errMeta struct {
		err        error
		source     namedInitializer
		silentExit bool
	}

	namedInitializer interface {
		named
		Initializer
		InitTimeout() time.Duration
	}
)
