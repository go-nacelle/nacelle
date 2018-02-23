package nacelle

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type (
	UtilSuite    struct{}
	TestStringer struct{}
)

func (TestStringer) String() string {
	return "bar"
}

func (s *UtilSuite) TestSerializeKey(t sweet.T) {
	Expect(serializeKey("foo")).To(Equal("foo"))
	Expect(serializeKey(TestStringer{})).To(Equal("bar"))
	Expect(serializeKey(TestConfigKey{})).To(Equal("TestConfigKey"))
	Expect(serializeKey(&TestConfigKey{})).To(Equal("TestConfigKey"))
}
