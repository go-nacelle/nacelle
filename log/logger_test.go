package log

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type LoggerSuite struct{}

func (s *LoggerSuite) TestNormalizeTimeValues(t sweet.T) {
	fields := Fields(map[string]interface{}{
		"foo":  "bar",
		"bar":  time.Unix(1503939881, 0),
		"baz":  time.Unix(1503939891, 0),
		"bonk": []bool{true, false, true},
	})

	// Modifies object in-place
	Expect(fields.normalizeTimeValues()).To(Equal(fields))

	// Non-time values remain the same
	Expect(fields["foo"]).To(Equal("bar"))
	Expect(fields["bonk"]).To(Equal([]bool{true, false, true}))

	// Times converted to ISO 8601
	Expect(fields["bar"]).To(Equal("2017-08-28T12:04:41.000-0500"))
	Expect(fields["baz"]).To(Equal("2017-08-28T12:04:51.000-0500"))
}

func (s *LoggerSuite) TestNormalizeTimeValuesOnNilFields(t sweet.T) {
	Expect(Fields(nil).normalizeTimeValues()).To(BeNil())
}
