package main

import (
	"fmt"
	"github.com/yarlett/corpustools"
	"log"
	"os"
	"time"
)

func main() {
	// Load the corpus as a lower-case sequence of characters.
	corpus := corpustools.CorpusFromFile("../data/test_corpus.txt", true, true)
	//corpus := corpustools.CorpusFromFile("/Users/yarlett/Corpora/Brown.txt", true)
	fmt.Println(corpus.Info())

	// Enumerate all the ngrams we want to explore.
	seqs := make([][]int, 0)
	for length := 3; length <= 10; length++ {
		ngrams := corpus.Ngrams(length)
		for _, ngram := range ngrams {
			if corpus.Frequency(ngram) >= 600 {
				seqs = append(seqs, ngram)
			}
		}
	}
	fmt.Printf("%d ngrams to be explored.\n", len(seqs))

	// Initialize the MDLSegmenter.
	mdlseg := corpustools.NewMDLSegmenter(corpus)

	// Initialize the output file.
	of, err := os.Create("mdl_segmentation_results.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer of.Close()
	of.Write([]byte("sequence,sequence_as_text,freq,dl_model,dl_data,dl_total\n"))

	// Compute baseline description length.
	dl_model_baseline, dl_data_baseline := mdlseg.DescriptionLength()
	of.Write([]byte(fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", "null", "null", 0, dl_model_baseline, dl_data_baseline, dl_model_baseline+dl_data_baseline)))

	// Identify the corpus subsequences which minimize the description length of the corpus.
	dlds := make(corpustools.Results, 0)
	for i, seq := range seqs {
		t := time.Now()
		mdlseg.AddNgram(seq)
		dl_model, dl_data := mdlseg.DescriptionLength()
		mdlseg.RemoveNgram(seq)
		dlds = append(dlds, corpustools.Result{Seq: seq, Val: dl_model + dl_data})
		// Report on performance.
		fmt.Printf("  Sequence %d/%d: %10v (%10v) --> %10.2f, %10.2f (took %v)\n", i+1, len(seqs), corpus.ToString(seq), seq, dl_model, dl_data, time.Now().Sub(t))
		// Write result to CSV file.
		of.Write([]byte(fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", seq, corpus.ToString(seq), corpus.Frequency(seq), dl_model, dl_data, dl_model+dl_data)))
	}
}