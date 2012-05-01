package main

import (
	"corpustools"
	"fmt"
	"time"
)

func main() {
	// Create a corpus from a file.
	corpus := corpustools.CorpusFromFile("/Users/dan/github/exponential_manifold_embedding/data/brown.txt", true)
	fmt.Println(corpus.Info())
	for order := 1; order < 10; order++ {
		t1 := time.Now()
		ngrams := corpus.Ngrams(order)
		t2 := time.Now()
		//fmt.Printf("%v", ngrams)
		fmt.Printf("%d %dgrams found in %v.\n", len(ngrams), order, t2.Sub(t1))
	}
}