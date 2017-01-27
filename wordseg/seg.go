package wordseg

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/chonla/go-trie/trie"
)

// IDict is dictionary interface
type IDict interface {
	LoadString(t string)
	LoadStringSet(ta []string)
	LoadFile(f string) error
	Has(v string) bool
	Clear()
	Depth() int
}

// Seg is segmentor
type Seg struct {
	Dict IDict
}

// NewSeg create a new Seg
func NewSeg(d IDict) *Seg {
	if d == nil {
		d = trie.NewDict()
	}
	return &Seg{Dict: d}
}

// UseDictFile to load dictionary file
func (s *Seg) UseDictFile(f string) error {
	return s.Dict.LoadFile(f)
}

// UseDictData to load dictionary from string set
func (s *Seg) UseDictData(ta []string) {
	s.Dict.LoadStringSet(ta)
}

// Clear to clean up dictionary
func (s *Seg) Clear() {
	s.Dict.Clear()
}

// SegmentText is to segment text into tokens
func (s *Seg) SegmentText(t string) []string {
	ts := s.groupText(t)
	out := []string{}

	for _, ti := range ts {
		if s.isThai(ti) {
			out = append(out, s.segmentThai(ti)...)
		} else {
			buf := strings.Trim(ti, " ")
			if len(buf) > 0 {
				out = append(out, buf)
			}
		}
	}

	return out
}

// groupText is to group text based by character type
func (s *Seg) groupText(t string) []string {
	b := []rune(t)
	out := []string{}
	var buf bytes.Buffer
	isthai := false
	for _, a := range b {
		c := string(a)
		if s.isThai(c) {
			if isthai {
				buf.WriteString(c)
			} else {
				out = append(out, buf.String())
				buf.Reset()
				buf.WriteString(c)
				isthai = true
			}
		} else {
			if isthai {
				out = append(out, buf.String())
				buf.Reset()
				buf.WriteString(c)
				isthai = false
			} else {
				buf.WriteString(c)
			}
		}
	}
	if buf.Len() > 0 {
		out = append(out, buf.String())
	}
	return out
}

// segmentThai is to segment a text containing non whitespace
func (s *Seg) segmentThai(t string) []string {
	b := []rune(t)
	var buf bytes.Buffer
	bufsize := 0
	out := []string{}
	l := len(b)
	recentCheckpoint := 0
	lastCheckpoint := 0
	depth := s.Dict.Depth()
	for cursor := 0; cursor < l; cursor++ {
		c := string(b[cursor])

		buf.WriteString(c)
		bufsize++

		if bufsize <= depth {
			if s.Dict.Has(buf.String()) {
				recentCheckpoint = cursor
			}

			if cursor == l-1 {
				goto flushBuffer
			}

			continue
		}

	flushBuffer:
		w := string(b[lastCheckpoint : recentCheckpoint+1])
		out = append(out, w)
		buf.Reset()
		bufsize = 0
		cursor = recentCheckpoint
		lastCheckpoint = recentCheckpoint + 1
	}

	return out
}

func (s *Seg) isThai(t string) bool {
	rt := []rune(t)
	for _, r := range rt {
		if !unicode.In(r, unicode.Thai) {
			return false
		}
	}
	return true
}
