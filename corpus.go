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
	return fmt.Sprintf("%d types and %d tokens in the corpus.\n", len(corpus.voc), len(corpus.seq))
}

func (corpus *Corpus) SetSuffixArray() {
	// Create and assign the corpus suffixes to be sorted, then sort them.
	suffixes := make([]*Suffix, 0)
	for i, _ := range(corpus.seq) { suffixes = append(suffixes, &Suffix{&corpus.seq, i}) }
	sort.Sort(BySuffix{suffixes})
	// Create the array of suffix indexes from the suffix objects.
	corpus.sfx = make([]int, len(suffixes))
	for _, suffix := range(suffixes) {
		corpus.sfx = append(corpus.sfx, suffix.corpus_ix)
	}
}

// Search methods.

func (corpus *Corpus) Search(ngram []int) (slo, shi int) {
	// // Get first suffix index hit by binary search.
	// six := corpus.BinarySearch(ngram, 0, len(corpus.sfx) - 1)
	// if six == -1 { return -1, -1 }
	// // Perform linear search left and right to get inclusive suffix range.
	// for slo = six; (slo > 0) && corpus.NgramCmp(slo - 1, ngram) == 0; slo-- {}
	// for shi = six; (shi < len(corpus.sfx) - 1) && corpus.NgramCmp(shi + 1, ngram) == 0; shi++ {}
	// return slo, shi

	slo = corpus.BinarySearchLeft(ngram, 0, len(corpus.sfx) - 1)
	shi = corpus.BinarySearchRight(ngram, slo, len(corpus.sfx) - 1)
	return slo, shi
}

func (corpus *Corpus) BinarySearch(ngram []int, smin, smax int) int {
	var (
		smid, cmp int
	)
	for ; smax >= smin; {
		smid = (smin + smax) / 2
		cmp = corpus.NgramCmp(smid, ngram)
		if cmp == -1 {
			smin = smid + 1
		} else if cmp == 1 {
			smax = smid - 1
		} else {
			return smid
		}
	}
	return -1
}

func (corpus *Corpus) BinarySearchLeft(ngram []int, smin, smax int) int {
	for ; smax >= smin; {
		smid := (smin + smax) / 2
		cmp := corpus.NgramCmp(smid, ngram)
		if cmp == -1 {
			smin = smid + 1
		} else if cmp == 1 {
			smax = smid - 1
		} else {
			if smid > 0 {
				cmp_left := corpus.NgramCmp(smid - 1, ngram)
				if cmp_left != 0 {
					return smid
				} else {
					smax = smid - 1
				}
			} else {
				return smid
			}
		}
	}
	return -1
}

func (corpus *Corpus) BinarySearchRight(ngram []int, smin, smax int) int {
	for ; smax >= smin; {
		smid := (smin + smax) / 2
		cmp := corpus.NgramCmp(smid, ngram)
		if cmp == -1 {
			smin = smid + 1
		} else if cmp == 1 {
			smax = smid - 1
		} else {
			if smid < len(corpus.sfx) - 1 {
				cmp_right := corpus.NgramCmp(smid + 1, ngram)
				if cmp_right != 0 {
					return smid
				} else {
					smin = smid + 1
				}
			} else {
				return smid
			}
		}
	}
	return -1
}

func (corpus *Corpus) Ngrams(order int) (ngrams [][]int) {
	sfx_lst := 0
	ngrams = append(ngrams, corpus.seq[0:order])
	for sfx_cur := 0; sfx_cur < len(corpus.sfx); sfx_cur++ {
		if (corpus.NgramSfxCmp(sfx_lst, sfx_cur, order) != 0) {
			if corpus.sfx[sfx_cur] + order < len(corpus.seq) {
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
	corpus = &Corpus{make(map[string]int), make([]int, 0), nil}
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