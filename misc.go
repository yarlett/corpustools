package corpustools

import (
	"math"
)

func SummarizeProbabilities(probs []float64) (float64, float64) {
	L := 0.0
	for i := 0; i < len(probs); i++ {
		L -= math.Log2(probs[i])
	}
	return L, L / float64(len(probs))
}
