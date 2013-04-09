package main

import (
	"fmt"
	"github.com/yarlett/corpustools"
	"time"
)

func main() {
	// Create a corpus from a text file.
	corpus := corpustools.CorpusFromFile("../data/brown.txt", true, false)
	fmt.Println(corpus.Info())

	// Iterate over various orders and generate the ngrams of this length.
	for n := 1; n <= 10; n++ {
		t1 := time.Now()
		ngrams := corpus.Ngrams(n)
		t2 := time.Now()
		fmt.Printf("%d %dgrams found in %v.\n", len(ngrams), n, t2.Sub(t1))
	}

	// Report the frequencies of ngrams.
	for n := 1; n <= 10; n++ {
		ngrams := corpus.Ngrams(n)
		for j, ngram := range ngrams {
			fmt.Printf("%dgram %d = %v (%v) has frequency of %d.\n", n, j, corpus.ToString(ngram), ngram, corpus.Frequency(ngram))
		}
	}
}