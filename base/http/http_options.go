package http

type (
	httpOptions struct {
		configToken interface{}
	}

	// HTTPServerConfigFunc is a function used to configure an instance of
	// an HTTP Server.
	HTTPServerConfigFunc func(*httpOptions)
)

// WithHTTPConfigToken sets the config token to use. This is useful if an application
// has multiple HTTP processes running with different configuration tags.
func WithHTTPConfigToken(token interface{}) HTTPServerConfigFunc {
	return func(o *httpOptions) { o.configToken = token }
}

func getHTTPOptions(configs []HTTPServerConfigFunc) *httpOptions {
	options := &httpOptions{
		configToken: HTTPConfigToken,
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
