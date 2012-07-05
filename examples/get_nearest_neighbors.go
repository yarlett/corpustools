package main

import (
	"corpustools"
	"fmt"
	"time"
)

func main() {
	// Create a corpus object from the test corpus.
	corpus := corpustools.CorpusFromFile("/Users/yarlett/Corpora/Brown.txt", true)
	fmt.Println(corpus.Info())

	// Get the list of comparison terms.
	t1 := time.Now()
	seqs := make([][]int, 0)
	for order := 1; order <= 3; order++ {
		for _, ngram := range corpus.Ngrams(order) {
			if corpus.Frequency(ngram) >= 30 {
				seqs = append(seqs, ngram)
			}
		}
	}
	t2 := time.Now()
	fmt.Printf("%d sequences in nearest neighbor set (took %v).\n", len(seqs), t2.Sub(t1))

	// Compute and report the nearest neighbors.
	t1 = time.Now()
	for i := 0; i < 100; i++ {
		_ = corpus.NearestNeighbors(seqs[i], seqs)
	}
	t2 = time.Now()
	// for i := 0; i < 15; i++ {
	// 	fmt.Printf("Sequence %d: %v %v: f=%d, description_length_delta = %.3f.\n", i, corpus.ToString(nns[i].Seq), nns[i].Seq, corpus.Frequency(nns[i].Seq), nns[i].Val)
	// }
	fmt.Printf("Took %v.\n", t2.Sub(t1))
}