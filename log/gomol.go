package log

import (
	"fmt"
	"os"
	"strings"
	"text/template"

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
	level, err := gomol.ToLogLevel(c.LogLevel)
	if err != nil {
		return nil, err
	}

	if err := gomol.ClearLoggers(); err != nil {
		return nil, err
	}

	gomol.SetLogLevel(level)
	gomol.SetConfig(&gomol.Config{SequenceAttr: "sequence_number"})

	switch c.LogEncoding {
	case "console":
		if err := setupConsoleLogger(c); err != nil {
			return nil, err
		}

	case "json":
		setupJSONLogger()
	}

	if err := gomol.InitLoggers(); err != nil {
		return nil, err
	}

	return NewGomolLogger(gomol.NewLogAdapter(nil), c.LogInitialFields), nil
}

func setupConsoleLogger(c *Config) error {
	consoleLogger, err := console.NewConsoleLogger(&console.ConsoleLoggerConfig{
		Colorize: true,
		Writer:   os.Stderr,
	})

	if err != nil {
		return err
	}

	tpl, err := newGomolConsoleTemplate(
		c.LogColorize,
		c.LogShortTime,
		c.LogMultilineFields,
		c.LogAttrBlacklist,
	)

	if err != nil {
		return err
	}

	consoleLogger.SetTemplate(tpl)
	gomol.AddLogger(consoleLogger)
	return nil
}

func setupJSONLogger() {
	gomol.AddLogger(newJSONLogger())
}

func newGomolConsoleTemplate(color, shortTime, multilineFields bool, blacklist []string) (*gomol.Template, error) {
	var (
		attrPrefix  = " "
		attrPadding = ""
		attrSuffix  = ""
	)

	if multilineFields {
		attrPrefix = "\n    "
		attrPadding = " "
		attrSuffix = "\n"
	}

	fieldsTemplate := fmt.Sprintf(
		""+
			`{{if .Attrs}}`+
			`{{range $key, $val := .Attrs}}`+
			`{{if shouldDisplayAttr $key}}`+
			`%s{{$key}}%s=%s{{$val}}`+
			`{{end}}`+
			`{{end}}`+
			`%s`+
			`{{end}}`,
		attrPrefix,
		attrPadding,
		attrPadding,
		attrSuffix,
	)

	timeFormat := "2006/01/02 15:04:05.000"
	if shortTime {
		timeFormat = "15:04:05"
	}

	text :=
		"" +
			`{{color}}` +
			`[{{ucase .LevelName | printf "%1.1s"}}] ` +
			fmt.Sprintf(`[{{.Timestamp.Format "%s"}}] {{.Message}}`, timeFormat) +
			`{{reset}}` +
			fieldsTemplate

	if !color {
		text = removeColor(text)
	}

	return gomol.NewTemplateWithFuncMap(text, template.FuncMap{
		"shouldDisplayAttr": shouldDisplayAttr(blacklist),
	})
}

func removeColor(text string) string {
	return strings.NewReplacer("{{color}}", "", "{{reset}}", "").Replace(text)
}

func shouldDisplayAttr(blacklist []string) func(string) bool {
	return func(attr string) bool {
		for _, cmp := range blacklist {
			if cmp == attr {
				return false
			}
		}

		return true
	}
}
