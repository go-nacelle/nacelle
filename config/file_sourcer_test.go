package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FileSourcerSuite struct{}

func (s *FileSourcerSuite) TestLoadJSON(t sweet.T) {
	sourcer, err := NewFileSourcer("test-files/values.json", ParseYAML)
	Expect(err).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadJSONNoParser(t sweet.T) {
	sourcer, err := NewFileSourcer("test-files/values.json", nil)
	Expect(err).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadYAML(t sweet.T) {
	sourcer, err := NewFileSourcer("test-files/values.yaml", ParseYAML)
	Expect(err).To(BeNil())
	testFileSourcer(sourcer)
}
func (s *FileSourcerSuite) TestLoadYAMLNoParser(t sweet.T) {
	sourcer, err := NewFileSourcer("test-files/values.yaml", nil)
	Expect(err).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadTOML(t sweet.T) {
	sourcer, err := NewFileSourcer("test-files/values.toml", ParseTOML)
	Expect(err).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadTOMLNoParser(t sweet.T) {
	sourcer, err := NewFileSourcer("test-files/values.toml", nil)
	Expect(err).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestOptionalFileSourcer(t sweet.T) {
	sourcer, err := NewOptionalFileSourcer("test-files/no-such-file.json", nil)
	Expect(err).To(BeNil())
	ensureMissing(sourcer, []string{"foo"})
}

func testFileSourcer(sourcer Sourcer) {
	ensureEquals(sourcer, []string{"foo"}, "bar")
	ensureMatches(sourcer, []string{"bar"}, "[1, 2, 3]")
	ensureMatches(sourcer, []string{"bonk"}, `{"x": 1, "y": 2, "z": 3}`)
	ensureMatches(sourcer, []string{"encoded"}, `{"w": 4}`)
	ensureMatches(sourcer, []string{"bonk.x"}, `1`)
	ensureMatches(sourcer, []string{"encoded.w"}, `4`)
	ensureMatches(sourcer, []string{"deeply.nested.struct"}, `[1, 2, 3]`)
}
