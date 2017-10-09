package log

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ConfigSuite struct{}

func (s *ConfigSuite) TestIsLegalBackend(t sweet.T) {
	Expect(isLegalBackend("gomol")).To(BeTrue())
	Expect(isLegalBackend("logrus")).To(BeTrue())
	Expect(isLegalBackend("zap")).To(BeTrue())
	Expect(isLegalBackend("gomolx")).To(BeFalse())
	Expect(isLegalBackend("paz")).To(BeFalse())
}

func (s *ConfigSuite) TestIsLegalLevel(t sweet.T) {
	Expect(isLegalLevel("debug")).To(BeTrue())
	Expect(isLegalLevel("info")).To(BeTrue())
	Expect(isLegalLevel("warning")).To(BeTrue())
	Expect(isLegalLevel("error")).To(BeTrue())
	Expect(isLegalLevel("fatal")).To(BeTrue())
	Expect(isLegalLevel("warn")).To(BeFalse())
	Expect(isLegalLevel("trace")).To(BeFalse())
	Expect(isLegalLevel("die")).To(BeFalse())
}

func (s *ConfigSuite) TestIsLegalEncoding(t sweet.T) {
	Expect(isLegalEncoding("json")).To(BeTrue())
	Expect(isLegalEncoding("console")).To(BeTrue())
	Expect(isLegalEncoding("file")).To(BeFalse())
	Expect(isLegalEncoding("yaml")).To(BeFalse())
}
