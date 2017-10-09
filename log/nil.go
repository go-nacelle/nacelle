package log

type NilLogger struct{}

func (l *NilLogger) WithFields(fields Fields) Logger                           { return l }
func (l *NilLogger) Debug(fields Fields, format string, args ...interface{})   {}
func (l *NilLogger) Info(fields Fields, format string, args ...interface{})    {}
func (l *NilLogger) Warning(fields Fields, format string, args ...interface{}) {}
func (l *NilLogger) Error(fields Fields, format string, args ...interface{})   {}
func (l *NilLogger) Fatal(fields Fields, format string, args ...interface{})   {}
func (l *NilLogger) Sync() error                                               { return nil }
