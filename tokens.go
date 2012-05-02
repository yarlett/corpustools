package corpustools

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// Streams the tokens within a text file.
func TokensFromFile(filename string, lowerCase bool) (tokens []string) {
	var (
		bfr *bufio.Reader
		tks []string
	)
	tokens = make([]string, 0)
	// Open the file for reading.
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	// Read the lines of the file one at a time.
	bfr = bufio.NewReaderSize(fh, 1024*16)
	for line, isprefix, err := bfr.ReadLine(); err != io.EOF; {
		// Error handling.
		if err != nil {
			log.Fatal(err)
		}
		if isprefix {
			log.Fatal("line too long for buffered reader")
		}
		// Convert the bytes in the line to nice tokens.
		tks = TokenizeLine(string(line), lowerCase)
		for _, tk := range tks {
			tokens = append(tokens, tk)
		}
		// Read from the file for the next iteration.
		line, isprefix, err = bfr.ReadLine()
	}
	return
}

// Converts a string (e.g. a line from a file) into an array of tokens.
func TokenizeLine(line string, lowerCase bool) (tks []string) {
	tks = strings.Split(line, " ")
	if lowerCase {
		for i, _ := range tks {
			tks[i] = strings.ToLower(tks[i])
		}
	}
	return
}