package nacelle

import (
	"github.com/efritz/nacelle/logging"
)

type (
	Logger        = logging.Logger
	ReplayLogger  = logging.ReplayLogger
	Fields        = logging.Fields
	LoggingConfig = logging.Config
	LogLevel      = logging.LogLevel
)

const (
	LevelFatal   = logging.LevelFatal
	LevelError   = logging.LevelError
	LevelWarning = logging.LevelWarning
	LevelInfo    = logging.LevelInfo
	LevelDebug   = logging.LevelDebug
)

var (
	NewNilLogger     = logging.NewNilLogger
	NewReplayAdapter = logging.NewReplayAdapter
	NewRollupAdapter = logging.NewRollupAdapter
)
