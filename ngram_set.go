package corpustools

import (
	"fmt"
)

type NgramSet struct {
	ngrams map[string][]int
}

func (ngs *NgramSet) Key(ngram []int) (key string) {
	key = fmt.Sprintf("%v", ngram)
	return
}

func (ngs *NgramSet) In(ngram []int) (in bool) {
	_, in = ngs.ngrams[ngs.Key(ngram)]
	return
}

func (ngs *NgramSet) Add(ngram []int) {
	ngs.ngrams[ngs.Key(ngram)] = ngram
	return
}

func (ngs *NgramSet) Remove(ngram []int) {
	delete(ngs.ngrams, ngs.Key(ngram))
	return
}

func (ngs *NgramSet) Size() (size int) {
	size = len(ngs.ngrams)
	return
}

func (ngs *NgramSet) LongestNgram() (longest int) {
	for _, ngram_seq := range ngs.ngrams {
		if len(ngram_seq) > longest {
			longest = len(ngram_seq)
		}
	}
	return
}
