package corpustools

import (
	"fmt"
	"math"
)

func (corpus *Corpus) DescriptionLengthDelta(seq[] int) float64 {
	// Get cost of internal transitions through sequence.
	internal_cost := 0.0
	for i := 0; i < len(seq) - 1; i++ {
		internal_cost += corpus.TransitionCost(seq[i : i+1], seq[i+1: i+2])
	}
	// Get the corpus indices at which the sequence occurs.
	indices := corpus.Indices(seq)
	// Iterate through the indices and measure the transition length without and with the sequence in the lexicon.
	l0, l1 := 0.0, 0.0
	T := make(map[string]bool, 0)
	for _, cpos := range indices {
		if cpos > 0 && cpos < len(corpus.seq) - len(seq) {
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
	transition_table_cost := float64(len(T) * (len(seq) + 2)) * 16.0
	//fmt.Printf("Coding advantage = %.2f, table cost = %.2f\n", l0 - l1, transition_table_cost)
	return (l0 - l1) - transition_table_cost
}



func (corpus *Corpus) TransitionCost(seq1, seq2 []int) float64 {
	return -math.Log2(corpus.ProbabilityTransition(seq1, seq2))
}


// func (corpus *Corpus) DescriptionLengthDelta(seq []int) float64 {
// 	// Maps to keep track of the number of distinct preceding and succeeding elements.
// 	P := make(map[string]int, 0)
// 	S := make(map[string]int, 0)
// 	// Description length of transitions internal to the ngram.
// 	tprobs := corpus.ProbabilityTransitions(seq, 1)[1: ]
// 	length_internal, _ = SummarizeProbabilities(tprobs)
// 	// Initialize the length of encoding the data.
// 	length_base := 0.0
// 	length_with := 0.0
// 	// Iterate through every occurrence of the sequence in question.
// 	slo, shi := corpus.Search(seq)
// 	for six := slo; six <= shi; six++ {


// 		cix := corpus.sfx[six]
// 		if (cix > 0) && (cix <= len(corpus.seq)-len(ngram)) {
// 			// Identify the preceding and succeeding elements.
// 			p := corpus.seq[cix-1 : cix]
// 			s := corpus.seq[cix+1 : cix+2]
// 			// Update the encountered elements.
// 			P[fmt.Sprintf("%v", p)] += 1
// 			S[fmt.Sprintf("%v", s)] += 1
// 			// Add the length of transitions without the ngram in the lexicon.
// 			length_base -= math.Log2(corpus.ConditionalProbability(p, ngram[:1]))
// 			length_base += length_internal
// 			length_base -= math.Log2(corpus.ConditionalProbability(ngram[len(ngram)-1:], s))
// 			// Add the length of transitions with the ngram in the lexicon.
// 			length_with -= math.Log2(corpus.ConditionalProbability(p, ngram))
// 			length_with -= math.Log2(corpus.ConditionalProbability(ngram, s))
// 		}
// 	}
// 	// Calculate the encoding length of the transitions table.
// 	encoding_length := float64(len(P)) * (1.0 + float64(len(ngram)) + 1.0) * 32.0
// 	encoding_length += float64(len(S)) * (float64(len(ngram)) + 1.0 + 1.0) * 32.0
// 	// Return the overall impact on description length.
// 	return encoding_length - length_base + length_with
// }
