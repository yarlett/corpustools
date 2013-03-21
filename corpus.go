package corpustools

import (
	"fmt"
	"runtime"
	"sort"
)

// The Corpus object and its methods.
type Corpus struct {
	voc map[string]int // Mapping from input string tokens to unique integers.
	seq []int          // The raw data of the corpus stored as a sequence of integers.
	sfx []int          // The suffix array, containing of slices of all suffixes of the corpus.
}

func (corpus *Corpus) Info() string {
	return fmt.Sprintf("%d types and %d tokens in the corpus; %d suffixes in the suffix array.", len(corpus.voc), len(corpus.seq), len(corpus.sfx))
}

//
// Sort interface methods.
//

func (corpus *Corpus) Len() int {
	return len(corpus.seq)
}

func (corpus *Corpus) Swap(i, j int) {
	corpus.sfx[i], corpus.sfx[j] = corpus.sfx[j], corpus.sfx[i]
}

func (corpus *Corpus) Less(i, j int) bool {
	if SeqCmp(corpus.seq[corpus.sfx[i]:], corpus.seq[corpus.sfx[j]:]) == -1 {
		return true
	}
	return false
}

//
// Suffix array.
//

func (corpus *Corpus) SetSuffixArray() {
	corpus.sfx = make([]int, len(corpus.seq))
	for i := 0; i < len(corpus.seq); i++ {
		corpus.sfx[i] = i
	}
	sort.Sort(corpus)
}

//
// Search methods.
//

// Returns the corpus indices where a given sequence occurs.
func (corpus *Corpus) Find(seq []int) (indices []int) {
	slo, shi := corpus.SuffixSearch(seq)
	indices = make([]int, shi-slo+1)
	i := 0
	for spos := slo; spos <= shi; spos++ {
		indices[i] = corpus.sfx[spos]
		i++
	}
	return
}

// Binary search over suffix array to find suffix range where a sequence is located.
func (corpus *Corpus) SuffixSearch(seq []int) (int, int) {
	slo, right_bound := corpus.binarySearchMin(seq, 0, len(corpus.sfx)-1)
	if slo == -1 {
		return -1, -1
	}
	shi, _ := corpus.binarySearchMax(seq, slo, right_bound)
	return slo, shi
}

