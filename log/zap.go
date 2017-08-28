package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapShim struct {
	logger *zap.SugaredLogger
}

func NewZapShim(c *Config) (Logger, error) {
	var (
		level        zap.AtomicLevel
		levelEncoder zapcore.LevelEncoder
		timeEncoder  zapcore.TimeEncoder
	)

	if err := level.UnmarshalText([]byte(c.LogLevel)); err != nil {
		return nil, err
	}

	if c.LogEncoding == "console" {
		levelEncoder = zapcore.CapitalColorLevelEncoder
		timeEncoder = zapConsoleTimeEncoder
	} else {
		levelEncoder = zapcore.LowercaseLevelEncoder
		timeEncoder = zapJSONTimeEncoder
	}

	config := zap.Config{
		Level:             level,
		DisableCaller:     c.LogDisableCaller,
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

	sugaredLogger := logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	return (&ZapShim{logger: sugaredLogger}).WithFields(c.LogInitialFields), nil
}

func (z *ZapShim) WithFields(fields Fields) Logger {
	if len(fields) == 0 {
		return z
	}

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
	z.getLogger(fields).Fatalf(format, args...)
}

func (z *ZapShim) Sync() error {
	return z.logger.Sync()
}

func (z *ZapShim) getLogger(fields Fields) *zap.SugaredLogger {
	if len(fields) == 0 {
		return z.logger
	}

	return z.logger.With(flatten(fields)...)
}

func flatten(fields Fields) []interface{} {
	flattened := []interface{}{}
	for key, value := range fields {
		flattened = append(flattened, key)
		flattened = append(flattened, value)
	}

	return flattened
}

func zapConsoleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(ConsoleTimeFormat))
}

func zapJSONTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(JSONTimeFormat))
}
