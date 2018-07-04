package nacelle

import (
	"github.com/efritz/nacelle/logging"
)

type (
	Logger    = logging.Logger
	LogLevel  = logging.LogLevel
	LogFields = logging.Fields
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
