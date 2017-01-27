package wordseg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockObject struct {
	mock.Mock
}

func (m *MockObject) LoadFile(f string) error {
	m.Called(f)
	return nil
}

func (m *MockObject) LoadString(t string) {
	m.Called(t)
}

func (m *MockObject) LoadStringSet(ta []string) {
	m.Called(ta)
}

func (m *MockObject) Has(v string) bool {
	m.Called(v)
	return true
}

func (m *MockObject) Clear() {
	m.Called()
}

func (m *MockObject) Depth() int {
	m.Called()
	return 0
}

func TestAssignDictionaryFile(t *testing.T) {
	d := new(MockObject)
	d.On("LoadFile", mock.AnythingOfType("string"))

	f := "somefile"
	s := NewSeg(d)
	s.UseDictFile(f)

	d.AssertCalled(t, "LoadFile", f)
}

func TestCleanUp(t *testing.T) {
	d := new(MockObject)
	d.On("Clear")

	s := NewSeg(d)
	s.Clear()

	d.AssertCalled(t, "Clear")
}

func TestEmptyDictShouldReturnIdenticalStringInArray(t *testing.T) {
	s := NewSeg(nil)
	defer s.Clear()

	r := s.SegmentText("test")

	assert.Equal(t, []string{"test"}, r)
}

func Benchmark_Seg_isThai(b *testing.B) {
	word := "สวัสดี"

	s := NewSeg(nil)
	for i := 0; i < b.N; i++ {
		s.isThai(word)
	}
}
