package log

import (
	"bytes"
	"time"

	"github.com/aphistic/gomol"
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type GomolJSONSuite struct{}

func (s *GomolJSONSuite) TestInitLogger(t sweet.T) {
	logger := newJSONLogger()
	Expect(logger.IsInitialized()).To(BeFalse())
	logger.InitLogger()
	Expect(logger.IsInitialized()).To(BeTrue())
}

func (s *GomolJSONSuite) TestShutdownLogger(t sweet.T) {
	logger := newJSONLogger()
	logger.InitLogger()
	Expect(logger.IsInitialized()).To(BeTrue())
	logger.ShutdownLogger()
	Expect(logger.IsInitialized()).To(BeFalse())
}

func (s *GomolJSONSuite) TestLogm(t sweet.T) {
	var (
		logger = newJSONLogger()
		buffer = bytes.NewBuffer(nil)
	)

	logger.stream = buffer

	logger.Logm(
		time.Unix(1503939881, 0),
		gomol.LevelFatal,
		map[string]interface{}{"attr1": 4321},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(MatchJSON(`{
		"level": "fatal",
		"message": "test 1234",
		"timestamp": "2017-08-28T12:04:41.000-0500",
		"attr1": 4321
	}`))
}

func (s *GomolJSONSuite) TestBaseAttrs(t sweet.T) {
	var (
		logger = newJSONLogger()
		buffer = bytes.NewBuffer(nil)
		base   = gomol.NewBase()
	)

	base.SetAttr("attr1", 7890)
	base.SetAttr("attr2", "val2")

	logger.SetBase(base)
	logger.stream = buffer

	logger.Logm(
		time.Unix(1503939881, 0),
		gomol.LevelDebug,
		map[string]interface{}{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(MatchJSON(`{
			"level": "debug",
			"message": "test 1234",
			"timestamp": "2017-08-28T12:04:41.000-0500",
			"attr1": 4321,
			"attr2": "val2",
			"attr3": "val3"
		}`))
}
