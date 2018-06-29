package http

type (
	options struct {
		configToken interface{}
	}

	// ConfigFunc is a function used to configure an instance of
	// an HTTP Server.
	ConfigFunc func(*options)
)

// WithConfigToken sets the config token to use. This is useful if an application
// has multiple HTTP processes running with different configuration tags.
func WithConfigToken(token interface{}) ConfigFunc {
	return func(o *options) { o.configToken = token }
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
