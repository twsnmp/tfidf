package tfidf

import (
	"testing"

	"github.com/twsnmp/tfidf/seg"
)

func TestTfidf(t *testing.T) {
	f := New()
	f.AddDocs("how are you", "are you fine", "how old are you", "are you ok", "i am ok", "i am file")

	text := "it is so cool"
	w := f.Cal(text)
	if p, ok := w["cool"]; !ok || p != 0.4864775372638283 {
		t.Fatalf("failed en test w=%v", w)
	}
	f = NewTokenizer(&seg.JaTokenizer{})
	f.AddDocs("寿司が食べたい。", "カレーが食べたくない。", "焼肉が食べたい。")

	text = "ラーメンが食べたい。"
	w = f.Cal(text)
	if p, ok := w["ラーメン"]; !ok || p != 0.2772588722239781 {
		t.Fatalf("failed ja test w=%v", w)
	}
}
