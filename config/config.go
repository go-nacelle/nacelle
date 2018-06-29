package config

import (
	"fmt"
)

type (
	// Config is a structure that maintains chunks of an application's
	// configuration values accessible via arbitrary key. This design is
	// meant for a library-driven architecture so that multiple pieces
	// of an application can register their own configuration requirements
	// independently from the core.
	Config interface {
		// Load attempts to read an external environment for values
		// and modifies the config's internal state on success.
		Load() []error

		// Register associates an empty configuration object. This key should
		// be unique to the application as it is generally error to register
		// the same key twice.
		Register(key interface{}, config interface{}) error

		// MustRegister calls Register and panics on error.
		MustRegister(key interface{}, config interface{})

		// Get retrieves a configuration object by its key. It is an error
		// to request a non-registered key or to Get before a call to Load.
		Get(key interface{}) (interface{}, error)

		// MustGet calls Get and panics on error.
		MustGet(key interface{}) interface{}

		// Fetch retrieve a configuration object by its key and copies the
		// values of fields into the target value. The same error conditions
		// of Get apply here. An error is also returned if the type of the
		// target value and the type of the registered config do not match.
		// If the target value conforms to the PostLoadConfig interface, the
		// PostLoad function may be called multiple times.
		Fetch(key interface{}, target interface{}) error

		// MustFetch calls Fetch and panics on error.
		MustFetch(key interface{}, target interface{})

		// ToMap will convert the configuration values into a printable
		// or loggable map.
		ToMap() (map[string]interface{}, error)
	}

	// PostLoadConfig is a marker interface for configuration objects
	// which should do some post-processing after being loaded. This
	// can perform additional casting (e.g. ints to time.Duration) and
	// more sophisticated validation (e.g. enum or exclusive values).
	PostLoadConfig interface {
		PostLoad() error
	}
)

var (
	// ErrAlreadyLoaded is returned on a second call to Config#Load.
	ErrAlreadyLoaded = fmt.Errorf("config already loaded")

	// ErrNotLoaded is returned on a call to Get without first calling Load.
	ErrNotLoaded = fmt.Errorf("config not loaded")
)
