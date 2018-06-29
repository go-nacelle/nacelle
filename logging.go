package nacelle

import (
	"errors"

	"github.com/efritz/nacelle/log"
)

type (
	Logger        = log.Logger
	ReplayLogger  = log.ReplayLogger
	Fields        = log.Fields
	LoggingConfig = log.Config
	LogLevel      = log.LogLevel

	loggingConfigToken string
)

const (
	LevelFatal   = log.LevelFatal
	LevelError   = log.LevelError
	LevelWarning = log.LevelWarning
	LevelInfo    = log.LevelInfo
	LevelDebug   = log.LevelDebug
)

var (
	NewReplayAdapter = log.NewReplayAdapter
	NewRollupAdapter = log.NewRollupAdapter

	LoggingConfigToken = loggingConfigToken("nacelle-logging")
	ErrBadConfig       = errors.New("logging config not registered properly")
)

func InitLogging(config Config) (Logger, error) {
	c := &LoggingConfig{}
	if err := config.Fetch(LoggingConfigToken, c); err != nil {
		return nil, ErrBadConfig
	}

	return log.InitGomolShim(c)
}

func emergencyLogger() Logger {
	logger, _ := log.InitGomolShim(&LoggingConfig{
		LogLevel:    "DEBUG",
		LogEncoding: "json",
	})

	return logger
}

func logEmergencyError(message string, err error) {
	l := emergencyLogger()
	l.Error(message, err.Error())
	l.Sync()
}

func logEmergencyErrors(message string, errs []error) {
	l := emergencyLogger()

	for _, err := range errs {
		l.Error(message, err.Error())
	}

	l.Sync()
}

func noopLogger(fields log.Fields, message string, args ...interface{}) {
	// Silence is golden
}
