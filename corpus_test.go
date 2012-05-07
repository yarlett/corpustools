package corpustools

import (
	"os"
	"strings"
	"testing"
)

// Iterate up to trigrams for test purposes.
var (
	MAX_NGRAM_LENGTH int = 3
)

// Load a corpus on which to perform testing and benchmarking.
var path, _ = os.Getwd()
var corpus = CorpusFromFile(strings.Join([]string{path, "/data/test_corpus.txt"}, ""), true)

// Length of corpus and suffix array should be the same.
func TestBasics(t *testing.T) {
	// Length of corpus and suffix array should be the same.
	if len(corpus.seq) != len(corpus.sfx) {
		t.Errorf("Corpus sequence and suffix arrays are not the same length (%d vs. %d)!", len(corpus.seq), len(corpus.sfx))
	}
	// i + 1th suffix ngram should be >= ith suffix ngram.
	for spos := 0; spos < len(corpus.sfx)-1; spos++ {
		if SeqCmp(corpus.seq[corpus.sfx[spos]:], corpus.seq[corpus.sfx[spos+1]:]) == 1 {
			t.Errorf("Suffix ordering error detected at positions %d and %d!\n", spos, spos+1)
		}
	}
	// Number of unigrams should equal the size of the vocabulary.
	unigrams := corpus.Ngrams(1)
	if len(unigrams) != len(corpus.voc) {
		t.Errorf("Number of unigrams is not equal to the size of the vocabulary (%d vs. %d)!", len(unigrams), len(corpus.voc))
	}
	// Check that suffix indices returned by fast search matches those returned by slow search.
	for _, unigram := range unigrams {
		slo_slow, shi_slow := corpus.slowSearch(unigram)
		slo_fast, shi_fast := corpus.SuffixSearch(unigram)
		if (slo_slow != slo_fast) || (shi_slow != shi_fast) {
			t.Errorf("Error in finding suffix indices! %v: (%d, %d) vs. (%d, %d).\n", unigram, slo_slow, shi_slow, slo_fast, shi_fast)
		}
	}
}

// Benchmark for making a corpus from a text file.
func BenchmarkCorpus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CorpusFromFile(strings.Join([]string{path, "/data/test_corpus.txt"}, ""), true)
	}
}

// Benchmark for generating ngrams in the corpus up to a specified length.
func BenchmarkNgrams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for length := 1; length < MAX_NGRAM_LENGTH; length++ {
			_ = corpus.Ngrams(length)
		}
	}
}

// Benchmark for generating frequencies of ngrams in the corpus up to a specified length.
func BenchmarkFreqs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for length := 1; length < MAX_NGRAM_LENGTH; length++ {
			ngrams := corpus.Ngrams(length)
			for _, ngram := range ngrams {
				_ = corpus.Frequency(ngram)
			}
		}
	}
}
