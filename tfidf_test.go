package tfidf

import (
	"fmt"
	"testing"

	"github.com/twsnmp/tfidf/seg"
)

func TestCal(t *testing.T) {
	f := New()
	f.AddDocs("how are you", "are you fine", "how old are you", "are you ok", "i am ok", "i am file")

	text := "it is so cool"
	w := f.Cal(text)
	if p, ok := w["cool"]; !ok || p != 0.7364775372638284 {
		t.Fatalf("failed en test w=%v", w)
	}
	f = NewTokenizer(seg.NewJaTokenizer())
	f.AddDocs("寿司が食べたい。", "カレーが食べたくない。", "焼肉が食べたい。")

	text = "ラーメンが食べたい。"
	w = f.Cal(text)
	if p, ok := w["ラーメン"]; !ok || p != 0.47725887222397817 {
		t.Fatalf("failed ja test w=%v", w)
	}
}
func TestGetTFIDF(t *testing.T) {
	docs := []string{
		"犬 可愛い 犬 大きい",
		"猫 小さい 猫 可愛い 可愛い",
		"虫 小さい 可愛くない",
	}
	f := New()
	f.AddDocs(docs...)
	dc := f.GetDocumentCount()
	if dc != 3 {
		t.Fatalf("failed GetDocumentCount=%v", dc)
	}
	dup := f.GetDupCount()
	if dup != 0 {
		t.Fatalf("failed GetDupCount=%v", dc)
	}
	at := f.GetAllTerms()
	if len(at) != 7 || at[6] != "虫" {
		t.Fatalf("failed GetAllTerms=%v", at)
	}
	v := f.GetTFIDF(docs...)
	if len(v) != 3 || len(v[0]) != 7 || v[0][0] != 0.3513662770270411 {
		t.Fatalf("failed GetTFIDF=%v", v)
	}
	fmt.Println(f.GetAllTerms())
	for _, e := range v {
		fmt.Println(e)
	}
}
