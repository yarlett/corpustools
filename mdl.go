package corpustools

// import (
// 	"fmt"
// 	"math"
// )

// func (corpus *Corpus) DescriptionLengthDelta(ngram []int) float64 {
// 	// Maps to keep track of the number of distinct preceding and succeeding elements.
// 	P := make(map[string]int, 0)
// 	S := make(map[string]int, 0)
// 	// Description length of transitions internal to the ngram.
// 	length_internal := 0.0
// 	for i := 0; i < len(ngram)-1; i++ {
// 		length_internal -= math.Log(corpus.ConditionalProbability(ngram[i:i+1], ngram[i+1:i+2]))
// 	}
// 	// Initialize the length of encoding the data.
// 	length_base := 0.0
// 	length_with := 0.0
// 	// Iterate through every occurrence of the ngram in question.
// 	slo, shi := corpus.Search(ngram)
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
