// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/nacelle/logging -i Logger -o mock_logger_test.go -f

package config

import (
	logging "github.com/efritz/nacelle/logging"
	"sync"
)

type MockLogger struct {
	DebugFunc             func(string, ...interface{})
	histDebug             []LoggerDebugParamSet
	DebugWithFieldsFunc   func(logging.Fields, string, ...interface{})
	histDebugWithFields   []LoggerDebugWithFieldsParamSet
	ErrorFunc             func(string, ...interface{})
	histError             []LoggerErrorParamSet
	ErrorWithFieldsFunc   func(logging.Fields, string, ...interface{})
	histErrorWithFields   []LoggerErrorWithFieldsParamSet
	FatalFunc             func(string, ...interface{})
	histFatal             []LoggerFatalParamSet
	FatalWithFieldsFunc   func(logging.Fields, string, ...interface{})
	histFatalWithFields   []LoggerFatalWithFieldsParamSet
	InfoFunc              func(string, ...interface{})
	histInfo              []LoggerInfoParamSet
	InfoWithFieldsFunc    func(logging.Fields, string, ...interface{})
	histInfoWithFields    []LoggerInfoWithFieldsParamSet
	LogWithFieldsFunc     func(logging.LogLevel, logging.Fields, string, ...interface{})
	histLogWithFields     []LoggerLogWithFieldsParamSet
	SyncFunc              func() error
	histSync              []LoggerSyncParamSet
	WarningFunc           func(string, ...interface{})
	histWarning           []LoggerWarningParamSet
	WarningWithFieldsFunc func(logging.Fields, string, ...interface{})
	histWarningWithFields []LoggerWarningWithFieldsParamSet
	WithFieldsFunc        func(logging.Fields) logging.Logger
	histWithFields        []LoggerWithFieldsParamSet
	mutex                 sync.RWMutex
}
type LoggerDebugParamSet struct {
	Arg0 string
	Arg1 []interface{}
}
type LoggerDebugWithFieldsParamSet struct {
	Arg0 logging.Fields
	Arg1 string
	Arg2 []interface{}
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
type LoggerSyncParamSet struct{}
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

func NewMockLogger() *MockLogger {
	m := &MockLogger{}
	m.DebugFunc = m.defaultDebugFunc
	m.DebugWithFieldsFunc = m.defaultDebugWithFieldsFunc
	m.ErrorFunc = m.defaultErrorFunc
	m.ErrorWithFieldsFunc = m.defaultErrorWithFieldsFunc
	m.FatalFunc = m.defaultFatalFunc
	m.FatalWithFieldsFunc = m.defaultFatalWithFieldsFunc
	m.InfoFunc = m.defaultInfoFunc
	m.InfoWithFieldsFunc = m.defaultInfoWithFieldsFunc
	m.LogWithFieldsFunc = m.defaultLogWithFieldsFunc
	m.SyncFunc = m.defaultSyncFunc
	m.WarningFunc = m.defaultWarningFunc
	m.WarningWithFieldsFunc = m.defaultWarningWithFieldsFunc
	m.WithFieldsFunc = m.defaultWithFieldsFunc
	return m
}
func (m *MockLogger) Debug(v0 string, v1 ...interface{}) {
	m.mutex.Lock()
	m.histDebug = append(m.histDebug, LoggerDebugParamSet{v0, v1})
	m.mutex.Unlock()
	m.DebugFunc(v0, v1...)
}
func (m *MockLogger) DebugFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histDebug)
}
func (m *MockLogger) DebugFuncCallParams() []LoggerDebugParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histDebug
}

func (m *MockLogger) DebugWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histDebugWithFields = append(m.histDebugWithFields, LoggerDebugWithFieldsParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.DebugWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) DebugWithFieldsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histDebugWithFields)
}
func (m *MockLogger) DebugWithFieldsFuncCallParams() []LoggerDebugWithFieldsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histDebugWithFields
}

func (m *MockLogger) Error(v0 string, v1 ...interface{}) {
	m.mutex.Lock()
	m.histError = append(m.histError, LoggerErrorParamSet{v0, v1})
	m.mutex.Unlock()
	m.ErrorFunc(v0, v1...)
}
func (m *MockLogger) ErrorFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histError)
}
func (m *MockLogger) ErrorFuncCallParams() []LoggerErrorParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histError
}

func (m *MockLogger) ErrorWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histErrorWithFields = append(m.histErrorWithFields, LoggerErrorWithFieldsParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.ErrorWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) ErrorWithFieldsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histErrorWithFields)
}
func (m *MockLogger) ErrorWithFieldsFuncCallParams() []LoggerErrorWithFieldsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histErrorWithFields
}

