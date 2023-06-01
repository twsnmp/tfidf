package seg

import (
	"github.com/ikawaha/kagome/tokenizer"
)

type JaTokenizer struct {
}

func (s *JaTokenizer) Seg(text string) []string {
	t := tokenizer.New()
	tokens := t.Tokenize(text)
	tokens = tokens[1 : len(tokens)-1]
	terms := make([]string, len(tokens))
	for i, token := range tokens {
		terms[i] = token.Surface
	}
	return terms
}

func (s *JaTokenizer) Free() {

}
