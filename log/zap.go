package log

import "go.uber.org/zap"

type ZapShim struct {
	logger *zap.SugaredLogger
}

func (z *ZapShim) WithFields(fields Fields) Logger {
	return &ZapShim{
		logger: z.getLogger(fields),
	}
}

func (z *ZapShim) Debug(fields Fields, format string, args ...interface{}) {
	z.getLogger(fields).Debugf(format, args...)
}

func (z *ZapShim) Info(fields Fields, format string, args ...interface{}) {
	z.getLogger(fields).Infof(format, args...)
}

func (z *ZapShim) Warning(fields Fields, format string, args ...interface{}) {
	z.getLogger(fields).Warnf(format, args...)
}

func (z *ZapShim) Error(fields Fields, format string, args ...interface{}) {
	z.getLogger(fields).Errorf(format, args...)
}

func (z *ZapShim) Fatal(fields Fields, format string, args ...interface{}) {
	// TODO - die
	z.getLogger(fields).Fatalf(format, args...)
}

func (z *ZapShim) getLogger(fields Fields) *zap.SugaredLogger {
	if len(fields) == 0 {
		return z.logger
	}

	flattened := []interface{}{}

	for key, value := range fields {
		flattened = append(flattened, key)
		flattened = append(flattened, value)
	}

	return z.logger.With(flattened...)
}
