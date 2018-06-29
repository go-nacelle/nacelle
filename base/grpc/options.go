package grpc

import "google.golang.org/grpc"

type (
	options struct {
		configToken   interface{}
		serverOptions []grpc.ServerOption
	}

	// ConfigFunc is a function used to configure an instance of
	// a GRPC Server.
	ConfigFunc func(*options)
)

// WithConfigToken sets the config token to use. This is useful if an application
// has multiple GRPC processes running with different configuration tags.
func WithConfigToken(token interface{}) ConfigFunc {
	return func(o *options) { o.configToken = token }
}

// WithServerOptions sets grpc options on the underlying server.
func WithServerOptions(opts ...grpc.ServerOption) ConfigFunc {
	return func(o *options) { o.serverOptions = append(o.serverOptions, opts...) }
}

func getOptions(configs []ConfigFunc) *options {
	options := &options{
		configToken: ConfigToken,
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
