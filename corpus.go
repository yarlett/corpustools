package corpustools

import (
	"fmt"
	"sort"
)

// The Corpus object and its methods.
type Corpus struct {
	voc map[string]int // Input strings get mapped to unique integers.
	seq []int          // The sequence of integers which represents the corpus.
	sfx []int          // The suffix array, containing indices into the corpus.
}

func (corpus *Corpus) Info() string {
	return fmt.Sprintf("%d types and %d tokens in the corpus.\n", len(corpus.voc), len(corpus.seq))
}

func (corpus *Corpus) SetSuffixArray() {
	// Create and assign the corpus suffixes to be sorted, then sort them.
	suffixes := make([]*Suffix, 0)
	for i := 0; i < len(corpus.seq); i++ {
		suffixes = append(suffixes, &Suffix{&corpus.seq, i})
	}
	sort.Sort(BySuffix{suffixes})
	// Create the array of suffix indexes from the suffix objects.
	for _, suffix := range suffixes {
		corpus.sfx = append(corpus.sfx, suffix.corpus_ix)
	}
}

func (corpus *Corpus) SearchNgram(ngram []int) {
	// Write me using binary search over the suffix array, to return lowest and highest suffix indices of the ngram in question.
}

func (corpus *Corpus) Ngrams(order int) (ngrams [][]int) {
	sfx_lst := 0
	ngrams = append(ngrams, corpus.seq[0:order])
	for sfx_cur := 0; sfx_cur < len(corpus.sfx); sfx_cur++ {
		if corpus.NgramCmp(sfx_lst, sfx_cur, order) != 0 {
			ngrams = append(ngrams, corpus.seq[corpus.sfx[sfx_cur]:corpus.sfx[sfx_cur]+order])
			sfx_lst = sfx_cur
		}
	}
	return
}

func (corpus *Corpus) NgramCmp(six1, six2, order int) int {
	if six1 == six2 { return 0 }
	for offset := 0; offset < order; offset++ {
		cix1 := corpus.sfx[six1] + offset
		cix2 := corpus.sfx[six2] + offset
		if cix1 < len(corpus.seq) && cix2 < len(corpus.seq) {
			if corpus.seq[cix1] < corpus.seq[cix2] { return -1 }
			if corpus.seq[cix2] < corpus.seq[cix1] { return 1 }
		} else {
			if cix1 >= len(corpus.seq) { return -1 }
			if cix2 >= len(corpus.seq) { return 1 }
		}
	}
	return 0
}

// Creates and returns a corpus from a text file.
func CorpusFromFile(filename string, lowerCase bool) (corpus *Corpus) {
	// Initialize the corpus.
	corpus = &Corpus{make(map[string]int), make([]int, 0), make([]int, 0)}
	// Get string array from tokenizer.
	tokens := TokensFromFile(filename, lowerCase)
	// Iterate through the string tokens.
	type_ctr := 0
	for _, token := range tokens {
		// Get the unique identifier for the token.
		_, found := corpus.voc[token]
		if !found {
			corpus.voc[token] = type_ctr
			type_ctr++
		}
		// Populate the corpus.
		corpus.seq = append(corpus.seq, corpus.voc[token])
	}
	// Compute the suffix array.
	corpus.SetSuffixArray()
	return
}