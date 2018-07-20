// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/nacelle/logging -i Logger -o mock_logger_test.go -f

package config

import (
	logging "github.com/efritz/nacelle/logging"
	sync "sync"
)

type MockLogger struct {
	statsDebugWithFieldsLock          sync.RWMutex
	statDebugWithFieldsFuncCallCount  int
	statDebugWithFieldsFuncCallParams []LoggerDebugWithFieldsParamSet
	DebugWithFieldsFunc               func(logging.Fields, string, ...interface{})

	statsSyncLock          sync.RWMutex
	statSyncFuncCallCount  int
	statSyncFuncCallParams []LoggerSyncParamSet
	SyncFunc               func() error

	statsFatalLock          sync.RWMutex
	statFatalFuncCallCount  int
	statFatalFuncCallParams []LoggerFatalParamSet
	FatalFunc               func(string, ...interface{})

	statsFatalWithFieldsLock          sync.RWMutex
	statFatalWithFieldsFuncCallCount  int
	statFatalWithFieldsFuncCallParams []LoggerFatalWithFieldsParamSet
	FatalWithFieldsFunc               func(logging.Fields, string, ...interface{})

	statsInfoLock          sync.RWMutex
	statInfoFuncCallCount  int
	statInfoFuncCallParams []LoggerInfoParamSet
	InfoFunc               func(string, ...interface{})

	statsInfoWithFieldsLock          sync.RWMutex
	statInfoWithFieldsFuncCallCount  int
	statInfoWithFieldsFuncCallParams []LoggerInfoWithFieldsParamSet
	InfoWithFieldsFunc               func(logging.Fields, string, ...interface{})

	statsLogWithFieldsLock          sync.RWMutex
	statLogWithFieldsFuncCallCount  int
	statLogWithFieldsFuncCallParams []LoggerLogWithFieldsParamSet
	LogWithFieldsFunc               func(logging.LogLevel, logging.Fields, string, ...interface{})

	statsDebugLock          sync.RWMutex
	statDebugFuncCallCount  int
	statDebugFuncCallParams []LoggerDebugParamSet
	DebugFunc               func(string, ...interface{})

	statsErrorLock          sync.RWMutex
	statErrorFuncCallCount  int
	statErrorFuncCallParams []LoggerErrorParamSet
	ErrorFunc               func(string, ...interface{})

	statsErrorWithFieldsLock          sync.RWMutex
	statErrorWithFieldsFuncCallCount  int
	statErrorWithFieldsFuncCallParams []LoggerErrorWithFieldsParamSet
	ErrorWithFieldsFunc               func(logging.Fields, string, ...interface{})

	statsWarningLock          sync.RWMutex
	statWarningFuncCallCount  int
	statWarningFuncCallParams []LoggerWarningParamSet
	WarningFunc               func(string, ...interface{})

	statsWarningWithFieldsLock          sync.RWMutex
	statWarningWithFieldsFuncCallCount  int
	statWarningWithFieldsFuncCallParams []LoggerWarningWithFieldsParamSet
	WarningWithFieldsFunc               func(logging.Fields, string, ...interface{})

	statsWithFieldsLock          sync.RWMutex
	statWithFieldsFuncCallCount  int
	statWithFieldsFuncCallParams []LoggerWithFieldsParamSet
	WithFieldsFunc               func(logging.Fields) logging.Logger
}
type LoggerInfoParamSet struct {
	Arg0 string
	Arg1 []interface{}
}
type LoggerInfoWithFieldsParamSet struct {
	Arg0 logging.Fields
	Arg1 string
	Arg2 []interface{}
}
type LoggerLogWithFieldsParamSet struct {
	Arg0 logging.LogLevel
	Arg1 logging.Fields
	Arg2 string
	Arg3 []interface{}
}
type LoggerDebugParamSet struct {
	Arg0 string
	Arg1 []interface{}
}
type LoggerErrorParamSet struct {
	Arg0 string
	Arg1 []interface{}
}
type LoggerErrorWithFieldsParamSet struct {
	Arg0 logging.Fields
	Arg1 string
	Arg2 []interface{}
}
type LoggerFatalParamSet struct {
	Arg0 string
	Arg1 []interface{}
}
type LoggerFatalWithFieldsParamSet struct {
	Arg0 logging.Fields
	Arg1 string
	Arg2 []interface{}
}
type LoggerWarningParamSet struct {
	Arg0 string
	Arg1 []interface{}
}
type LoggerWarningWithFieldsParamSet struct {
	Arg0 logging.Fields
	Arg1 string
	Arg2 []interface{}
}
type LoggerWithFieldsParamSet struct {
	Arg0 logging.Fields
}
type LoggerDebugWithFieldsParamSet struct {
	Arg0 logging.Fields
	Arg1 string
	Arg2 []interface{}
}
type LoggerSyncParamSet struct{}

