package corpustools

import (
	"bufio"
	"fmt"
	"log"
	"io"
	"os"
	"strings"
)

// Streams the tokens within a text file.
func TokenStreamer(filename string) {
	var (
		bfr *bufio.Reader
		tks []string
		numtoks int64
	)
	// Open the file for reading.
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	// Read the lines of the file one at a time.
	bfr = bufio.NewReaderSize(fh, 1024 * 16)
	for line, isprefix, err := bfr.ReadLine(); err != io.EOF; {
		// Error handling.
		if err != nil {
			log.Fatal(err)
		}
		if isprefix {
			log.Fatal("line too long for buffered reader")
		}
		// Convert the bytes in the line to nice tokens.
		tks = TokenizeLine(string(line))
		for _, tk := range(tks) {
			fmt.Println(tk)
			numtoks++
		}
		// Read from the file for the next iteration.
		line, isprefix, err = bfr.ReadLine()
	}
	fmt.Println(numtoks)
}

// Converts a string (e.g. a line from a file) into an array of tokens.
func TokenizeLine(line string) (tks []string) {
	return strings.Split(line, " ")
}