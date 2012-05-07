package main

import (
	"corpustools"
	"fmt"
	"time"
)

func main() {
	// Create a corpus from a text file.
	corpus := corpustools.CorpusFromFile("/Users/yarlett/Corpora/BNCAll.txt", true)
	fmt.Println(corpus.Info())

	// Iterate over various orders and generate the ngrams of this length.
	for length := 1; length <= 10; length++ {
		t1 := time.Now()
		ngrams := corpus.Ngrams(length)
		t2 := time.Now()
		fmt.Printf("%d %dgrams found in %v.\n", len(ngrams), length, t2.Sub(t1))
	}
}