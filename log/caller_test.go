package log

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"
)

type CallerSuite struct{}

var (
	testFields1 = Fields{"A": 1}
	testFields2 = Fields{"B": 2}
	testFields3 = Fields{"C": 3}
)

func (s *CallerSuite) testBasic(init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		Expect(err).To(BeNil())

		logger.Info("X")
		logger.InfoWithFields(Fields{"empty": false}, "Y")
		logger.Info("Z")
		logger.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	Expect(lines).To(HaveLen(3))

	var (
		data1 = Fields{}
		data2 = Fields{}
		data3 = Fields{}
	)

	Expect(json.Unmarshal([]byte(lines[0]), &data1)).To(BeNil())
	Expect(json.Unmarshal([]byte(lines[1]), &data2)).To(BeNil())
	Expect(json.Unmarshal([]byte(lines[2]), &data3)).To(BeNil())

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 27

	Expect(data1["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start+0)))
	Expect(data2["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start+1)))
	Expect(data3["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start+2)))
}

func (s *CallerSuite) testReplay(init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		Expect(err).To(BeNil())

		// Non-replayed messages are below log level - not emitted
		adapter := NewReplayAdapter(logger, LevelDebug, LevelInfo)
		adapter.Debug("X")
		adapter.InfoWithFields(Fields{"empty": false}, "Y")
		adapter.Debug("Z")
		adapter.Replay(LevelWarning)
		adapter.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	Expect(lines).To(HaveLen(4))

	var (
		data1 = Fields{}
		data2 = Fields{}
		data3 = Fields{}
	)

	Expect(json.Unmarshal([]byte(lines[1]), &data1)).To(BeNil())
	Expect(json.Unmarshal([]byte(lines[2]), &data2)).To(BeNil())
	Expect(json.Unmarshal([]byte(lines[3]), &data3)).To(BeNil())

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 63

	Expect(data1["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start+0)))
	Expect(data2["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start+1)))
	Expect(data3["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start+2)))
}

func (s *CallerSuite) testRollup(init func(*Config) (Logger, error)) {
	stderr := captureStderr(func() {
		logger, err := init(&Config{LogLevel: "info", LogEncoding: "json"})
		Expect(err).To(BeNil())

		clock := glock.NewMockClock()
		adapter := adaptShim(newRollupShim(logger, clock, time.Second))
		adapter.Info("A")
		adapter.Info("A")
		adapter.Info("A")
		adapter.Info("A")
		adapter.Info("A")
		clock.BlockingAdvance(time.Second)
		adapter.Sync()
	})

	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	Expect(lines).To(HaveLen(2))

	var (
		data1 = Fields{}
		data2 = Fields{}
	)

	Expect(json.Unmarshal([]byte(lines[0]), &data1)).To(BeNil())
	Expect(json.Unmarshal([]byte(lines[1]), &data2)).To(BeNil())

	// Note: this value refers to the line number containing `logger.Info("X")` in
	// the function literal above. If code is added before that line, this value
	// must be updated.
	start := 100

	Expect(data1["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start)))
	Expect(data2["caller"]).To(Equal(fmt.Sprintf("log/caller_test.go:%d", start)))
}

func (s *CallerSuite) testFields(init func(*Config) (Logger, error)) {
	s.testBasic(func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return logger.WithFields(testFields1).WithFields(testFields2).WithFields(testFields3), nil
	})
}

func (s *CallerSuite) testReplayAdapter(init func(*Config) (Logger, error)) {
	s.testBasic(func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return NewReplayAdapter(NewReplayAdapter(NewReplayAdapter(logger))), nil
	})
}

func (s *CallerSuite) testRollupAdapter(init func(*Config) (Logger, error)) {
	s.testBasic(func(config *Config) (Logger, error) {
		logger, err := init(config)
		if err != nil {
			return nil, err
		}

		return NewRollupAdapter(NewRollupAdapter(NewRollupAdapter(logger, time.Second), time.Second), time.Second), nil
	})
}

//
// Real Tests

func (s *CallerSuite) TestTrimPath(t sweet.T) {
	Expect(trimPath("")).To(Equal(""))
	Expect(trimPath("/")).To(Equal("/"))
	Expect(trimPath("/foo")).To(Equal("/foo"))
	Expect(trimPath("/foo/bar")).To(Equal("foo/bar"))
	Expect(trimPath("/foo/bar/baz")).To(Equal("bar/baz"))
	Expect(trimPath("/foo/bar/baz/bonk")).To(Equal("baz/bonk"))
}

func (s *CallerSuite) TestGomol(t sweet.T)                   { s.testBasic(InitGomolShim) }
func (s *CallerSuite) TestLogrus(t sweet.T)                  { s.testBasic(InitLogrusShim) }
func (s *CallerSuite) TestZap(t sweet.T)                     { s.testBasic(InitZapShim) }
func (s *CallerSuite) TestGomolWithFields(t sweet.T)         { s.testFields(InitGomolShim) }
func (s *CallerSuite) TestLogrusWithFields(t sweet.T)        { s.testFields(InitLogrusShim) }
func (s *CallerSuite) TestZapWithFields(t sweet.T)           { s.testFields(InitZapShim) }
func (s *CallerSuite) TestGomolWithReplayAdapter(t sweet.T)  { s.testReplayAdapter(InitGomolShim) }
func (s *CallerSuite) TestLogrusWithReplayAdapter(t sweet.T) { s.testReplayAdapter(InitLogrusShim) }
func (s *CallerSuite) TestZapWithReplayAdapter(t sweet.T)    { s.testReplayAdapter(InitZapShim) }
func (s *CallerSuite) TestGomolWithRollupAdapter(t sweet.T)  { s.testRollupAdapter(InitGomolShim) }
func (s *CallerSuite) TestLogrusWithRollupAdapter(t sweet.T) { s.testRollupAdapter(InitLogrusShim) }
func (s *CallerSuite) TestZapWithRollupAdapter(t sweet.T)    { s.testRollupAdapter(InitZapShim) }
func (s *CallerSuite) TestGomolReplay(t sweet.T)             { s.testReplay(InitGomolShim) }
func (s *CallerSuite) TestLogrusReplay(t sweet.T)            { s.testReplay(InitLogrusShim) }
func (s *CallerSuite) TestZapReplay(t sweet.T)               { s.testReplay(InitZapShim) }
func (s *CallerSuite) TestGomolRollup(t sweet.T)             { s.testRollup(InitGomolShim) }
func (s *CallerSuite) TestLogrusRollup(t sweet.T)            { s.testRollup(InitLogrusShim) }
func (s *CallerSuite) TestZapRollup(t sweet.T)               { s.testRollup(InitZapShim) }