func (m *MockLogger) Fatal(v0 string, v1 ...interface{}) {
	m.mutex.Lock()
	m.histFatal = append(m.histFatal, LoggerFatalParamSet{v0, v1})
	m.mutex.Unlock()
	m.FatalFunc(v0, v1...)
}
func (m *MockLogger) FatalFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histFatal)
}
func (m *MockLogger) FatalFuncCallParams() []LoggerFatalParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histFatal
}

func (m *MockLogger) FatalWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histFatalWithFields = append(m.histFatalWithFields, LoggerFatalWithFieldsParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.FatalWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) FatalWithFieldsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histFatalWithFields)
}
func (m *MockLogger) FatalWithFieldsFuncCallParams() []LoggerFatalWithFieldsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histFatalWithFields
}

func (m *MockLogger) Info(v0 string, v1 ...interface{}) {
	m.mutex.Lock()
	m.histInfo = append(m.histInfo, LoggerInfoParamSet{v0, v1})
	m.mutex.Unlock()
	m.InfoFunc(v0, v1...)
}
func (m *MockLogger) InfoFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histInfo)
}
func (m *MockLogger) InfoFuncCallParams() []LoggerInfoParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histInfo
}

func (m *MockLogger) InfoWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histInfoWithFields = append(m.histInfoWithFields, LoggerInfoWithFieldsParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.InfoWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) InfoWithFieldsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histInfoWithFields)
}
func (m *MockLogger) InfoWithFieldsFuncCallParams() []LoggerInfoWithFieldsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histInfoWithFields
}

func (m *MockLogger) LogWithFields(v0 logging.LogLevel, v1 logging.Fields, v2 string, v3 ...interface{}) {
	m.mutex.Lock()
	m.histLogWithFields = append(m.histLogWithFields, LoggerLogWithFieldsParamSet{v0, v1, v2, v3})
	m.mutex.Unlock()
	m.LogWithFieldsFunc(v0, v1, v2, v3...)
}
func (m *MockLogger) LogWithFieldsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histLogWithFields)
}
func (m *MockLogger) LogWithFieldsFuncCallParams() []LoggerLogWithFieldsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histLogWithFields
}

func (m *MockLogger) Sync() error {
	m.mutex.Lock()
	m.histSync = append(m.histSync, LoggerSyncParamSet{})
	m.mutex.Unlock()
	return m.SyncFunc()
}
func (m *MockLogger) SyncFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histSync)
}
func (m *MockLogger) SyncFuncCallParams() []LoggerSyncParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histSync
}

func (m *MockLogger) Warning(v0 string, v1 ...interface{}) {
	m.mutex.Lock()
	m.histWarning = append(m.histWarning, LoggerWarningParamSet{v0, v1})
	m.mutex.Unlock()
	m.WarningFunc(v0, v1...)
}
func (m *MockLogger) WarningFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histWarning)
}
func (m *MockLogger) WarningFuncCallParams() []LoggerWarningParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histWarning
}

func (m *MockLogger) WarningWithFields(v0 logging.Fields, v1 string, v2 ...interface{}) {
	m.mutex.Lock()
	m.histWarningWithFields = append(m.histWarningWithFields, LoggerWarningWithFieldsParamSet{v0, v1, v2})
	m.mutex.Unlock()
	m.WarningWithFieldsFunc(v0, v1, v2...)
}
func (m *MockLogger) WarningWithFieldsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histWarningWithFields)
}
func (m *MockLogger) WarningWithFieldsFuncCallParams() []LoggerWarningWithFieldsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histWarningWithFields
}

func (m *MockLogger) WithFields(v0 logging.Fields) logging.Logger {
	m.mutex.Lock()
	m.histWithFields = append(m.histWithFields, LoggerWithFieldsParamSet{v0})
	m.mutex.Unlock()
	return m.WithFieldsFunc(v0)
}
func (m *MockLogger) WithFieldsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histWithFields)
}
func (m *MockLogger) WithFieldsFuncCallParams() []LoggerWithFieldsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histWithFields
}

func (m *MockLogger) defaultDebugFunc(v0 string, v1 ...interface{}) {
	return
}
func (m *MockLogger) defaultDebugWithFieldsFunc(v0 logging.Fields, v1 string, v2 ...interface{}) {
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
func (m *MockLogger) defaultInfoFunc(v0 string, v1 ...interface{}) {
	return
}
func (m *MockLogger) defaultInfoWithFieldsFunc(v0 logging.Fields, v1 string, v2 ...interface{}) {
	return
}
func (m *MockLogger) defaultLogWithFieldsFunc(v0 logging.LogLevel, v1 logging.Fields, v2 string, v3 ...interface{}) {
	return
}
func (m *MockLogger) defaultSyncFunc() error {
	return nil
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
