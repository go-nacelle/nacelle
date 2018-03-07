package log

import (
	"os"
	"strings"

	"github.com/aphistic/gomol"
	console "github.com/aphistic/gomol-console"
)

type GomolShim struct {
	logger gomol.WrappableLogger
}

var gomolLevels = map[LogLevel]gomol.LogLevel{
	LevelDebug:   gomol.LevelDebug,
	LevelInfo:    gomol.LevelInfo,
	LevelWarning: gomol.LevelWarning,
	LevelError:   gomol.LevelError,
	LevelFatal:   gomol.LevelFatal,
}

//
// Shim

func NewGomolLogger(logger *gomol.LogAdapter, initialFields Fields) Logger {
	return adaptShim((&GomolShim{logger}).WithFields(initialFields))
}

func (g *GomolShim) WithFields(fields Fields) logShim {
	if len(fields) == 0 {
		return g
	}

	return &GomolShim{gomol.NewLogAdapterFor(g.logger, gomol.NewAttrsFromMap(fields))}
}

func (g *GomolShim) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	g.logger.Log(gomolLevels[level], gomol.NewAttrsFromMap(addCaller(fields).normalizeTimeValues()), format, args...)

	if level == LevelFatal {
		g.logger.ShutdownLoggers()
		os.Exit(1)
	}
}

func (g *GomolShim) Sync() error {
	return gomol.ShutdownLoggers()
}

//
// Init

func InitGomolShim(c *Config) (Logger, error) {
	level, _ := gomol.ToLogLevel(c.LogLevel)
	gomol.SetLogLevel(level)

	if c.LogEncoding == "console" {
		consoleLogger, err := console.NewConsoleLogger(&console.ConsoleLoggerConfig{
			Colorize: true,
			Writer:   os.Stderr,
		})

		if err != nil {
			return nil, err
		}

		tpl, err := newGomolConsoleTemplate(c.LogColorize)
		if err != nil {
			return nil, err
		}

		consoleLogger.SetTemplate(tpl)
		gomol.AddLogger(consoleLogger)
	} else {
		gomol.AddLogger(newJSONLogger())
	}

	if err := gomol.InitLoggers(); err != nil {
		return nil, err
	}

	return NewGomolLogger(gomol.NewLogAdapter(nil), c.LogInitialFields), nil
}

func newGomolConsoleTemplate(color bool) (*gomol.Template, error) {
	text := "" +
		`[{{.Timestamp.Format "2006-01-02 15:04:05.000"}}] ` +
		`{{color}}{{printf "%5s" (ucase .LevelName)}}{{reset}} ` +
		"{{.Message}}" +
		"{{if .Attrs}}{{range $key, $val := .Attrs}} {{color}}{{$key}}{{reset}}={{$val}}{{end}}{{end}}"

	if !color {
		text = removeColor(text)
	}

	return gomol.NewTemplate(text)
}

func removeColor(text string) string {
	return strings.NewReplacer("{{color}}", "", "{{reset}}", "").Replace(text)
}
