package nacelle

import "github.com/efritz/nacelle/logging"

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

func emergencyLogger() Logger {
	logger, _ := logging.InitGomolShim(&LoggingConfig{
		LogLevel:    "DEBUG",
		LogEncoding: "json",
	})

	return logger
}
