package process

import (
	"fmt"
	"sync"
	"time"

	"github.com/efritz/glock"
)

type (
	Health interface {
		Reasons() []Reason
		LastChange() time.Duration
		AddReason(key interface{}) error
		RemoveReason(key interface{}) error
	}

	health struct {
		reasons    map[interface{}]Reason
		lastChange time.Time
		mutex      sync.RWMutex
		clock      glock.Clock
	}

	Reason struct {
		Key   interface{}
		Added time.Time
	}
)

func NewHealth(configs ...HealthConfigFunc) Health {
	h := &health{
		reasons: map[interface{}]Reason{},
		clock:   glock.NewRealClock(),
	}

	for _, f := range configs {
		f(h)
	}

	h.lastChange = h.clock.Now()
	return h
}

func (h *health) Reasons() []Reason {
	reasons := []Reason{}
	for _, reason := range h.reasons {
		reasons = append(reasons, reason)
	}

	return reasons
}

func (h *health) LastChange() time.Duration {
	return h.clock.Now().Sub(h.lastChange)
}

func (h *health) AddReason(key interface{}) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.reasons[key]; ok {
		return fmt.Errorf("reason %s already registered", key)
	}

	now := h.clock.Now()

	if len(h.reasons) == 0 {
		h.lastChange = now
	}

	h.reasons[key] = Reason{
		Key:   key,
		Added: now,
	}

	return nil
}

func (h *health) RemoveReason(key interface{}) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.reasons[key]; !ok {
		return fmt.Errorf("reason %s not registered", key)
	}

	delete(h.reasons, key)

	if len(h.reasons) == 0 {
		h.lastChange = h.clock.Now()
	}

	return nil
}
