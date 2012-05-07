package corpustools

import (
	"fmt"
	"math"
)

func (corpus *Corpus) DescriptionLengthDelta(seq []int) float64 {
	// Get cost of internal transitions through sequence.
	internal_cost := 0.0
	for i := 0; i < len(seq)-1; i++ {
		internal_cost += corpus.TransitionCost(seq[i:i+1], seq[i+1:i+2])
	}
	// Get the corpus indices at which the sequence occurs.
	indices := corpus.Find(seq)
	// Iterate through the indices and measure the transition length without and with the sequence in the lexicon.
	l0, l1 := 0.0, 0.0
	T := make(map[string]bool, 0)
	for _, cpos := range indices {
		if cpos > 0 && cpos < len(corpus.seq)-len(seq) {
			// Identify the items occurring before and after the sequence of interest.
			before := corpus.seq[cpos-1 : cpos]
			after := corpus.seq[cpos+len(seq) : cpos+len(seq)+1]
			// Add the transition cost without the sequence in the lexicon.
			l0 += corpus.TransitionCost(before, seq[:1])
			l0 += internal_cost
			l0 += corpus.TransitionCost(seq[len(seq)-1:], after)
			// Add the transition cost with the sequence in the lexicon.
			l1 += corpus.TransitionCost(before, seq)
			l1 += corpus.TransitionCost(seq, after)
			// Update the count of the total number of transitions the posited sequence is involved in.
			transition_key := fmt.Sprintf("%v --> %v", before, seq)
			T[transition_key] = true
			transition_key = fmt.Sprintf("%v --> %v", seq, after)
			T[transition_key] = true
		}
	}
	//transition_table_cost := float64(len(T)*(len(seq)+2)) * 16.0
	return (l0 - l1) // - transition_table_cost
}

func (corpus *Corpus) TransitionCost(seq1, seq2 []int) float64 {
	return -math.Log2(corpus.ProbabilityTransition(seq1, seq2))
}
