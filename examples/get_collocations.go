package main

import (
	"fmt"
	"github.com/yarlett/corpustools"
	"sort"
)

func main() {
	// Create a corpus from a text file.
	corpus := corpustools.CorpusFromFile("../data/test_corpus.txt", true, false)
	fmt.Println(corpus.Info())

	// Compute the mutual information associated with ngrams of varying length.
	min_freq := 2
	for n := 2; n <= 4; n++ {
		results := make(corpustools.Results, 0)
		for _, ngram := range corpus.Ngrams(n) {
			if corpus.Frequency(ngram) >= min_freq {
				I := corpus.MutualInformation(ngram)
				results = append(results, corpustools.Result{ngram, I})
			}
		}
		sort.Sort(corpustools.ResultsReverseSort{results})
		fmt.Printf("%dgrams with the highest mutual information:\n", n)
		for i, result := range results {
			fmt.Printf("%d: %v (%v) --> %v\n", i+1, corpus.ToString(result.Seq), result.Seq, result.Val)
		}
		fmt.Printf("\n")
	}
}