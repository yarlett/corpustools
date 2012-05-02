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

		//
		t3 := time.Now()
		for i, ngram := range(ngrams) {
			f := corpus.Frequency(ngram)
		t4 := time.Now()

		fmt.Printf("%d %dgrams found in %v.\n", len(ngrams), order, t2.Sub(t1))

		
		// for i, ngram := range(ngrams) {
		// 	l := corpus.DescriptionLengthDelta(ng)
		// 	fmt.Printf("(%d/%d) %v: %.2f\n", i + 1, len(ngrams), ng, l)
		// }

		// for _, ngram := range(ngrams) {
		// 	slo, shi := corpus.Search(ngram)
		// 	fmt.Printf("%v --> (%d, %d)\n", ngram, slo, shi)
		// }
	}
}