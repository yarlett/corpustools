package corpustools

import (
	"fmt"
	"sort"
)

// The Corpus object and its methods.
type Corpus struct {
	voc map[string]int // Mapping from input string tokens to unique integers.
	seq []int          // The raw data of the corpus as a sequence of integers.
	sfx []int          // The suffix array, containing indices into the corpus.
}

func (corpus *Corpus) Info() string {
	return fmt.Sprintf("%d types and %d tokens in the corpus, %d suffixes in the suffix array.\n", len(corpus.voc), len(corpus.seq), len(corpus.sfx))
}

func (corpus *Corpus) SetSuffixArray() {
	// Create and assign the corpus suffixes to be sorted, then sort them.
	suffixes := make([]*Suffix, 0)
	for i, _ := range(corpus.seq) { suffixes = append(suffixes, &Suffix{&corpus.seq, i}) }
	sort.Sort(BySuffix{suffixes})
	// Create the array of suffix indexes from the suffix objects.
	corpus.sfx = make([]int, 0)
	for _, suffix := range(suffixes) {
		corpus.sfx = append(corpus.sfx, suffix.corpus_ix)
	}
}

// Suffix array search methods.

// Slow linear search over corpus. Only to be used for testing.
func (corpus *Corpus) SearchSlow(ngram []int) (slo, shi int) {
	slo = len(corpus.sfx) - 1
	shi = 0
	for spos := 0; spos < len(corpus.sfx); spos++ {
		cmp := corpus.NgramCmp(spos, ngram)
		if cmp == 0 {
			if spos < slo { slo = spos }
			if spos > shi { shi = spos }
		}
		if cmp == 1 {
			break
		}
	}
	return
}

// Binary search over suffix array to find suffix range corresponding to a specified ngram.
func (corpus *Corpus) Search(ngram []int) (int, int) {
	slo, right_bound := corpus.BinarySearchLeftmost(ngram, 0, len(corpus.sfx) - 1)
	if slo == -1 { return -1, -1 }
	shi, _ := corpus.BinarySearchRightmost(ngram, slo, right_bound)
	return slo, shi
}

// Finds the first suffix pointer to the ngram using deferred detection of equality. Also returns a rightmost bound for the ngram.
func (corpus *Corpus) BinarySearchLeftmost(ngram []int, smin, smax int) (int, right_bound int) {
	right_bound = smax
	for ; smax > smin; {
		// Compare the ngram found at the search location with the desired ngram. 
		smid := (smin + smax) / 2
		cmp := corpus.NgramCmp(smid, ngram)
		// Update the right bound.
		if cmp == 1 && smid < right_bound { right_bound = smid }
		// Update the search.
		if cmp == -1 {
			smin = smid + 1
		} else {
			smax = smid
		}
	}
	if ((smax == smin) && (corpus.NgramCmp(smin, ngram) == 0)) {
		return smin, right_bound
	}
	return -1, -1
}

// Finds the last suffix pointer to the ngram using deferred detection of equality. Also returns a leftmost bound for the ngram.
func (corpus *Corpus) BinarySearchRightmost(ngram []int, smin, smax int) (int, left_bound int) {
	left_bound = smin
	for ; smax > smin; {
		// Compare the ngram found at the search location with the desired ngram. 
		smid := ((smin + smax) / 2) + 1
		cmp := corpus.NgramCmp(smid, ngram)
		// Update the right bound.
		if cmp == -1 && smid > left_bound { left_bound = smid }
		// Update the search.
		if cmp == 1 {
			smax = smid - 1
		} else {
			smin = smid
		}
	}
	if ((smax == smin) && (corpus.NgramCmp(smin, ngram) == 0)) {
		return smin, left_bound
	}
	return -1, -1
}

func (corpus *Corpus) Ngrams(order int) (ngrams [][]int) {
	sfx_lst := 0
	ngrams = append(ngrams, corpus.seq[0:order])
	for sfx_cur := 0; sfx_cur < len(corpus.sfx); sfx_cur++ {
		if (corpus.NgramSfxCmp(sfx_lst, sfx_cur, order) != 0) {
			if corpus.sfx[sfx_cur] + order - 1 < len(corpus.seq) {
				ngrams = append(ngrams, corpus.seq[corpus.sfx[sfx_cur]:corpus.sfx[sfx_cur]+order])
			}
			sfx_lst = sfx_cur
		}
	}
	return
}

func (corpus *Corpus) NgramCmp(six int, ngram []int) int {
	for offset := 0; offset < len(ngram); offset++ {
		cix := corpus.sfx[six] + offset
		if cix < len(corpus.seq) {
			if corpus.seq[cix] < ngram[offset] { return -1 }
			if corpus.seq[cix] > ngram[offset] { return 1 }
		} else { return -1 }
	}
	return 0
}

func (corpus *Corpus) NgramSfxCmp(six1, six2, order int) int {
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

// Frequency and probability methods.

func (corpus *Corpus) Frequency(ngram []int) int {
	slo, shi := corpus.Search(ngram)
	if slo == -1 { return 0 }
	return (shi - slo) + 1
}

func (corpus *Corpus) Probability(ngram []int) float64 {
	f := corpus.Frequency(ngram)
	return float64(f) / float64(len(corpus.seq) - (len(ngram) - 1))
}

func (corpus *Corpus) ConditionalProbability(ngram1 []int, ngram2 []int) float64 {
	ngram_all := make([]int, len(ngram1) + len(ngram2))
	copy(ngram_all[: len(ngram1)], ngram1)
	copy(ngram_all[len(ngram1):], ngram2)
	return float64(corpus.Frequency(ngram_all)) / float64(corpus.Frequency(ngram1))
}

// Creates and returns a corpus from a text file.
func CorpusFromFile(filename string, lowerCase bool) (corpus *Corpus) {
	// Initialize the corpus.
	corpus = &Corpus{voc: make(map[string]int), seq: make([]int, 0), sfx: nil}
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