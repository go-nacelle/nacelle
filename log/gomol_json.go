package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aphistic/gomol"
)

type jsonLogger struct {
	stream        io.Writer
	base          *gomol.Base
	isInitialized bool
}

func newJSONLogger() *jsonLogger {
	return &jsonLogger{
		stream: os.Stderr,
	}
}

func (l *jsonLogger) SetBase(base *gomol.Base) {
	l.base = base
}

func (l *jsonLogger) InitLogger() error {
	l.isInitialized = true
	return nil
}

func (l *jsonLogger) IsInitialized() bool {
	return l.isInitialized
}

func (l *jsonLogger) ShutdownLogger() error {
	l.isInitialized = false
	return nil
}

func (l *jsonLogger) Logm(timestamp time.Time, level gomol.LogLevel, attrs map[string]interface{}, msg string) error {
	mergedAttrs := map[string]interface{}{}

	if l.base != nil && l.base.BaseAttrs != nil {
		for key, val := range l.base.BaseAttrs.Attrs() {
			mergedAttrs[key] = val
		}
	}

	for key, val := range attrs {
		mergedAttrs[key] = val
	}

	mergedAttrs["message"] = msg
	mergedAttrs["timestamp"] = timestamp.Format(JSONTimeFormat)
	mergedAttrs["level"] = level.String()

	out, err := json.Marshal(mergedAttrs)
	if err != nil {
		return err
	}

	fmt.Fprint(l.stream, string(out)+"\n")
	return nil
}
