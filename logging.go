package nacelle

import (
	"github.com/efritz/nacelle/logging"
)

type (
	Logger        = logging.Logger
	Fields        = logging.Fields
	LogLevel      = logging.LogLevel
	LoggingConfig = logging.Config
)

const (
	LevelFatal   = logging.LevelFatal
	LevelError   = logging.LevelError
	LevelWarning = logging.LevelWarning
	LevelInfo    = logging.LevelInfo
	LevelDebug   = logging.LevelDebug
)

var (
	NewNilLogger       = logging.NewNilLogger
	NewReplayAdapter   = logging.NewReplayAdapter
	NewRollupAdapter   = logging.NewRollupAdapter
	LogEmergencyError  = logging.LogEmergencyError
	LogEmergencyErrors = logging.LogEmergencyErrors
	EmergencyLogger    = logging.EmergencyLogger
)
