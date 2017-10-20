package log

type NilLogger struct{}

func (l *NilLogger) WithFields(fields Fields) Logger                                     { return l }
func (l *NilLogger) Debug(format string, args ...interface{})                            {}
func (l *NilLogger) Info(format string, args ...interface{})                             {}
func (l *NilLogger) Warning(format string, args ...interface{})                          {}
func (l *NilLogger) Error(format string, args ...interface{})                            {}
func (l *NilLogger) Fatal(format string, args ...interface{})                            {}
func (l *NilLogger) DebugWithFields(fields Fields, format string, args ...interface{})   {}
func (l *NilLogger) InfoWithFields(fields Fields, format string, args ...interface{})    {}
func (l *NilLogger) WarningWithFields(fields Fields, format string, args ...interface{}) {}
func (l *NilLogger) ErrorWithFields(fields Fields, format string, args ...interface{})   {}
func (l *NilLogger) FatalWithFields(fields Fields, format string, args ...interface{})   {}
func (l *NilLogger) Sync() error                                                         { return nil }