var _ logging.Logger = NewMockLogger()

func NewMockLogger() *MockLogger {
	m := &MockLogger{}
	m.DebugWithFieldsFunc = m.defaultDebugWithFieldsFunc
	m.SyncFunc = m.defaultSyncFunc
	m.ErrorWithFieldsFunc = m.defaultErrorWithFieldsFunc
	m.FatalFunc = m.defaultFatalFunc
	m.FatalWithFieldsFunc = m.defaultFatalWithFieldsFunc
	m.InfoFunc = m.defaultInfoFunc
	m.InfoWithFieldsFunc = m.defaultInfoWithFieldsFunc
	m.LogWithFieldsFunc = m.defaultLogWithFieldsFunc
	m.DebugFunc = m.defaultDebugFunc
	m.ErrorFunc = m.defaultErrorFunc
	m.WithFieldsFunc = m.defaultWithFieldsFunc
	m.WarningFunc = m.defaultWarningFunc
	m.WarningWithFieldsFunc = m.defaultWarningWithFieldsFunc
	return m
}
func (m *MockLogger) ErrorWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.statsErrorWithFieldsLock.Lock()
	m.statErrorWithFieldsFuncCallCount++
	m.statErrorWithFieldsFuncCallParams = append(m.statErrorWithFieldsFuncCallParams, LoggerErrorWithFieldsParamSet{v0, v1, v2})
	m.statsErrorWithFieldsLock.Unlock()
	m.ErrorWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) ErrorWithFieldsFuncCallCount() int {
	m.statsErrorWithFieldsLock.RLock()
	defer m.statsErrorWithFieldsLock.RUnlock()
	return m.statErrorWithFieldsFuncCallCount
}
func (m *MockLogger) ErrorWithFieldsFuncCallParams() []LoggerErrorWithFieldsParamSet {
	m.statsErrorWithFieldsLock.RLock()
	defer m.statsErrorWithFieldsLock.RUnlock()
	return m.statErrorWithFieldsFuncCallParams
}

func (m *MockLogger) Fatal(v0 string, v1 ...interface{}) {
	m.statsFatalLock.Lock()
	m.statFatalFuncCallCount++
	m.statFatalFuncCallParams = append(m.statFatalFuncCallParams, LoggerFatalParamSet{v0, v1})
	m.statsFatalLock.Unlock()
	m.FatalFunc(v0, v1...)
}
func (m *MockLogger) FatalFuncCallCount() int {
	m.statsFatalLock.RLock()
	defer m.statsFatalLock.RUnlock()
	return m.statFatalFuncCallCount
}
func (m *MockLogger) FatalFuncCallParams() []LoggerFatalParamSet {
	m.statsFatalLock.RLock()
	defer m.statsFatalLock.RUnlock()
	return m.statFatalFuncCallParams
}

func (m *MockLogger) FatalWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.statsFatalWithFieldsLock.Lock()
	m.statFatalWithFieldsFuncCallCount++
	m.statFatalWithFieldsFuncCallParams = append(m.statFatalWithFieldsFuncCallParams, LoggerFatalWithFieldsParamSet{v0, v1, v2})
	m.statsFatalWithFieldsLock.Unlock()
	m.FatalWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) FatalWithFieldsFuncCallCount() int {
	m.statsFatalWithFieldsLock.RLock()
	defer m.statsFatalWithFieldsLock.RUnlock()
	return m.statFatalWithFieldsFuncCallCount
}
func (m *MockLogger) FatalWithFieldsFuncCallParams() []LoggerFatalWithFieldsParamSet {
	m.statsFatalWithFieldsLock.RLock()
	defer m.statsFatalWithFieldsLock.RUnlock()
	return m.statFatalWithFieldsFuncCallParams
}

func (m *MockLogger) Info(v0 string, v1 ...interface{}) {
	m.statsInfoLock.Lock()
	m.statInfoFuncCallCount++
	m.statInfoFuncCallParams = append(m.statInfoFuncCallParams, LoggerInfoParamSet{v0, v1})
	m.statsInfoLock.Unlock()
	m.InfoFunc(v0, v1...)
}
func (m *MockLogger) InfoFuncCallCount() int {
	m.statsInfoLock.RLock()
	defer m.statsInfoLock.RUnlock()
	return m.statInfoFuncCallCount
}
func (m *MockLogger) InfoFuncCallParams() []LoggerInfoParamSet {
	m.statsInfoLock.RLock()
	defer m.statsInfoLock.RUnlock()
	return m.statInfoFuncCallParams
}

