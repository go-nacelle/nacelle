package process

type (
	Process interface {
		Initializer
		Start() error
		Stop() error
	}
)
