package corpustools

import (
	"os"
	"strings"
	"testing"
)

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
	// Number of unigrams should equal the size of the vocabulary.
	unigrams := corpus.Ngrams(1)
	if len(unigrams) != len(corpus.voc) {
		t.Errorf("Number of unigrams is not equal to the size of the vocabulary (%d vs, %d)!", len(unigrams), len(corpus.voc))
	}
	// Check that suffix indices returned by fast search matches those returned by slow search.
	for _, unigram := range(unigrams) {
		slo1, shi1 := corpus.SearchSlow(unigram)
		slo2, shi2 := corpus.Search(unigram)
		if (slo1 != slo2) || (shi1 != shi2) {
			t.Errorf("Error in finding suffix indices! %v: (%d, %d) vs. (%d, %d).\n", unigram, slo1, shi1, slo2, shi2)
		}
	}
}

// Test that all suffixes are in order.
func TestSuffixOrdering(t *testing.T) {
	for length := 0; length < MAX_NGRAM_LENGTH; length++ {
		for six := 0; six < len(corpus.sfx) - 1; six++ {
			if corpus.NgramSfxCmp(six, six + 1, length) == 1 { t.Errorf("Suffix ordering error detected at positions %d and %d!\n", six, six + 1) }
		}
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
func BenchmarkNgramFreqs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for length := 1; length < MAX_NGRAM_LENGTH; length++ {
			ngrams := corpus.Ngrams(length)
			for _, ngram := range(ngrams) {
				_ = corpus.Frequency(ngram)
			}
		}
	}
}