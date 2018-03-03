package log

import (
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&LoggerSuite{})
		s.AddSuite(&CallerSuite{})
		s.AddSuite(&ConfigSuite{})
		s.AddSuite(&GomolJSONSuite{})
		s.AddSuite(&ReplaySuite{})
		s.AddSuite(&RollupSuite{})
	})
}

//
// Mocks

type testShim struct {
	messages []*logMessage
}

func (ts *testShim) WithFields(fields Fields) logShim {
	return ts
}

func (ts *testShim) Log(level LogLevel, format string, args ...interface{}) {
	ts.LogWithFields(level, nil, format, args...)
}

func (ts *testShim) LogWithFields(level LogLevel, fields Fields, format string, args ...interface{}) {
	ts.messages = append(ts.messages, &logMessage{
		level:  level,
		fields: fields,
		format: format,
		args:   args,
	})
}

func (ts *testShim) Sync() error {
	return nil
}
