package log

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CallerSuite struct{}

func (s *CallerSuite) TestTrimPath(t sweet.T) {
	Expect(trimPath("")).To(Equal(""))
	Expect(trimPath("/")).To(Equal("/"))
	Expect(trimPath("/foo")).To(Equal("/foo"))
	Expect(trimPath("/foo/bar")).To(Equal("foo/bar"))
	Expect(trimPath("/foo/bar/baz")).To(Equal("bar/baz"))
	Expect(trimPath("/foo/bar/baz/bonk")).To(Equal("baz/bonk"))
}