func (m *MockLogger) InfoWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.statsInfoWithFieldsLock.Lock()
	m.statInfoWithFieldsFuncCallCount++
	m.statInfoWithFieldsFuncCallParams = append(m.statInfoWithFieldsFuncCallParams, LoggerInfoWithFieldsParamSet{v0, v1, v2})
	m.statsInfoWithFieldsLock.Unlock()
	m.InfoWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) InfoWithFieldsFuncCallCount() int {
	m.statsInfoWithFieldsLock.RLock()
	defer m.statsInfoWithFieldsLock.RUnlock()
	return m.statInfoWithFieldsFuncCallCount
}
func (m *MockLogger) InfoWithFieldsFuncCallParams() []LoggerInfoWithFieldsParamSet {
	m.statsInfoWithFieldsLock.RLock()
	defer m.statsInfoWithFieldsLock.RUnlock()
	return m.statInfoWithFieldsFuncCallParams
}

func (m *MockLogger) LogWithFields(v0 logging.LogLevel, v1 logging.Fields, v2 string, v3 ...interface{}) {
	m.statsLogWithFieldsLock.Lock()
	m.statLogWithFieldsFuncCallCount++
	m.statLogWithFieldsFuncCallParams = append(m.statLogWithFieldsFuncCallParams, LoggerLogWithFieldsParamSet{v0, v1, v2, v3})
	m.statsLogWithFieldsLock.Unlock()
	m.LogWithFieldsFunc(v0, v1, v2, v3...)
}
func (m *MockLogger) LogWithFieldsFuncCallCount() int {
	m.statsLogWithFieldsLock.RLock()
	defer m.statsLogWithFieldsLock.RUnlock()
	return m.statLogWithFieldsFuncCallCount
}
func (m *MockLogger) LogWithFieldsFuncCallParams() []LoggerLogWithFieldsParamSet {
	m.statsLogWithFieldsLock.RLock()
	defer m.statsLogWithFieldsLock.RUnlock()
	return m.statLogWithFieldsFuncCallParams
}

func (m *MockLogger) Debug(v0 string, v1 ...interface{}) {
	m.statsDebugLock.Lock()
	m.statDebugFuncCallCount++
	m.statDebugFuncCallParams = append(m.statDebugFuncCallParams, LoggerDebugParamSet{v0, v1})
	m.statsDebugLock.Unlock()
	m.DebugFunc(v0, v1...)
}
func (m *MockLogger) DebugFuncCallCount() int {
	m.statsDebugLock.RLock()
	defer m.statsDebugLock.RUnlock()
	return m.statDebugFuncCallCount
}
func (m *MockLogger) DebugFuncCallParams() []LoggerDebugParamSet {
	m.statsDebugLock.RLock()
	defer m.statsDebugLock.RUnlock()
	return m.statDebugFuncCallParams
}

func (m *MockLogger) Error(v0 string, v1 ...interface{}) {
	m.statsErrorLock.Lock()
	m.statErrorFuncCallCount++
	m.statErrorFuncCallParams = append(m.statErrorFuncCallParams, LoggerErrorParamSet{v0, v1})
	m.statsErrorLock.Unlock()
	m.ErrorFunc(v0, v1...)
}
func (m *MockLogger) ErrorFuncCallCount() int {
	m.statsErrorLock.RLock()
	defer m.statsErrorLock.RUnlock()
	return m.statErrorFuncCallCount
}
func (m *MockLogger) ErrorFuncCallParams() []LoggerErrorParamSet {
	m.statsErrorLock.RLock()
	defer m.statsErrorLock.RUnlock()
	return m.statErrorFuncCallParams
}

func (m *MockLogger) WithFields(v0 logging.Fields) logging.Logger {
	m.statsWithFieldsLock.Lock()
	m.statWithFieldsFuncCallCount++
	m.statWithFieldsFuncCallParams = append(m.statWithFieldsFuncCallParams, LoggerWithFieldsParamSet{v0})
	m.statsWithFieldsLock.Unlock()
	return m.WithFieldsFunc(v0)
}
func (m *MockLogger) WithFieldsFuncCallCount() int {
	m.statsWithFieldsLock.RLock()
	defer m.statsWithFieldsLock.RUnlock()
	return m.statWithFieldsFuncCallCount
}
func (m *MockLogger) WithFieldsFuncCallParams() []LoggerWithFieldsParamSet {
	m.statsWithFieldsLock.RLock()
	defer m.statsWithFieldsLock.RUnlock()
	return m.statWithFieldsFuncCallParams
}

