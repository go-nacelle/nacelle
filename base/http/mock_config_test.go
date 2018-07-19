// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/nacelle/config -i Config -o mock_config_test.go -f

package http

import (
	config "github.com/efritz/nacelle/config"
	tag "github.com/efritz/nacelle/config/tag"
	sync "sync"
)

type MockConfig struct {
	statsLoadLock          sync.RWMutex
	statLoadFuncCallCount  int
	statLoadFuncCallParams []ConfigLoadParamSet
	LoadFunc               func(interface{}, ...tag.Modifier) error

	statsMustLoadLock          sync.RWMutex
	statMustLoadFuncCallCount  int
	statMustLoadFuncCallParams []ConfigMustLoadParamSet
	MustLoadFunc               func(interface{}, ...tag.Modifier)
}
type ConfigLoadParamSet struct {
	Arg0 interface{}
	Arg1 []tag.Modifier
}
type ConfigMustLoadParamSet struct {
	Arg0 interface{}
	Arg1 []tag.Modifier
}

var _ config.Config = NewMockConfig()

func NewMockConfig() *MockConfig {
	m := &MockConfig{}
	m.LoadFunc = m.defaultLoadFunc
	m.MustLoadFunc = m.defaultMustLoadFunc
	return m
}
func (m *MockConfig) Load(v0 interface{}, v1 ...tag.Modifier) error {
	m.statsLoadLock.Lock()
	m.statLoadFuncCallCount++
	m.statLoadFuncCallParams = append(m.statLoadFuncCallParams, ConfigLoadParamSet{v0, v1})
	m.statsLoadLock.Unlock()
	return m.LoadFunc(v0, v1...)
}
func (m *MockConfig) LoadFuncCallCount() int {
	m.statsLoadLock.RLock()
	defer m.statsLoadLock.RUnlock()
	return m.statLoadFuncCallCount
}
func (m *MockConfig) LoadFuncCallParams() []ConfigLoadParamSet {
	m.statsLoadLock.RLock()
	defer m.statsLoadLock.RUnlock()
	return m.statLoadFuncCallParams
}

func (m *MockConfig) MustLoad(v0 interface{}, v1 ...tag.Modifier) {
	m.statsMustLoadLock.Lock()
	m.statMustLoadFuncCallCount++
	m.statMustLoadFuncCallParams = append(m.statMustLoadFuncCallParams, ConfigMustLoadParamSet{v0, v1})
	m.statsMustLoadLock.Unlock()
	m.MustLoadFunc(v0, v1...)
}
func (m *MockConfig) MustLoadFuncCallCount() int {
	m.statsMustLoadLock.RLock()
	defer m.statsMustLoadLock.RUnlock()
	return m.statMustLoadFuncCallCount
}
func (m *MockConfig) MustLoadFuncCallParams() []ConfigMustLoadParamSet {
	m.statsMustLoadLock.RLock()
	defer m.statsMustLoadLock.RUnlock()
	return m.statMustLoadFuncCallParams
}

func (m *MockConfig) defaultLoadFunc(v0 interface{}, v1 ...tag.Modifier) error {
	return nil
}
func (m *MockConfig) defaultMustLoadFunc(v0 interface{}, v1 ...tag.Modifier) {
	return
}
