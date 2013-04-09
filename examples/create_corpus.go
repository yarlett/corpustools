package main

import (
	"fmt"
	"github.com/yarlett/corpustools"
	"os"
	"strings"
)

func main() {
	// Get the path for the test corpus.
	var path, _ = os.Getwd()
	path_parts := strings.Split(path, "/")
	path_parts = path_parts[: len(path_parts) - 1]
	for _, part := range []string{"data", "brown.txt"} {
		path_parts = append(path_parts, part)
	}
	corpusfile := strings.Join(path_parts, "/")

	// Create a corpus object from the test corpus.
	lowerCase, returnChars := true, false
	corpus := corpustools.CorpusFromFile(corpusfile, lowerCase, returnChars)
	fmt.Println(corpus.Info())
}