package nacelle

import "sync"

type logShim struct {
	sync.RWMutex
	logger Logger
}

func (s *logShim) setLogger(logger Logger) {
	s.Lock()
	defer s.Unlock()

	s.logger = logger
}

func (s *logShim) Printf(format string, args ...interface{}) {
	s.RLock()
	defer s.RUnlock()

	if s.logger != nil {
		s.logger.Info(format, args...)
	}
}
