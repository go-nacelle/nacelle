package secret

import "github.com/go-nacelle/nacelle"

type logAdapter struct {
	nacelle.Logger
}

func (l *logAdapter) Printf(format string, args ...interface{}) {
	l.Logger.Debug(format, args...)
}