func (m *MockLogger) Warning(v0 string, v1 ...interface{}) {
	m.statsWarningLock.Lock()
	m.statWarningFuncCallCount++
	m.statWarningFuncCallParams = append(m.statWarningFuncCallParams, LoggerWarningParamSet{v0, v1})
	m.statsWarningLock.Unlock()
	m.WarningFunc(v0, v1...)
}
func (m *MockLogger) WarningFuncCallCount() int {
	m.statsWarningLock.RLock()
	defer m.statsWarningLock.RUnlock()
	return m.statWarningFuncCallCount
}
func (m *MockLogger) WarningFuncCallParams() []LoggerWarningParamSet {
	m.statsWarningLock.RLock()
	defer m.statsWarningLock.RUnlock()
	return m.statWarningFuncCallParams
}

func (m *MockLogger) WarningWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.statsWarningWithFieldsLock.Lock()
	m.statWarningWithFieldsFuncCallCount++
	m.statWarningWithFieldsFuncCallParams = append(m.statWarningWithFieldsFuncCallParams, LoggerWarningWithFieldsParamSet{v0, v1, v2})
	m.statsWarningWithFieldsLock.Unlock()
	m.WarningWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) WarningWithFieldsFuncCallCount() int {
	m.statsWarningWithFieldsLock.RLock()
	defer m.statsWarningWithFieldsLock.RUnlock()
	return m.statWarningWithFieldsFuncCallCount
}
func (m *MockLogger) WarningWithFieldsFuncCallParams() []LoggerWarningWithFieldsParamSet {
	m.statsWarningWithFieldsLock.RLock()
	defer m.statsWarningWithFieldsLock.RUnlock()
	return m.statWarningWithFieldsFuncCallParams
}

func (m *MockLogger) DebugWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.statsDebugWithFieldsLock.Lock()
	m.statDebugWithFieldsFuncCallCount++
	m.statDebugWithFieldsFuncCallParams = append(m.statDebugWithFieldsFuncCallParams, LoggerDebugWithFieldsParamSet{v0, v1, v2})
	m.statsDebugWithFieldsLock.Unlock()
	m.DebugWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) DebugWithFieldsFuncCallCount() int {
	m.statsDebugWithFieldsLock.RLock()
	defer m.statsDebugWithFieldsLock.RUnlock()
	return m.statDebugWithFieldsFuncCallCount
}
func (m *MockLogger) DebugWithFieldsFuncCallParams() []LoggerDebugWithFieldsParamSet {
	m.statsDebugWithFieldsLock.RLock()
	defer m.statsDebugWithFieldsLock.RUnlock()
	return m.statDebugWithFieldsFuncCallParams
}

func (m *MockLogger) Sync() error {
	m.statsSyncLock.Lock()
	m.statSyncFuncCallCount++
	m.statSyncFuncCallParams = append(m.statSyncFuncCallParams, LoggerSyncParamSet{})
	m.statsSyncLock.Unlock()
	return m.SyncFunc()
}
func (m *MockLogger) SyncFuncCallCount() int {
	m.statsSyncLock.RLock()
	defer m.statsSyncLock.RUnlock()
	return m.statSyncFuncCallCount
}
func (m *MockLogger) SyncFuncCallParams() []LoggerSyncParamSet {
	m.statsSyncLock.RLock()
	defer m.statsSyncLock.RUnlock()
	return m.statSyncFuncCallParams
}

func (m *MockLogger) defaultInfoFunc(v0 string, v1 ...interface{}) {
	return
}
func (m *MockLogger) defaultInfoWithFieldsFunc(v0 logging.Fields, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultLogWithFieldsFunc(v0 logging.LogLevel, v1 logging.Fields, v2 string, v3 ...interface{}) {
	return
}
func (m *MockLogger) defaultDebugFunc(v0 string, v1 ...interface{}) {
	return
}
func (m *MockLogger) defaultErrorFunc(v0 string, v1 ...interface{}) {
	return
}
func (m *MockLogger) defaultErrorWithFieldsFunc(v0 logging.Fields, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultFatalFunc(v0 string, v1 ...interface{}) {
	return
}
func (m *MockLogger) defaultFatalWithFieldsFunc(v0 logging.Fields, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultWarningFunc(v0 string, v1 ...interface{}) {
	return
}
func (m *MockLogger) defaultWarningWithFieldsFunc(v0 logging.Fields, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultWithFieldsFunc(v0 logging.Fields) logging.Logger {
	return nil
}
func (m *MockLogger) defaultDebugWithFieldsFunc(v0 logging.Fields, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultSyncFunc() error {
	return nil
}
