package corpustools

import (
	"testing"
)

var (
	MAX_NGRAM_LENGTH int = 5
)

// Test that all suffixes are in order.
func TestSuffixOrdering(t *testing.T) {
	corpus := CorpusFromFile("/Users/dan/github/exponential_manifold_embedding/data/brown.txt", true)
	for six := 0; six < len(corpus.sfx) - 1; six++ {
		cmp := corpus.NgramSfxCmp(six, six + 1, 3)
		if cmp == 1 { t.Error(six, six + 1, cmp) }
	}
}

// Benchmark for generating ngrams in the corpus up to a specified length.
func BenchmarkNgrams(b *testing.B) {
	b.StopTimer()
	corpus := CorpusFromFile("/Users/dan/github/exponential_manifold_embedding/data/brown.txt", true)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for length := 1; length < MAX_NGRAM_LENGTH; length++ {
			_ = corpus.Ngrams(length)
		}
	}
}

// Benchmark for generating frequencies of ngrams in the corpus up to a specified length.
func BenchmarkNgramFreqs(b *testing.B) {
	b.StopTimer()
	corpus := CorpusFromFile("/Users/dan/github/exponential_manifold_embedding/data/brown.txt", true)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for length := 1; length < MAX_NGRAM_LENGTH; length++ {
			ngrams := corpus.Ngrams(length)
			for _, ngram := range(ngrams) {
				_ = corpus.Frequency(ngram)
			}
		}
	}
}