// Slow linear search over corpus. Only useful for testing so not exported.
func (corpus *Corpus) slowSearch(seq []int) (slo, shi int) {
	slo = len(corpus.sfx) - 1
	shi = 0
	for spos := 0; spos < len(corpus.sfx); spos++ {
		if SeqCmpLimited(corpus.seq[corpus.sfx[spos]:], seq, len(seq)) == 0 {
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

// Returns the lowest suffix pointer to a sequence using binary search and deferred detection of equality for speed.
// Also returns a rightmost bound for the sequence which can be used to constrain the maximum search.
func (corpus *Corpus) binarySearchMin(seq []int, smin, smax int) (int, right_bound int) {
	right_bound = smax
	for smax > smin {
		// Compare the ngram found at the search location with the desired ngram. 
		smid := (smin + smax) / 2
		cmp := SeqCmpLimited(corpus.seq[corpus.sfx[smid]:], seq, len(seq))
		// Update the right bound.
		if cmp == 1 && smid < right_bound {
			right_bound = smid
		}
		// Update the search range.
		if cmp == -1 {
			smin = smid + 1
		} else {
			smax = smid
		}
	}
	if (smax == smin) && SeqCmpLimited(seq, corpus.seq[corpus.sfx[smin]:], len(seq)) == 0 {
		return smin, right_bound
	}
	return -1, -1
}

// Returns the highest suffix pointer to a sequence using binary search and deferred detection of equality for speed.
// Also returns a leftmost bound for the sequence which can be used to constrain the minimum search.
func (corpus *Corpus) binarySearchMax(seq []int, smin, smax int) (int, left_bound int) {
	left_bound = smin
	for smax > smin {
		// Compare the ngram found at the search location with the desired ngram. 
		smid := ((smin + smax) / 2) + 1
		cmp := SeqCmpLimited(corpus.seq[corpus.sfx[smid]:], seq, len(seq))
		// Update the right bound.
		if cmp == -1 && smid > left_bound {
			left_bound = smid
		}
		// Update the search range.
		if cmp == 1 {
			smax = smid - 1
		} else {
			smin = smid
		}
	}
	if (smax == smin) && SeqCmpLimited(seq, corpus.seq[corpus.sfx[smin]:], len(seq)) == 0 {
		return smin, left_bound
	}
	return -1, -1
}

//
// Utility methods.
//

// Returns a copy of the corpus.
func (corpus *Corpus) Corpus() (seq []int) {
	for cpos := 0; cpos < len(corpus.seq); cpos++ {
		seq = append(seq, corpus.seq[cpos])
	}
	return
}

// Converts a corpus sequence back into its input form.
func (corpus *Corpus) ToString(seq []int) (strings []string) {
	for pos := 0; pos < len(seq); pos++ {
		str := "**UNKNOWN**"
		for type_str, type_int := range corpus.voc {
			if type_int == seq[pos] {
				str = type_str
				break
			}
		}
		strings = append(strings, str)
	}
	return
}

//
// Ngram methods.
//

func (corpus *Corpus) Ngrams(order int) (ngrams [][]int) {
	for spos := 0; spos < len(corpus.sfx); spos++ {
		corpus_slice := corpus.seq[corpus.sfx[spos]:]
		if (len(corpus_slice) >= order) && (len(ngrams) == 0 || SeqCmpLimited(ngrams[len(ngrams)-1], corpus_slice, order) != 0) {
			ngrams = append(ngrams, corpus_slice[:order])
		}
	}
	return
}

//
// Frequency and probability methods.
//

// Returns the number of times a sequence occurs in the corpus.
func (corpus *Corpus) Frequency(seq []int) int {
	slo, shi := corpus.SuffixSearch(seq)
	if slo == -1 {
		return 0
	}
	return (shi - slo) + 1
}

// Returns the probability of a sequence in the corpus.
func (corpus *Corpus) Probability(seq []int) float64 {
	f := corpus.Frequency(seq)
	return float64(f) / float64(len(corpus.seq)-(len(seq)-1))
}

// Returns the P(seq2 | seq1) in the corpus.
func (corpus *Corpus) ProbabilityTransition(seq1, seq2 []int) float64 {
	return float64(corpus.Frequency(SeqJoin(seq1, seq2))) / float64(corpus.Frequency(seq1))
}

// Returns the probability of walking through a sequence using the corpus as training data. Useful for bigram language modeling.
func (corpus *Corpus) ProbabilityTransitions(seq []int, predictor_length int) (probs []float64) {
	// Iterate through the sequence.
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

//
// Nearest neighbor methods.
//

func (corpus *Corpus) NearestNeighbors(seq []int, seqs [][]int) (results Results) {
	// Set the maximum number of threads to be used to the number of CPU cores available.
	numprocs := runtime.NumCPU()
	runtime.GOMAXPROCS(numprocs)
	// Precompute the base vector and magnitude.
	base_vector := corpus.CoocVector(seq)
	base_mag := base_vector.Mag()
	// Initialize the channels goroutines will use to send results back.
	channel := make(chan Result, len(seqs))
	// Start the goroutines.
	for i := 0; i < len(seqs); i++ {
		go corpus.NearestNeighborWorker(base_vector, base_mag, seqs[i], channel)
	}
	// Drain the channels of results.
	for i := 0; i < len(seqs); i++ {
		result, _ := <-channel
		results = append(results, result)
	}
	// Return the sorted results.
	sort.Sort(results)
	return
}

func (corpus *Corpus) NearestNeighborWorker(base_vector *Cooc, base_mag float64, seq []int, results_channel chan Result) {
	// Compute the similarity between the base vector and a specified sequence.
	cooc := corpus.CoocVector(seq)
	results_channel <- Result{Seq: seq, Val: base_vector.Prod(cooc) / (base_mag * cooc.Mag())}
}

// Returns a co-occurrence vector for a sequence.
func (corpus *Corpus) CoocVector(seq []int) (cooc *Cooc) {
	// Get suffix range where the sequence occurs.
	slo, shi := corpus.SuffixSearch(seq)
	// Get the frequency counts.
	cooc = &Cooc{seq: seq, dat: make(map[int]float64)}
	for spos := slo; spos <= shi; spos++ {
		cpos := corpus.sfx[spos]
		// Increment the count of the type occurring before the sequence.
		if cpos > 0 {
			cooc.Inc(corpus.seq[cpos-1])
		}
		// Increment the count of the type occurring after the sequence.
		if cpos < len(corpus.seq)-1 {
			cooc.Inc(corpus.seq[cpos+1])
		}
	}
	return
}

//
// Functions to create a corpus.
//

// Creates and returns a corpus from a text file.
func CorpusFromFile(filename string, lowerCase bool, returnChars bool) (corpus *Corpus) {
	// Initialize the corpus.
	corpus = &Corpus{voc: make(map[string]int), seq: make([]int, 0), sfx: nil}
	// Get string array from tokenizer.
	tokens := TokensFromFile(filename, lowerCase, returnChars)
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
