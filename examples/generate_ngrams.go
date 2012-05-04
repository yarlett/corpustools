package main

import (
	"corpustools"
	"fmt"
	"os"
	"strings"
	"time"
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

	// Iterate over various orders and generate the ngrams of this length.
	for length := 1; length <= 10; length++ {
		t1 := time.Now()
		ngrams := corpus.Ngrams(length)
		t2 := time.Now()
		fmt.Printf("%d %dgrams found in %v.\n", len(ngrams), length, t2.Sub(t1))
	}
}