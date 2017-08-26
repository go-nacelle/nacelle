package log

import (
	"fmt"
	"runtime"
	"strings"
)

func getCaller() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s:%d", trimPath(file), line)
}

func trimPath(path string) string {
	//
	// For details, see https://github.com/uber-go/zap/blob/e15639dab1b6ca5a651fe7ebfd8d682683b7d6a8/zapcore/entry.go#L101

	if idx := strings.LastIndexByte(path, '/'); idx >= 0 {
		if idx := strings.LastIndexByte(path[:idx], '/'); idx >= 0 {
			// Keep everything after the penultimate separator.
			return path[idx+1:]
		}
	}

	return path
}
