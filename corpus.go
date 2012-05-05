package corpustools

import (
	"fmt"
	"sort"
)

// The Corpus object and its methods.
type Corpus struct {
	voc map[string]int // Mapping from input string tokens to unique integers.
	seq []int          // The raw data of the corpus stored as a sequence of integers.
	sfx [][]int        // The suffix array, containing of slices of all suffixes of the corpus.
}

func (corpus *Corpus) Info() string {
	return fmt.Sprintf("%d types and %d tokens in the corpus, %d suffixes in the suffix array.\n", len(corpus.voc), len(corpus.seq), len(corpus.sfx))
}

func (corpus *Corpus) SetSuffixArray() {
	// Create the list of suffixes from the corpus.
	suffixes := make(Seqs, 0)
	for cpos, _ := range corpus.seq {
		suffixes = append(suffixes, corpus.seq[cpos:])
	}
	sort.Sort(suffixes)
	// Assign the sorted suffixes to the corpus.
	corpus.sfx = make([][]int, 0)
	for _, seq := range suffixes {
		corpus.sfx = append(corpus.sfx, seq)
	}
}

// Slow linear search over corpus. Only to be used for testing.
func (corpus *Corpus) SearchSlow(seq []int) (slo, shi int) {
	slo = len(corpus.sfx) - 1
	shi = 0
	for spos := 0; spos < len(corpus.sfx); spos++ {
		if (len(corpus.sfx[spos]) >= len(seq)) && (SeqCmp(seq, corpus.sfx[spos][:len(seq)])) == 0 {
			if spos < slo {
				slo = spos
			}
			if spos > shi {
				shi = spos
			}
		}
	}
	return
}

// Returns the corpus indices where a given sequence occurs.
func (corpus *Corpus) Indices(seq []int) (indices []int) {
	for cpos := 0; cpos < len(corpus.seq) - (len(seq) - 1); cpos++ {
		if SeqCmpLimited(seq, corpus.seq[cpos:], len(seq)) == 0 {
			indices = append(indices, cpos)	
		}
	}
	return
}

// Binary search over suffix array to find suffix range corresponding to a specified ngram.
func (corpus *Corpus) Search(seq []int) (int, int) {
	slo, right_bound := corpus.BinarySearchLeftmost(seq, 0, len(corpus.sfx)-1)
	if slo == -1 {
		return -1, -1
	}
	shi, _ := corpus.BinarySearchRightmost(seq, slo, right_bound)
	return slo, shi
}

// Finds the first suffix pointer to the ngram using deferred detection of equality. Also returns a rightmost bound for the ngram.
func (corpus *Corpus) BinarySearchLeftmost(seq []int, smin, smax int) (int, right_bound int) {
	right_bound = smax
	for smax > smin {
		// Compare the ngram found at the search location with the desired ngram. 
		smid := (smin + smax) / 2
		cmp := SeqCmpLimited(corpus.sfx[smid], seq, len(seq))
		// Update the right bound.
		if cmp == 1 && smid < right_bound {
			right_bound = smid
		}
		// Update the search.
		if cmp == -1 {
			smin = smid + 1
		} else {
			smax = smid
		}
	}
	if (smax == smin) && SeqCmpLimited(seq, corpus.sfx[smin], len(seq)) == 0 {
		return smin, right_bound
	}
	return -1, -1
}

// Finds the last suffix pointer to the ngram using deferred detection of equality. Also returns a leftmost bound for the ngram.
func (corpus *Corpus) BinarySearchRightmost(seq []int, smin, smax int) (int, left_bound int) {
	left_bound = smin
	for smax > smin {
		// Compare the ngram found at the search location with the desired ngram. 
		smid := ((smin + smax) / 2) + 1
		cmp := SeqCmpLimited(corpus.sfx[smid], seq, len(seq))
		// Update the right bound.
		if cmp == -1 && smid > left_bound {
			left_bound = smid
		}
		// Update the search.
		if cmp == 1 {
			smax = smid - 1
		} else {
			smin = smid
		}
	}
	if (smax == smin) && SeqCmpLimited(seq, corpus.sfx[smin], len(seq)) == 0 {
		return smin, left_bound
	}
	return -1, -1
}

// Returns a copy of the corpus.
func (corpus *Corpus) Corpus() (seq []int) {
	for cpos := 0; cpos < len(corpus.seq); cpos++ {
		seq = append(seq, corpus.seq[cpos])
	}
	return
}

func (corpus *Corpus) Ngrams(order int) (ngrams [][]int) {
	for spos := 0; spos < len(corpus.sfx); spos++ {
		if (len(corpus.sfx[spos]) >= order) && ((len(ngrams) == 0) || SeqCmpLimited(ngrams[len(ngrams)-1], corpus.sfx[spos], order) != 0) {
			ngrams = append(ngrams, corpus.sfx[spos][:order])
		}
	}
	return
}

// Frequency and probability methods.

func (corpus *Corpus) Frequency(seq []int) int {
	slo, shi := corpus.Search(seq)
	if slo == -1 {
		return 0
	}
	return (shi - slo) + 1
}

func (corpus *Corpus) Probability(seq []int) float64 {
	f := corpus.Frequency(seq)
	return float64(f) / float64(len(corpus.seq)-(len(seq)-1))
}

func (corpus *Corpus) ProbabilityTransition(seq1, seq2 []int) float64 {
	return float64(corpus.Frequency(SeqJoin(seq1, seq2))) / float64(corpus.Frequency(seq1))
}

func (corpus *Corpus) ProbabilityTransitions(seq []int, predictor_length int) (probs []float64) {
	// 
	for pos := 0; pos < len(seq)-predictor_length-1; pos++ {
		// Assign conditioning and outcome elements.
		cond := seq[pos : pos+predictor_length]
		outcome := seq[pos+predictor_length : pos+predictor_length+1]
		// Assign probability of first element.
		if pos == 0 {
			probs = append(probs, corpus.Probability(cond))
		}
		// Assign transition probabilities.
		probs = append(probs, corpus.ProbabilityTransition(cond, outcome))
	}
	return
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
