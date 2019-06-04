package http

import "github.com/go-nacelle/nacelle/config"

type (
	options struct {
		tagModifiers []config.TagModifier
	}

	// ConfigFunc is a function used to configure an instance of a Worker.
	ConfigFunc func(*options)
)

// WithTagModifiers applies the given tag modifiers on config load.
func WithTagModifiers(modifiers ...config.TagModifier) ConfigFunc {
	return func(o *options) { o.tagModifiers = append(o.tagModifiers, modifiers...) }
}

func getOptions(configs []ConfigFunc) *options {
	options := &options{}
	for _, f := range configs {
		f(options)
	}

	return options
}
