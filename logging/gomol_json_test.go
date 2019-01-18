package logging

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aphistic/gomol"
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type GomolJSONSuite struct{}

func (s *GomolJSONSuite) TestInitLogger(t sweet.T) {
	logger := newJSONLogger(nil)
	Expect(logger.IsInitialized()).To(BeFalse())
	logger.InitLogger()
	Expect(logger.IsInitialized()).To(BeTrue())
}

func (s *GomolJSONSuite) TestShutdownLogger(t sweet.T) {
	logger := newJSONLogger(nil)
	logger.InitLogger()
	Expect(logger.IsInitialized()).To(BeTrue())
	logger.ShutdownLogger()
	Expect(logger.IsInitialized()).To(BeFalse())
}

func (s *GomolJSONSuite) TestLogm(t sweet.T) {
	var (
		logger    = newJSONLogger(nil)
		buffer    = bytes.NewBuffer(nil)
		timestamp = time.Unix(1503939881, 0)
	)

	logger.stream = buffer

	logger.Logm(
		timestamp,
		gomol.LevelFatal,
		Fields{"attr1": 4321},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(MatchJSON(fmt.Sprintf(`{
		"level": "fatal",
		"message": "test 1234",
		"timestamp": "%s",
		"attr1": 4321
	}`, timestamp.Format(JSONTimeFormat))))
}

func (s *GomolJSONSuite) TestBaseAttrs(t sweet.T) {
	var (
		logger    = newJSONLogger(nil)
		buffer    = bytes.NewBuffer(nil)
		base      = gomol.NewBase()
		timestamp = time.Unix(1503939881, 0)
	)

	base.SetAttr("attr1", 7890)
	base.SetAttr("attr2", "val2")

	logger.SetBase(base)
	logger.stream = buffer

	logger.Logm(
		timestamp,
		gomol.LevelDebug,
		Fields{
			"attr1": 4321,
			"attr3": "val3",
		},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(MatchJSON(fmt.Sprintf(`{
			"level": "debug",
			"message": "test 1234",
			"timestamp": "%s",
			"attr1": 4321,
			"attr2": "val2",
			"attr3": "val3"
		}`, timestamp.Format(JSONTimeFormat))))
}

func (s *GomolJSONSuite) TestCustomFieldNames(t sweet.T) {
	var (
		logger = newJSONLogger(map[string]string{
			"timestamp": "@timestamp",
			"level":     "log_level",
		})
		buffer    = bytes.NewBuffer(nil)
		timestamp = time.Unix(1503939881, 0)
	)

	logger.stream = buffer

	logger.Logm(
		timestamp,
		gomol.LevelFatal,
		Fields{"attr1": 4321},
		"test 1234",
	)

	Expect(string(buffer.Bytes())).To(MatchJSON(fmt.Sprintf(`{
		"log_level": "fatal",
		"message": "test 1234",
		"@timestamp": "%s",
		"attr1": 4321
	}`, timestamp.Format(JSONTimeFormat))))
}
