package tfidf

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"sort"

	"github.com/twsnmp/tfidf/seg"
	"github.com/twsnmp/tfidf/util"
)

// TFIDF tfidf model
type TFIDF struct {
	docIndex  map[string]int         // train document index in TermFreqs
	termFreqs []map[string]int       // term frequency for each train document
	termDocs  map[string]int         // documents number for each term in train data
	n         int                    // number of documents in train data
	dup       int                    // number of dup documents
	stopWords map[string]interface{} // words to be filtered
	tokenizer seg.Tokenizer          // tokenizer, space is used as default
	allTerms  []string               // all terms in train data
}

// New new model with default
func New() *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		dup:       0,
		tokenizer: &seg.EnTokenizer{},
	}
}

// NewTokenizer new with specified tokenizer
func NewTokenizer(tokenizer seg.Tokenizer) *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		dup:       0,
		tokenizer: tokenizer,
		allTerms:  []string{},
	}
}

// AddStopWords add stop words to be filtered
func (f *TFIDF) AddStopWords(words ...string) {
	if f.stopWords == nil {
		f.stopWords = make(map[string]interface{})
	}

	for _, word := range words {
		f.stopWords[word] = nil
	}
}

// AddStopWordsFile add stop words file to be filtered, with one word a line
func (f *TFIDF) AddStopWordsFile(file string) (err error) {
	lines, err := util.ReadLines(file, "")
	if err != nil {
		return
	}
	f.AddStopWords(lines...)
	return
}

// AddDocs add train documents
func (f *TFIDF) AddDocs(docs ...string) {
	for _, doc := range docs {
		h := hash(doc)
		docPos := f.docHashPos(h)
		if docPos >= 0 {
			termFreq := f.termFreqs[docPos]
			for term := range termFreq {
				f.termDocs[term]++
			}
			f.dup++
			continue
		}
		termFreq := f.termFreq(doc)
		if len(termFreq) == 0 {
			continue
		}

		f.docIndex[h] = f.n
		f.n++

		f.termFreqs = append(f.termFreqs, termFreq)

		for term := range termFreq {
			f.termDocs[term]++
		}
	}
	f.allTerms = []string{}
	for t := range f.termDocs {
		f.allTerms = append(f.allTerms, t)
	}
	sort.Strings(f.allTerms)
}

// Cal calculate tf-idf weight for specified document
func (f *TFIDF) Cal(doc string) (weight map[string]float64) {
	weight = make(map[string]float64)

	var termFreq map[string]int

	docPos := f.docPos(doc)
	o := 0
	if docPos < 0 {
		termFreq = f.termFreq(doc)
		o = 1
	} else {
		termFreq = f.termFreqs[docPos]
	}

	docTerms := 0
	for _, freq := range termFreq {
		docTerms += freq
	}
	for term, freq := range termFreq {
		weight[term] = tfidf(freq, docTerms, f.termDocs[term], f.n, o)
	}

	return weight
}

// GetTFIDF : calculate tf-idf vector for specified documents
func (f *TFIDF) GetTFIDF(limit int, docs ...string) [][]float64 {
	ret := [][]float64{}
	for _, doc := range docs {
		var termFreq map[string]int
		docPos := f.docPos(doc)
		o := 0
		if docPos < 0 {
			termFreq = f.termFreq(doc)
			o = 1
		} else {
			termFreq = f.termFreqs[docPos]
		}
		docTerms := 0
		for _, freq := range termFreq {
			docTerms += freq
		}
		w := make(map[string]float64)
		for term, freq := range termFreq {
			w[term] = tfidf(freq, docTerms, f.termDocs[term], f.n, o)
		}
		vector := []float64{}
		for _, t := range f.allTerms {
			if v, ok := w[t]; ok {
				vector = append(vector, v)
			} else {
				vector = append(vector, 0.0)
			}
		}
		if len(vector) < limit {
			ret = append(ret, vector)
		} else {
			rv := make([]float64, limit)
			for i, e := range vector {
				rv[i%limit] += e
			}
			ret = append(ret, rv)
		}
	}
	return ret
}

// GetDocumentCount : get number of documents in train data
func (f *TFIDF) GetDocumentCount() int {
	return f.n + f.dup
}

// GetSDupCount : get number of duplicate documents in train data
func (f *TFIDF) GetDupCount() int {
	return f.dup
}

// GetAllTerms : get sorted all terms in train data
func (f *TFIDF) GetAllTerms() []string {
	return f.allTerms
}

func (f *TFIDF) termFreq(doc string) (m map[string]int) {
	m = make(map[string]int)

	tokens := f.tokenizer.Seg(doc)
	if len(tokens) == 0 {
		return
	}

	for _, term := range tokens {
		if _, ok := f.stopWords[term]; ok {
			continue
		}

		m[term]++
	}

	return
}

func (f *TFIDF) docHashPos(hash string) int {
	if pos, ok := f.docIndex[hash]; ok {
		return pos
	}

	return -1
}

func (f *TFIDF) docPos(doc string) int {
	return f.docHashPos(hash(doc))
}

func hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func tfidf(termFreq, docTerms, termDocs, N, o int) float64 {
	tf := float64(termFreq) / float64(docTerms)
	idf := math.Log(float64(o+N)/float64(o+termDocs)) + 1
	return tf * idf
}
