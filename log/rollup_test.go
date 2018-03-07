package log

import (
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type RollupSuite struct{}

func (s *RollupSuite) TestRollupSimilarMessages(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newRollupShim(adaptShim(shim), clock, time.Second)
	)

	for i := 1; i <= 20; i++ {
		// Logged, starting window
		adapter.LogWithFields(LevelDebug, nil, "a")
		Expect(shim.messages).To(HaveLen(2*i - 1))

		// Stashed
		adapter.LogWithFields(LevelDebug, nil, "a")
		adapter.LogWithFields(LevelDebug, nil, "a")
		Expect(shim.messages).To(HaveLen(2*i - 1))

		// Flushed
		clock.BlockingAdvance(time.Second)
		Eventually(func() []*logMessage { return shim.messages }).Should(HaveLen(2 * i))
		Expect(shim.messages[2*i-1].fields[FieldRollup]).To(Equal(2))
	}
}

func (s *RollupSuite) TestRollupInactivity(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newRollupShim(adaptShim(shim), clock, time.Second)
	)

	for i := 0; i < 20; i++ {
		adapter.LogWithFields(LevelDebug, nil, "a")
		clock.Advance(time.Second * 2)
	}

	// All messages present
	Eventually(func() []*logMessage { return shim.messages }).Should(HaveLen(20))
}

func (s *RollupSuite) TestRollupFlushesRelativeToFirstMessage(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newRollupShim(adaptShim(shim), clock, time.Second)
	)

	adapter.LogWithFields(LevelDebug, nil, "a")
	clock.Advance(time.Millisecond * 500)

	for i := 0; i < 90; i++ {
		adapter.LogWithFields(LevelDebug, nil, "a")
		clock.Advance(time.Millisecond * 5)
	}

	clock.BlockingAdvance(time.Millisecond * 50)
	Eventually(func() []*logMessage { return shim.messages }).Should(HaveLen(2))
}

func (s *RollupSuite) TestAllDistinctMessages(t sweet.T) {
	var (
		shim    = &testShim{}
		clock   = glock.NewMockClock()
		adapter = newRollupShim(adaptShim(shim), clock, time.Second)
	)

	for i := 0; i < 10; i++ {
		adapter.LogWithFields(LevelDebug, nil, "a")
		adapter.LogWithFields(LevelDebug, nil, "b")
		adapter.LogWithFields(LevelDebug, nil, "c")
	}

	Expect(shim.messages).To(HaveLen(3))
}
