package config

//go:generate go-mockgen github.com/efritz/nacelle/config -i Sourcer -o mock_sourcer_test.go -f

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type MultiSourcerSuite struct{}

func (s *MultiSourcerSuite) TestMultiSourcerBasic(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"env"} }
	s2.TagsFunc = func() []string { return []string{"env"} }

	s1.GetFunc = func(values []string) (string, bool, bool) {
		if values[0] == "foo" {
			return "bar", false, true
		}

		return "", false, false
	}

	s2.GetFunc = func(values []string) (string, bool, bool) {
		if values[0] == "bar" {
			return "baz", false, true
		}

		return "", false, false
	}

	multi := NewMultiSourcer(s2, s1)
	ensureEquals(multi, []string{"foo"}, "bar")
	ensureEquals(multi, []string{"bar"}, "baz")
	ensureMissing(multi, []string{"baz"})
}

func (s *MultiSourcerSuite) TestMultiSourcerPriority(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"env"} }
	s2.TagsFunc = func() []string { return []string{"env"} }

	s1.GetFunc = func(values []string) (string, bool, bool) {
		return "bar", false, true
	}

	s2.GetFunc = func(values []string) (string, bool, bool) {
		return "baz", false, true
	}

	multi := NewMultiSourcer(s2, s1)
	ensureEquals(multi, []string{"foo"}, "bar")
}

func (s *MultiSourcerSuite) TestMultiSourcerTags(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s4 := NewMockSourcer()
	s5 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"a"} }
	s2.TagsFunc = func() []string { return []string{"b"} }
	s3.TagsFunc = func() []string { return []string{"c"} }
	s4.TagsFunc = func() []string { return []string{"a", "b", "d"} }
	s5.TagsFunc = func() []string { return []string{"e"} }

	multi := NewMultiSourcer(s5, s4, s3, s2, s1)
	tags := multi.Tags()
	Expect(tags).To(HaveLen(5))
	Expect(tags).To(ConsistOf("a", "b", "c", "d", "e"))
}

func (s *MultiSourcerSuite) TestMultiSourcerDifferentTags(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"a"} }
	s2.TagsFunc = func() []string { return []string{"b"} }
	s3.TagsFunc = func() []string { return []string{"a"} }

	s1.GetFunc = func(values []string) (string, bool, bool) {
		Expect(values).To(Equal([]string{"foo"}))
		return "", true, false
	}

	s2.GetFunc = func(values []string) (string, bool, bool) {
		Expect(values).To(Equal([]string{"bar"}))
		return "", true, false
	}

	s3.GetFunc = func(values []string) (string, bool, bool) {
		Expect(values).To(Equal([]string{"foo"}))
		return "", false, false
	}

	multi := NewMultiSourcer(s3, s2, s1)
	_, skip, ok := multi.Get([]string{"foo", "bar"})
	Expect(ok).To(BeFalse())
	Expect(skip).To(BeFalse())
	Expect(s1.GetFuncCallCount()).To(Equal(1))
	Expect(s2.GetFuncCallCount()).To(Equal(1))
	Expect(s3.GetFuncCallCount()).To(Equal(1))
}

func (s *MultiSourcerSuite) TestMultiSourceSkip(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"a"} }
	s2.TagsFunc = func() []string { return []string{"b"} }
	s3.TagsFunc = func() []string { return []string{"a"} }

	s1.GetFunc = func(values []string) (string, bool, bool) {
		return "", true, false
	}

	s2.GetFunc = func(values []string) (string, bool, bool) {
		return "", true, false
	}

	s3.GetFunc = func(values []string) (string, bool, bool) {
		return "", true, false
	}

	multi := NewMultiSourcer(s3, s2, s1)
	_, skip, ok := multi.Get([]string{"", ""})
	Expect(ok).To(BeFalse())
	Expect(skip).To(BeTrue())
}
