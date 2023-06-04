package seg

import (
	"strings"

	"github.com/gravwell/gravwell/v3/timegrinder"
)

type LogTokenizer struct {
	tg *timegrinder.TimeGrinder
	r  *strings.Replacer
}

func NewLogTokenizer(useTimeGrinder bool) *LogTokenizer {
	var tg *timegrinder.TimeGrinder
	if useTimeGrinder {
		tg, _ = timegrinder.New(timegrinder.Config{
			EnableLeftMostSeed: true,
		})
	}
	r := &LogTokenizer{
		r: strings.NewReplacer(
			"\"", " ", "'", " ", "(", " ", ")", " ",
			"{", " ", "}", " ", "[", " ", "]", " ",
			":", " ", ";", " ", "<", " ", ">", " ", "->", " ", ",", " ", "=", " ",
		),
		tg: tg,
	}
	return r
}

func (s *LogTokenizer) Seg(text string) []string {
	if s.tg != nil {
		if start, end, ok := s.tg.Match([]byte(text)); ok {
			text = text[:start] + text[end:]
		}
	}
	return strings.Fields(s.r.Replace(text))
}

func (s *LogTokenizer) Free() {

}
