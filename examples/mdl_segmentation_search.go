package main

import (
	"corpustools"
	"fmt"
	"sort"
	"time"
)

func main() {
	// Load the corpus.
	corpus := corpustools.CorpusFromFile("/Users/dan/github/exponential_manifold_embedding/data/brown.txt", true)
	//corpus := corpustools.CorpusFromFile("/Users/yarlett/Corpora/Brown.txt", true)
	fmt.Println(corpus.Info())

	// Initialize the MDLSegmenter.
	mdlseg := corpustools.NewMDLSegmenter(corpus)

	// Enumerate all the ngrams we want to explore.
	seqs := make([][]int, 0)
	for length := 2; length <= 5; length++ {
		ngrams := corpus.Ngrams(length)
		for _, ngram := range ngrams {
			if corpus.Frequency(ngram) >= 30 {
				seqs = append(seqs, ngram)
			}
		}
	}
	fmt.Printf("%d ngrams to be explored.\n", len(seqs))

	// Compute baseline description length.
	description_length_model_baseline, description_length_data_baseline := mdlseg.DescriptionLength()

	// Identify the corpus subsequences which minimize the description length of the corpus.
	dlds := make(corpustools.Results, 0)
	for i, seq := range seqs {
		t := time.Now()
		mdlseg.AddNgram(seq)
		dl_model, dl_data := mdlseg.DescriptionLength()
		mdlseg.RemoveNgram(seq)
		dlds = append(dlds, corpustools.Result{Seq: seq, Val: dl_model + dl_data})
		fmt.Printf("  Sequence %d/%d: %10v (%10v) --> %10.2f, %10.2f (took %v)\n", i+1, len(seqs), corpus.ToString(seq), seq, dl_model, dl_data, time.Now().Sub(t))
	}
	// Sort the description lengths in descending order.
	sort.Sort(dlds)
	// Print the results.
	fmt.Printf("\n")
	fmt.Printf("Baseline description length is %.2f + %.2f = %.2f.\n", description_length_model_baseline, description_length_data_baseline, description_length_model_baseline + description_length_data_baseline)
	for i := 0; i < 100; i++ {
		fmt.Printf("Sequence %d: %10v %10v: f=%6d, description_length = %.2f.\n", i, corpus.ToString(dlds[i].Seq), dlds[i].Seq, corpus.Frequency(dlds[i].Seq), dlds[i].Val)
	}
}