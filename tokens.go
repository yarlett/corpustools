package corpustools

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

// Streams the tokens within a text file.
func TokensFromFile(filename string, lowerCase bool, returnChars bool) (tokens []string) {
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
		tks = TokenizeLine(string(line), lowerCase, returnChars)
		for _, tk := range tks {
			tokens = append(tokens, tk)
		}
		// Read from the file for the next iteration.
		line, isprefix, err = bfr.ReadLine()
	}
	return
}

// Converts a string (e.g. a line from a file) into an array of tokens.
func TokenizeLine(line string, lowerCase bool, returnChars bool) (tokens []string) {
	// Lower case everything if required.
	if lowerCase {
		line = strings.ToLower(line)
	}
	// Split line into characters.
	if returnChars {
		// Create a map of acceptable characters.
		var okChars = make(map[rune]bool)
		for _, rn := range "abcdefghijklmnopqrstuvwxyz0123456789 ,;:." {
			okChars[rn] = true
			okChars[unicode.ToUpper(rn)] = true
		}
		// Add rune to tokens if it is acceptable.
		for _, rn := range line {
			if okChars[rn] {
				tokens = append(tokens, string(rn))
			} else {
				tokens = append(tokens, "XXX")
			}
		}
		// Or else split line into "words" by splitting on space.
	} else {
		tokens = strings.Split(line, " ")
	}
	return
}
