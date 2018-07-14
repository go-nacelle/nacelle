package process

import (
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type HealthSuite struct{}

func (s *HealthSuite) TestReasons(t sweet.T) {
	now := time.Now()
	clock := glock.NewMockClock()
	clock.SetCurrent(now)

	health := NewHealth(WithHealthClock(clock))
	Expect(health.LastChange()).To(BeZero())

	clock.Advance(time.Hour)
	Expect(health.AddReason("foo")).To(BeNil())
	clock.Advance(time.Minute)
	Expect(health.AddReason("bar")).To(BeNil())
	clock.Advance(time.Minute)
	Expect(health.AddReason("baz")).To(BeNil())
	Expect(health.RemoveReason("bar")).To(BeNil())

	Expect(health.Reasons()).To(ConsistOf([]Reason{
		Reason{Key: "foo", Added: now.Add(time.Hour)},
		Reason{Key: "baz", Added: now.Add(time.Hour + time.Minute*2)},
	}))
}

func (s *HealthSuite) TestLastChangedTime(t sweet.T) {
	clock := glock.NewMockClock()
	clock.SetCurrent(time.Now())

	health := NewHealth(WithHealthClock(clock))
	Expect(health.LastChange()).To(BeZero())

	// Changed
	clock.Advance(time.Hour)
	Expect(health.AddReason("foo")).To(BeNil())
	clock.Advance(time.Minute * 2)
	Expect(health.LastChange()).To(Equal(time.Minute * 2))

	// No change
	Expect(health.AddReason("bar")).To(BeNil())
	clock.Advance(time.Minute * 4)
	Expect(health.LastChange()).To(Equal(time.Minute * 6))

	// No change
	Expect(health.RemoveReason("foo")).To(BeNil())
	clock.Advance(time.Minute * 2)
	Expect(health.LastChange()).To(Equal(time.Minute * 8))

	// Changed
	Expect(health.RemoveReason("bar")).To(BeNil())
	Expect(health.LastChange()).To(Equal(time.Minute * 0))
}

func (s *HealthSuite) TestAddReasonError(t sweet.T) {
	health := NewHealth()
	health.AddReason("foo")
	Expect(health.AddReason("foo")).To(MatchError("reason foo already registered"))
}

func (s *HealthSuite) TestRemoveReasonError(t sweet.T) {
	health := NewHealth()
	Expect(health.RemoveReason("foo")).To(MatchError("reason foo not registered"))
}
