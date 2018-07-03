package logging

import (
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type LoggerSuite struct{}

func (s *LoggerSuite) TestNormalizeTimeValues(t sweet.T) {
	var (
		t1     = time.Unix(1503939881, 0)
		t2     = time.Unix(1503939891, 0)
		fields = Fields{
			"foo":  "bar",
			"bar":  t1,
			"baz":  t2,
			"bonk": []bool{true, false, true},
		}
	)

	// Modifies object in-place
	Expect(fields.normalizeTimeValues()).To(Equal(fields))

	// Non-time values remain the same
	Expect(fields["foo"]).To(Equal("bar"))
	Expect(fields["bonk"]).To(Equal([]bool{true, false, true}))

	// Times converted to ISO 8601
	Expect(time.Parse(JSONTimeFormat, fields["bar"].(string))).To(Equal(t1))
	Expect(time.Parse(JSONTimeFormat, fields["baz"].(string))).To(Equal(t2))
}

func (s *LoggerSuite) TestNormalizeTimeValuesOnNilFields(t sweet.T) {
	Expect(Fields(nil).normalizeTimeValues()).To(BeNil())
}
