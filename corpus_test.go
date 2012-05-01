package corpustools

import (
	"testing"
)

func TestCorpus(t *testing.T) {
	corpus := Corpus{seq: []int{5, 4, 3, 2, 1}}
	corpus.SetSuffixArray()
	for i, ix := range corpus.sfx {
		t.Error(i, ix)
	}
}
