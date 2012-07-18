package corpustools

import (
	"fmt"
	"math"
)

type MDLSegmenter struct {
	sequence []int
	ngrams   NgramSet
}

//
// Methods to add or remove ngrams countenanced as valid segments.
//

func (mdlseg *MDLSegmenter) AddNgram(ngram []int) {
	mdlseg.ngrams.Add(ngram)
}

func (mdlseg *MDLSegmenter) RemoveNgram(ngram []int) {
	mdlseg.ngrams.Remove(ngram)
}

//
// Methods to return the description length of the corpus given the currently countenanced segments.
//

func(mdlseg *MDLSegmenter) DescriptionLength() (description_length_model, description_length_data float64) {
	// Segment the sequence.
	segmentation := mdlseg.Segment()
	// Get the unigram and bigram frequencies.
	N, U, B := mdlseg.SegmentationStats(segmentation)
	// Add the description length of the bigram statistics.
	description_length_model = mdlseg.DescriptionLengthModel(N, U, B)
	// Add the description length of the sequence given the bigram statistics.
	description_length_data = mdlseg.DescriptionLengthData(segmentation[0], N, U, B)
	return
}

// Computes the number of bits required to encode the bigram transitions in the segmented corpus.
func (mdlseg *MDLSegmenter) DescriptionLengthModel(N int, U map[string]int, B map[string]map[string]int) (description_length float64) {
	// Number of bits required to represent the biggest integer required by the stream.
	bits_per_int := math.Log2(float64(N))
	// Transmit each transition.
	for ng1, fmap := range B {
		// Number of bits to transmit the preceding symbol and the number of succeeding records to follow.
		description_length += mdlseg.HuffmanBits(ng1, N, U) + bits_per_int
		// Transmit each {succeeding symbol, frequency} pair.
		for ng2, _ := range fmap {
			description_length += mdlseg.HuffmanBits(ng2, N, U) + bits_per_int
		}
	}
	return
}

// Computes the number of bits required to encode the segmented sequence given the model (bigram transitions) has been transmitted.
func (mdlseg *MDLSegmenter) DescriptionLengthData(first_symbol string, N int, U map[string]int, B map[string]map[string]int) (description_length float64) {
	// Number of bits required to transmit first symbol.
	description_length += mdlseg.HuffmanBits(first_symbol, N, U)
	// Now transmit remaining symbols conditional on last.
	for ng1, fmap := range B {
		n := float64(U[ng1])
		for _, f_ng1_ng2 := range fmap {
			description_length += (float64(f_ng1_ng2) * (-math.Log2(float64(f_ng1_ng2) / n)))
		}
	}
	return
}

func (mdlseg *MDLSegmenter) HuffmanBits(ngram string, N int, U map[string]int) (number_bits float64) {
	return -math.Log2(float64(U[ngram]) / float64(N))
}


//
// Methods for greedily segmenting the training sequence based on the currently countenanced segments.
//

// Returns a segmented copy of the training sequence given the current ngrams. Uses a simple greedy segmentation approach.
func (mdlseg *MDLSegmenter) Segment() (segmentation []string) {
	// Iterate through the sequence and construct the segmented version.
	for pos := 0; pos < len(mdlseg.sequence); {
		match := mdlseg.Match(mdlseg.sequence[pos:])
		segmentation = append(segmentation, match)
		pos += len(match)
	}
	return
}

// Returns the longest sequence that matches the currently countenanced ngrams.
func (mdlseg *MDLSegmenter) Match(sequence []int) (match string) {
	for _, ngram_seq := range mdlseg.ngrams.ngrams {
		if len(ngram_seq) <= len(sequence) {
			// Determine whether the ngram matches the front of the sequence.
			matches := true
			for i := 0; i < len(ngram_seq); i++ {
				if ngram_seq[i] != sequence[i] {
					matches = false
					break
				}
			}
			// If the ngram matches and is longer than the current one, store it.
			if matches && len(ngram_seq) > len(match) {
				match = fmt.Sprintf("%v", ngram_seq)
			}
		}
	}
	// If no match found amongst countenanced ngrams, then return the first element of the sequence.
	if match == "" {
		match = fmt.Sprintf("%v", sequence[:1])
	}
	return
}

//
// Methods to return the unigram and bigram statistics of a segmented stream.
//

// Returns the statistics associated with a segmentation.
func (mdlseg *MDLSegmenter) SegmentationStats(segmentation []string) (N int, U map[string]int, B map[string]map[string]int) {
	U = make(map[string]int)
	B = make(map[string]map[string]int)
	// Compute unigram and bigram frequencies.
	for i := 0; i < len(segmentation); i++ {
		N += 1
		U[segmentation[i]] += 1
		if i < len(segmentation) - 1 {
			_, exists := B[segmentation[i]]
			if !exists {
				B[segmentation[i]] = make(map[string]int)
			}
			B[segmentation[i]][segmentation[i+1]] += 1
		}
	}
	return
}

// Returns an initialized MDLSegmenter based on the sequence contained in a corpus that is passed in.
func NewMDLSegmenter(corpus *Corpus) MDLSegmenter {
	return MDLSegmenter{sequence: corpus.seq, ngrams: NgramSet{ngrams: make(map[string][]int)}}
}

