// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/nacelle/config -i Sourcer -o mock_sourcer_test.go -f

package config

import "sync"

type MockSourcer struct {
	GetFunc  func([]string) (string, bool, bool)
	histGet  []SourcerGetParamSet
	TagsFunc func() []string
	histTags []SourcerTagsParamSet
	mutex    sync.RWMutex
}
type SourcerGetParamSet struct {
	Arg0 []string
}
type SourcerTagsParamSet struct{}

func NewMockSourcer() *MockSourcer {
	m := &MockSourcer{}
	m.GetFunc = m.defaultGetFunc
	m.TagsFunc = m.defaultTagsFunc
	return m
}
func (m *MockSourcer) Get(v0 []string) (string, bool, bool) {
	m.mutex.Lock()
	m.histGet = append(m.histGet, SourcerGetParamSet{v0})
	m.mutex.Unlock()
	return m.GetFunc(v0)
}
func (m *MockSourcer) GetFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histGet)
}
func (m *MockSourcer) GetFuncCallParams() []SourcerGetParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histGet
}

func (m *MockSourcer) Tags() []string {
	m.mutex.Lock()
	m.histTags = append(m.histTags, SourcerTagsParamSet{})
	m.mutex.Unlock()
	return m.TagsFunc()
}
func (m *MockSourcer) TagsFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histTags)
}
func (m *MockSourcer) TagsFuncCallParams() []SourcerTagsParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histTags
}

func (m *MockSourcer) defaultGetFunc(v0 []string) (string, bool, bool) {
	return "", false, false
}
func (m *MockSourcer) defaultTagsFunc() []string {
	return nil
}
