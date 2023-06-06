package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	iforest "github.com/codegaudi/go-iforest"
	"github.com/twsnmp/golof/lof"
	"github.com/twsnmp/tfidf"
	"github.com/twsnmp/tfidf/seg"
	"github.com/twsnmp/tfidf/util"
)

var version = "1.0.0"
var commit = ""

var filter = ""
var count = 10
var tokenizer = "en"
var ad = ""
var userTG = false
var showVersion = false

type logEnt struct {
	Pos   int
	Score float64
}

func init() {
	flag.StringVar(&filter, "f", "", "regexp filter")
	flag.StringVar(&ad, "a", "", "anomaly detection(sum|if|lof)")
	flag.StringVar(&tokenizer, "t", "en", "tokenizer (en|log|ja)")
	flag.IntVar(&count, "c", 10, "show top n count")
	flag.BoolVar(&userTG, "g", false, "use time grinder")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()
}

func main() {
	if showVersion {
		fmt.Printf("version=%s(%s)", version, commit)
		return
	}
	if flag.NArg() < 1 {
		log.Fatalln("no input file")
	}
	lines := []string{}
	for i := 0; i < flag.NArg(); i++ {
		ls, err := util.ReadLines(flag.Arg(0), filter)
		if err != nil {
			log.Fatalln(err)
		}
		lines = append(lines, ls...)
	}
	if len(lines) < 1 {
		log.Fatalln("no lines")
	}
	logs := []logEnt{}
	var f *tfidf.TFIDF
	switch tokenizer {
	case "ja":
		f = tfidf.NewTokenizer(seg.NewJaTokenizer())
	case "log":
		f = tfidf.NewTokenizer(seg.NewLogTokenizer(userTG))
	default:
		f = tfidf.New()
	}
	f.AddDocs(lines...)
	switch ad {
	case "lof":
		vectors := f.GetTFIDF(lines...)
		samples := lof.GetSamplesFromFloat64s(vectors)
		lofGetter := lof.NewLOF(5)
		if err := lofGetter.Train(samples); err != nil {
			log.Fatalf("LOF err=%v", err)
		}
		for i, s := range samples {
			s := lofGetter.GetLOF(s, "fast")
			logs = append(logs, logEnt{
				Pos:   i,
				Score: s,
			})
		}
	case "if":
		vectors := f.GetTFIDF(lines...)
		rand.Seed(time.Now().UnixNano())
		ss := 256
		if ss > len(vectors) {
			ss = len(vectors)
		}
		ifo, err := iforest.NewIForest(vectors, 1000, ss)
		if err != nil {
			log.Fatalf("NewIForest err=%v", err)
		}
		for i, v := range vectors {
			s := ifo.CalculateAnomalyScore(v)
			logs = append(logs, logEnt{
				Pos:   i,
				Score: s,
			})
		}
	default:
		for i := 0; i < len(lines); i++ {
			w := f.Cal(lines[i])
			logs = append(logs, logEnt{
				Pos:   i,
				Score: sumTFIDF(w),
			})
		}
	}
	sort.Slice(logs, func(i, j int) bool { return logs[i].Score > logs[j].Score })
	fmt.Printf("%-5s\t%s\n", "Score", "Log")
	for i := 0; i < len(logs) && i < count; i++ {
		fmt.Printf("%5.3f\t%s\n", logs[i].Score, lines[logs[i].Pos])
	}
}

func sumTFIDF(m map[string]float64) float64 {
	r := 0.0
	for _, v := range m {
		r += v
	}
	return r
}
