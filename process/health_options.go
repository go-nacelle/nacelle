package process

import "github.com/efritz/glock"

type HealthConfigFunc func(*health)

func WithHealthClock(clock glock.Clock) HealthConfigFunc {
	return func(h *health) { h.clock = clock }
}
