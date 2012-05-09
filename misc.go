package corpustools

import (
	"math"
)

//
// Various utility functions related to sequences.
//

func SeqJoin(seq1, seq2 []int) (joined []int) {
	for i := 0; i < len(seq1); i++ {
		joined = append(joined, seq1[i])
	}
	for i := 0; i < len(seq2); i++ {
		joined = append(joined, seq2[i])
	}
	return
}

func SeqCmp(seq1, seq2 []int) int {
	// Get lengths of sequences, and shortest length.
	len1, len2 := len(seq1), len(seq2)
	length := len1
	if len2 < length {
		length = len2
	}
	// Make comparisons over defined extent of ngrams.
	for pos := 0; pos < length; pos++ {
		if seq1[pos] < seq2[pos] {
			return -1
		}
		if seq1[pos] > seq2[pos] {
			return 1
		}
	}
	// All assigned elements are the same, so make comparison based on length (shorter is lesser).
	if len1 == len2 {
		return 0
	} else if len1 < len2 {
		return -1
	}
	return 1
}

func SeqCmpLimited(seq1, seq2 []int, comparison_length int) int {
	// Get lengths of sequences, and shortest length.
	len1, len2 := len(seq1), len(seq2)
	length := len1
	if len2 < length {
		length = len2
	}
	if comparison_length < length {
		length = comparison_length
	}
	// Make comparisons over defined extent of ngrams.
	for pos := 0; pos < length; pos++ {
		if seq1[pos] < seq2[pos] {
			return -1
		}
		if seq1[pos] > seq2[pos] {
			return 1
		}
	}
	// All assigned elements are the same, so make comparison based on length (shorter is lesser).
	if length == comparison_length {
		return 0
	} else if len1 < len2 {
		return -1
	}
	return 1
}

//
// Probability functions.
//

func SummarizeProbabilities(probs []float64) (float64, float64) {
	L := 0.0
	for i := 0; i < len(probs); i++ {
		L -= math.Log2(probs[i])
	}
	return L, L / float64(len(probs))
}
