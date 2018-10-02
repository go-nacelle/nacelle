package config

//go:generate go-mockgen github.com/efritz/nacelle/config -i Config -o mock_config_test.go -f
//go:generate go-mockgen github.com/efritz/nacelle/logging -i Logger -o mock_logger_test.go -f

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/nacelle/config/tag"
	"github.com/efritz/nacelle/logging"
)

type LoggingConfigSuite struct{}

func (s *LoggingConfigSuite) TestLoadLogs(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &TestSimpleConfig{}
	)

	config.LoadFunc = func(target interface{}, modifiers ...tag.Modifier) error {
		target.(*TestSimpleConfig).X = "foo"
		target.(*TestSimpleConfig).Y = 123
		target.(*TestSimpleConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	}

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.InfoWithFieldsFuncCallCount()).To(Equal(1))

	params := logger.InfoWithFieldsFuncCallParams()[0]
	Expect(params.Arg1).To(Equal("Config loaded from environment"))
	Expect(params.Arg0).To(Equal(logging.Fields{
		"X": "foo",
		"Y": "123",
		"Q": `["bar","baz","bonk"]`,
	}))
}

func (s *LoggingConfigSuite) TestMask(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &TestMaskConfig{}
	)

	config.LoadFunc = func(target interface{}, modifiers ...tag.Modifier) error {
		target.(*TestMaskConfig).X = "foo"
		target.(*TestMaskConfig).Y = 123
		target.(*TestMaskConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	}

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.InfoWithFieldsFuncCallCount()).To(Equal(1))

	params := logger.InfoWithFieldsFuncCallParams()[0]
	Expect(params.Arg1).To(Equal("Config loaded from environment"))
	Expect(params.Arg0).To(Equal(logging.Fields{
		"X": "foo",
	}))
}

func (s *LoggingConfigSuite) TestBadMaskTag(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &TestBadMaskTagConfig{}
	)

	Expect(lc.Load(chunk)).To(MatchError("" +
		"failed to serialize config" +
		" (" +
		"field 'X' has an invalid mask tag" +
		")",
	))
}

func (s *LoggingConfigSuite) TestMustLoadLogs(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &TestSimpleConfig{}
	)

	lc.MustLoad(chunk)
	Expect(logger.InfoWithFieldsFuncCallCount()).To(Equal(1))
}
