package process

import (
	"time"

	"github.com/efritz/glock"
)

func makeTimeoutChan(clock glock.Clock, timeout time.Duration) <-chan time.Time {
	if timeout == 0 {
		return nil
	}

	return clock.After(timeout)
}

func makeErrChan(f func() error) <-chan error {
	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		ch <- f()
	}()

	return ch
}
