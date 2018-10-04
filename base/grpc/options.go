package grpc

import (
	"github.com/efritz/nacelle/config"
	"google.golang.org/grpc"
)

type (
	options struct {
		tagModifiers  []config.TagModifier
		serverOptions []grpc.ServerOption
	}

	// ConfigFunc is a function used to configure an instance of
	// a GRPC Server.
	ConfigFunc func(*options)
)

// WithTagModifiers applies the given tag modifiers on config load.
func WithTagModifiers(modifiers ...config.TagModifier) ConfigFunc {
	return func(o *options) { o.tagModifiers = append(o.tagModifiers, modifiers...) }
}

// WithServerOptions sets grpc options on the underlying server.
func WithServerOptions(opts ...grpc.ServerOption) ConfigFunc {
	return func(o *options) { o.serverOptions = append(o.serverOptions, opts...) }
}

func getOptions(configs []ConfigFunc) *options {
	options := &options{}
	for _, f := range configs {
		f(options)
	}

	return options
}
