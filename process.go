package nacelle

type (
	Process interface {
		Init(config Config) error
		Start() error
		Stop() error
	}

	Initializer interface {
		Init(config Config) error
	}

	InitializerFunc func(config Config) error
)

func (f InitializerFunc) Init(config Config) error {
	return f(config)
}
