package log

type NilShim struct{}

func NewNilLogger() Logger {
	return adaptShim(&NilShim{})
}

func (n *NilShim) WithFields(Fields) logShim                              { return n }
func (n *NilShim) LogWithFields(LogLevel, Fields, string, ...interface{}) {}
func (n *NilShim) Sync() error                                            { return nil }
