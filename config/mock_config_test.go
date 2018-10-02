// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/nacelle/config -i Config -o mock_config_test.go -f

package config

import (
	tag "github.com/efritz/nacelle/config/tag"
	"sync"
)

type MockConfig struct {
	LoadFunc     func(interface{}, ...tag.Modifier) error
	histLoad     []ConfigLoadParamSet
	MustLoadFunc func(interface{}, ...tag.Modifier)
	histMustLoad []ConfigMustLoadParamSet
	mutex        sync.RWMutex
}
type ConfigLoadParamSet struct {
	Arg0 interface{}
	Arg1 []tag.Modifier
}
type ConfigMustLoadParamSet struct {
	Arg0 interface{}
	Arg1 []tag.Modifier
}

func NewMockConfig() *MockConfig {
	m := &MockConfig{}
	m.LoadFunc = m.defaultLoadFunc
	m.MustLoadFunc = m.defaultMustLoadFunc
	return m
}
func (m *MockConfig) Load(v0 interface{}, v1 ...tag.Modifier) error {
	m.mutex.Lock()
	m.histLoad = append(m.histLoad, ConfigLoadParamSet{v0, v1})
	m.mutex.Unlock()
	return m.LoadFunc(v0, v1...)
}
func (m *MockConfig) LoadFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histLoad)
}
func (m *MockConfig) LoadFuncCallParams() []ConfigLoadParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histLoad
}

func (m *MockConfig) MustLoad(v0 interface{}, v1 ...tag.Modifier) {
	m.mutex.Lock()
	m.histMustLoad = append(m.histMustLoad, ConfigMustLoadParamSet{v0, v1})
	m.mutex.Unlock()
	m.MustLoadFunc(v0, v1...)
}
func (m *MockConfig) MustLoadFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histMustLoad)
}
func (m *MockConfig) MustLoadFuncCallParams() []ConfigMustLoadParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histMustLoad
}

func (m *MockConfig) defaultLoadFunc(v0 interface{}, v1 ...tag.Modifier) error {
	return nil
}
func (m *MockConfig) defaultMustLoadFunc(v0 interface{}, v1 ...tag.Modifier) {
	return
}
