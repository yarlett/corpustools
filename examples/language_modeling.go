package main

import (
	"corpustools"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Get the path for the test corpus.
	var path, _ = os.Getwd()
	path_parts := strings.Split(path, "/")
	path_parts = path_parts[: len(path_parts) - 1]
	for _, part := range []string{"data", "test_corpus.txt"} {
		path_parts = append(path_parts, part)
	}
	corpusfile := strings.Join(path_parts, "/")

	// Create a corpus object from the test corpus.
	lower_case_tokens := true
	corpus := corpustools.CorpusFromFile(corpusfile, lower_case_tokens)
	fmt.Println(corpus.Info())

	// Calculate the mean cross-entropy of the corpus trained on itself
	corpus_sequence := corpus.Corpus()
	for predictor_length := 0; predictor_length <= 5; predictor_length++ {
		probs := corpus.ProbabilityTransitions(corpus_sequence, predictor_length)
		_, L_mn := corpustools.SummarizeProbabilities(probs)
		fmt.Printf("The mean cross-entropy of the corpus with itself using length %d predictors is %.2f bits.\n", predictor_length, L_mn)
	}
}