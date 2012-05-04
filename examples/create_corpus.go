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
}