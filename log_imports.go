package nacelle

import "github.com/go-nacelle/log/v2"

type (
	LogFields    = log.LogFields
	Logger       = log.Logger
	LogLevel     = log.LogLevel
	ReplayLogger = log.ReplayLogger
)

const (
	LevelDebug   = log.LevelDebug
	LevelError   = log.LevelError
	LevelFatal   = log.LevelFatal
	LevelInfo    = log.LevelInfo
	LevelWarning = log.LevelWarning
)

var (
	EmergencyLogger    = log.EmergencyLogger
	LogEmergencyError  = log.LogEmergencyError
	LogEmergencyErrors = log.LogEmergencyErrors
	NewNilLogger       = log.NewNilLogger
	NewReplayAdapter   = log.NewReplayAdapter
	NewRollupAdapter   = log.NewRollupAdapter
)
