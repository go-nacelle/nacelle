package config

import "github.com/efritz/nacelle/config/tag"

type (
	// Config is a structure that can populate the exported fields of a
	// struct based on tags (e.g. the `env` tag for the EnvConfig struct).
	Config interface {
		// Load populates a configuration object. The given tag modifiers
		// are applied to the configuration object pre-load. If the target
		// value conforms to the PostLoadConfig interface, the PostLoad
		// function may be called multiple times.
		Load(interface{}, ...tag.TagModifier) error

		// MustInject calls Injects and panics on error.
		MustLoad(interface{}, ...tag.TagModifier)
	}

	// PostLoadConfig is a marker interface for configuration objects
	// which should do some post-processing after being loaded. This
	// can perform additional casting (e.g. ints to time.Duration) and
	// more sophisticated validation (e.g. enum or exclusive values).
	PostLoadConfig interface {
		PostLoad() error
	}
)
