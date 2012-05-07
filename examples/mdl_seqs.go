package main

import (
	"corpustools"
	"fmt"
	//"os"
	//"strings"
	"sort"
)

func main() {
	// // Get the path for the test corpus.
	// var path, _ = os.Getwd()
	// path_parts := strings.Split(path, "/")
	// path_parts = path_parts[: len(path_parts) - 1]
	// for _, part := range []string{"data", "test_corpus.txt"} {
	// 	path_parts = append(path_parts, part)
	// }
	// corpusfile := strings.Join(path_parts, "/")

	// // Create a corpus object from the test corpus.
	// lower_case_tokens := true
	// corpus := corpustools.CorpusFromFile(corpusfile, lower_case_tokens)
	// fmt.Println(corpus.Info())

	corpus := corpustools.CorpusFromFile("/Users/yarlett/Corpora/Brown.txt", true)
	fmt.Println(corpus.Info())

	// Enumerate all the subsequences we want to explore.
	seqs := make([][]int, 0)
	for length := 2; length <= 5; length++ {
		ngrams := corpus.Ngrams(length)
		for _, ngram := range ngrams {
			if corpus.Frequency(ngram) >= 5 {
				seqs = append(seqs, ngram)
			}
		}
	}
	fmt.Printf("%d sequences to be explored.\n", len(seqs))

	// Identify the corpus subsequences which minimize the description length of the corpus.
	dlds := make(corpustools.Results, 0)
	for i, seq := range seqs {
		dld := corpus.DescriptionLengthDelta(seq)
		dlds = append(dlds, corpustools.Result{Seq: seq, Val: dld})
		// if dld > 0.0 {
		// 	fmt.Printf("Sequence %d/%d: %v %v: description_length_delta = %.3f.\n", i + 1, len(seqs), corpus.ToString(seq), seq, dld)
		// }
		if i > 0 && i%1000 == 0 {
			fmt.Printf("  %d seqs processed...\n", i)
		}
	}
	sort.Sort(dlds)

	for i := 0; i < 15; i++ {
		fmt.Printf("Sequence %d: %v %v: f=%d, description_length_delta = %.3f.\n", i, corpus.ToString(dlds[i].Seq), dlds[i].Seq, corpus.Frequency(dlds[i].Seq), dlds[i].Val)
	}
}