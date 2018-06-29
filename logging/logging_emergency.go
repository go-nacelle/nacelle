package logging

func LogEmergencyError(message string, err error) {
	l := EmergencyLogger()
	l.Error(message, err.Error())
	l.Sync()
}

func LogEmergencyErrors(message string, errs []error) {
	l := EmergencyLogger()

	for _, err := range errs {
		l.Error(message, err.Error())
	}

	l.Sync()
}

func EmergencyLogger() Logger {
	logger, _ := InitGomolShim(&Config{
		LogLevel:    "DEBUG",
		LogEncoding: "json",
	})

	return logger
}
