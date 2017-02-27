package wordseg

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/chonla/go-trie/trie"
)

const (
	// AlgoLongest is Longest Matching
	AlgoLongest = 0
	// AlgoMaximum is Maximum Matching
	AlgoMaximum = 1
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
	Algo int
}

// NewSeg create a new Seg
func NewSeg(d IDict) *Seg {
	if d == nil {
		d = trie.NewDict()
	}
	return &Seg{Dict: d, Algo: AlgoLongest}
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
			if s.Dict.Depth() > 0 {
				res := s.segmentThai(ti)
				out = append(out, res...)
			} else {
				out = append(out, ti)
			}
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
				if buf.Len() > 0 {
					out = append(out, buf.String())
					buf.Reset()
				}
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
	if s.Algo == AlgoLongest {
		return s.segmentThaiLongest(t)
	}
	if s.Algo == AlgoMaximum {
		return s.segmentThaiMaximum(t)
	}
	return []string{t}
}

// segmentThaiMaximum is to segment a text containing non whitespace
func (s *Seg) segmentThaiMaximum(t string) []string {
	/*
		b := []rune(t)
		var buf bytes.Buffer
		bufsize := 0
		out := []string{}
		l := len(b)
		depth := s.Dict.Depth()
		safe := []int{}

			for cursor := 0; cursor < l; cursor++ {
				c := string(b[cursor])

				buf.WriteString(c)
				bufsize++

				if bufsize <= depth {
					if s.Dict.Has(buf.String()) {
						safe = append(safe, cursor)
					} else {
					}
				} else {

				}
			}
		return out
	*/

	// Maximum matching is not ready yet
	return []string{t}
}

// segmentThaiLongest is to segment a text containing non whitespace
func (s *Seg) segmentThaiLongest(t string) []string {
	var buf bytes.Buffer
	bufsize := 0
	out := []string{}
	recentCheckpoint := 0
	lastCheckpoint := 0
	depth := s.Dict.Depth()

	cluster := s.createCluster(t)
	l := len(cluster)

	for cursor := 0; cursor < l; cursor++ {
		c := cluster[cursor]

		buf.WriteString(c)
		bufsize += len([]rune(c))

		// if bufsize > depth, it is 100% not in the dictionary
		if bufsize <= depth && s.Dict.Has(buf.String()) {
			recentCheckpoint = cursor
		}

		if cursor >= l-1 {
			w := ""
			if recentCheckpoint >= lastCheckpoint {
				w = strings.Join(cluster[lastCheckpoint:recentCheckpoint+1], "")
			} else {
				w = strings.Join(cluster[lastCheckpoint:], "")
				recentCheckpoint = cursor
			}
			out = append(out, w)
			buf.Reset()
			bufsize = 0
			cursor = recentCheckpoint
			lastCheckpoint = recentCheckpoint + 1
		}
	}

	return out
}

func (s *Seg) createCluster(t string) []string {
	// เ แ โ ใ ไ ไม้หันอากาศ
	needFollowingChar := []string{"\u0e40", "\u0e41", "\u0e42", "\u0e43", "\u0e44", "\u0e31"}
	// ะ ไม้หันอากาศ า ำ สระอิ สระอี สระอึ สระอือ สระอุ สระอู สระอาหางยาว ไม้ไต่คู้ ไม้เอก ไม้โท ไม้ตรี ไม้จัตวา ทันฑฆาต
	needPrecedingChar := []string{"\u0e30", "\u0e31", "\u0e32", "\u0e33", "\u0e34", "\u0e35", "\u0e36", "\u0e37", "\u0e38", "\u0e39", "\u0e45", "\u0e47", "\u0e48", "\u0e49", "\u0e4a", "\u0e4b", "\u0e4c"}

	b := []rune(t)
	var buf bytes.Buffer
	l := len(b)
	for cursor := 0; cursor < l; cursor++ {
		c := string(b[cursor])

		if s.isInGroup(c, needPrecedingChar) {
			buf.WriteString("<")
		}

		buf.WriteString(c)

		if !s.isInGroup(c, needFollowingChar) {
			buf.WriteString(" ")
		}
	}

	return strings.Split(strings.Replace(strings.Replace(buf.String(), " <", "", -1), "<", "", -1), " ")
}

func (s *Seg) isInGroup(c string, g []string) bool {
	for _, v := range g {
		if v == c {
			return true
		}
	}
	return false
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
