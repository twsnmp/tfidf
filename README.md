# TFIDF

GO version of the TFIDF package

GO言語版のTFIDFパッケージ

[![Godoc Reference](https://godoc.org/github.com/twsnmp/tfidf?status.svg)](http://godoc.org/github.com/twsnmp/tfidf)
[![Go Report Card](https://goreportcard.com/badge/twsnmp/tfidf)](https://goreportcard.com/report/twsnmp/tfidf)

## Introduction

+ tokenizer support, contains english and Japanease Tokenizer.
+ TFIDF, calculate tfidf value of giving document.
+ Cosine, calculate Cosine value of giving documents pair.

+ 英語と日本語の単語分割に対応
+ 与えられたドキュメントのTFIDFを計算
+ コサイン類似度（Cosine Similarity）を計算

## Guide

```
go get github.com/twsnmp/tfidf
```


```go
package main

import (
	"fmt"

	"github.com/twsnmp/tfidf"
	"github.com/twsnmp/tfidf/seg"
	"github.com/twsnmp/tfidf/similarity"
)

func main() {

	f := tfidf.New()
	f.AddDocs("how are you", "are you fine", "how old are you", "are you ok", "i am ok", "i am file")

	t1 := "it is so cool"
	w1 := f.Cal(t1)
	fmt.Printf("weight of %s is %+v.\n", t1, w1)

	t2 := "you are so beautiful"
	w2 := f.Cal(t2)
	fmt.Printf("weight of %s is %+v.\n", t2, w2)

	sim := similarity.Cosine(w1, w2)
	fmt.Printf("cosine between %s and %s is %f .\n", t1, t2, sim)

	tokenizer := &seg.JaTokenizer{}
	defer tokenizer.Free()

	f = tfidf.NewTokenizer(tokenizer)

	f.AddDocs("寿司が食べたい。", "カレーが食べたくない。", "焼肉が食べたい。")

	t1 = "ラーメンが食べたい。"
	w1 = f.Cal(t1)
	fmt.Printf("weight of %s is %+v.\n", t1, w1)

	t2 = "ごはんが食べたい。"
	w2 = f.Cal(t2)
	fmt.Printf("weight of %s is %+v.\n", t2, w2)

	sim = similarity.Cosine(w1, w2)
	fmt.Printf("cosine between %s and %s is %f .\n", t1, t2, sim)
}

}

```
