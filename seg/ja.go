package seg

import (
	"github.com/ikawaha/kagome/tokenizer"
)

type JaTokenizer struct {
	t tokenizer.Tokenizer
}

func NewJaTokenizer() *JaTokenizer {
	return &JaTokenizer{
		t: tokenizer.New(),
	}
}

func (s *JaTokenizer) Seg(text string) []string {
	tokens := s.t.Tokenize(text)
	tokens = tokens[1 : len(tokens)-1]
	terms := make([]string, len(tokens))
	for i, token := range tokens {
		terms[i] = token.Surface
	}
	return terms
}

func (s *JaTokenizer) Free() {

}
