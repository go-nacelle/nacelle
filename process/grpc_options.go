package process

import "google.golang.org/grpc"

type (
	grpcOptions struct {
		configToken   interface{}
		serverOptions []grpc.ServerOption
	}

	// GRPCServerConfigFunc is a function used to configure an instance of
	// a GRPC Server.
	GRPCServerConfigFunc func(*grpcOptions)
)

// WithGRPCConfigToken sets the config token to use. This is useful if an application
// has multiple GRPC processes running with different configuration tags.
func WithGRPCConfigToken(token interface{}) GRPCServerConfigFunc {
	return func(o *grpcOptions) { o.configToken = token }
}

// WithGRPCServerOptions sets grpc options on the underlying server.
func WithGRPCServerOptions(options ...grpc.ServerOption) GRPCServerConfigFunc {
	return func(o *grpcOptions) { o.serverOptions = append(o.serverOptions, options...) }
}

func getGRPCOptions(configs []GRPCServerConfigFunc) *grpcOptions {
	options := &grpcOptions{
		configToken: GRPCConfigToken,
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
