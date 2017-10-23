package discovery

import "github.com/efritz/nacelle"

type logAdapter struct {
	logger nacelle.Logger
}

func (l *logAdapter) Printf(format string, args ...interface{}) {
	// TODO - these kinds of adapters have useless caller attributes
	l.logger.Debug(format, args...)
}
