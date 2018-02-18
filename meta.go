package nacelle

import "fmt"

type (
	initializerMeta struct {
		Initializer
		name string
	}

	processMeta struct {
		Process
		name       string
		priority   int
		silentExit bool
	}

	InitializerConfigFunc func(*initializerMeta)
	ProcessConfigFunc     func(*processMeta)
)

func (m *initializerMeta) Name() string {
	if m.name == "" {
		return "unnamed initializer"
	}

	return fmt.Sprintf("initializer %s", m.name)
}

func (m *processMeta) Name() string {
	if m.name == "" {
		return "unnamed process"
	}

	return fmt.Sprintf("process %s", m.name)
}

//
// Configuration Functinos

func WithInitializerName(name string) InitializerConfigFunc {
	return func(meta *initializerMeta) { meta.name = name }
}

func WithProcessName(name string) ProcessConfigFunc {
	return func(meta *processMeta) { meta.name = name }
}

func WithPriority(priority int) ProcessConfigFunc {
	return func(meta *processMeta) { meta.priority = priority }
}

func WithSilentExit() ProcessConfigFunc {
	return func(meta *processMeta) { meta.silentExit = true }
}
