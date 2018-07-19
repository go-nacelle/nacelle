package http

import "github.com/efritz/nacelle/config/tag"

type (
	options struct {
		tagModifiers []tag.Modifier
	}

	// ConfigFunc is a function used to configure an instance of a Worker.
	ConfigFunc func(*options)
)

// WithTagModifiers applies the givne tag modifiers on config load.
func WithTagModifiers(modifiers ...tag.Modifier) ConfigFunc {
	return func(o *options) { o.tagModifiers = append(o.tagModifiers, modifiers...) }
}

func getOptions(configs []ConfigFunc) *options {
	options := &options{}
	for _, f := range configs {
		f(options)
	}

	return options
}
