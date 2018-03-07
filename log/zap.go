package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapShim struct {
	logger *zap.SugaredLogger
}

//
// Shim

func NewZapLogger(logger *zap.SugaredLogger, initialFields Fields) Logger {
	return adaptShim((&ZapShim{logger}).WithFields(initialFields))
}

func (z *ZapShim) WithFields(fields Fields) logShim {
	if len(fields) == 0 {
		return z
	}

	return &ZapShim{z.getLogger(fields)}
}

func (z *ZapShim) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	logger := z.getLogger(addCaller(fields))

	switch level {
	case LevelDebug:
		logger.Debugf(format, args...)
	case LevelInfo:
		logger.Infof(format, args...)
	case LevelWarning:
		logger.Warnf(format, args...)
	case LevelError:
		logger.Errorf(format, args...)
	case LevelFatal:
		logger.Fatalf(format, args...)
	}
}

func (z *ZapShim) Sync() error {
	return z.logger.Sync()
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

//
// Init

func InitZapShim(c *Config) (Logger, error) {
	var (
		level        zap.AtomicLevel
		levelEncoder zapcore.LevelEncoder
		timeEncoder  zapcore.TimeEncoder
	)

	if err := level.UnmarshalText([]byte(c.LogLevel)); err != nil {
		return nil, err
	}

	if c.LogEncoding == "console" {
		if c.LogColorize {
			levelEncoder = zapcore.CapitalColorLevelEncoder
		} else {
			levelEncoder = zapcore.CapitalLevelEncoder
		}

		timeEncoder = zapConsoleTimeEncoder
	} else {
		levelEncoder = zapcore.LowercaseLevelEncoder
		timeEncoder = zapJSONTimeEncoder
	}

	config := zap.Config{
		Level:             level,
		DisableCaller:     true,
		Encoding:          c.LogEncoding,
		Development:       false,
		DisableStacktrace: true,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			CallerKey:      "caller",
			EncodeLevel:    levelEncoder,
			EncodeTime:     timeEncoder,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return NewZapLogger(logger.Sugar(), c.LogInitialFields), nil
}

func zapConsoleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(ConsoleTimeFormat))
}

func zapJSONTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(JSONTimeFormat))
}
