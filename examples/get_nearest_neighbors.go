package main

import (
	"fmt"
	"github.com/yarlett/corpustools"
	"time"
)

func main() {
	// Create a corpus object from the test corpus.
	corpus := corpustools.CorpusFromFile("../data/brown.txt", true, false)
	fmt.Println(corpus.Info())

	// Get the list of comparison terms.
	t1 := time.Now()
	seqs := make([][]int, 0)
	for order := 1; order <= 3; order++ {
		for _, ngram := range corpus.Ngrams(order) {
			if corpus.Frequency(ngram) >= 20 {
				seqs = append(seqs, ngram)
			}
		}
	}
	t2 := time.Now()
	fmt.Printf("%d sequences in nearest neighbor set (took %v).\n", len(seqs), t2.Sub(t1))

	// Compute and report the nearest neighbors.
	t1 = time.Now()
	for i := 0; i < 100; i++ {
		results := corpus.NearestNeighbors(seqs[i], seqs)
		fmt.Printf("Top 10 nearest neighbors of '%v' are...\n", corpus.ToString(seqs[i]))
		for j := 0; j < 10; j++ {
			fmt.Printf("'%v' score=%v\n", corpus.ToString(results[j].Seq), results[j].Val) 
		}
		fmt.Println()
	}
	t2 = time.Now()
	fmt.Printf("Took %v.\n", t2.Sub(t1))
}