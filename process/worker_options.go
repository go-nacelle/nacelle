package process

type (
	workerOptions struct {
		configToken interface{}
	}

	// WorkerConfigFunc is a function used to configure an instance of a Worker.
	WorkerConfigFunc func(*workerOptions)
)

// WithWorkerConfigToken sets the config token to use. This is useful if an application
// has multiple Worker processes running with different configuration tags.
func WithWorkerConfigToken(token interface{}) WorkerConfigFunc {
	return func(o *workerOptions) { o.configToken = token }
}

func getWorkerOptions(configs []WorkerConfigFunc) *workerOptions {
	options := &workerOptions{
		configToken: WorkerConfigToken,
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
