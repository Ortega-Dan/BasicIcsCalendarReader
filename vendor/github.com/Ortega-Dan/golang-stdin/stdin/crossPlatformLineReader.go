package stdin

import (
	"bufio"
	"os"
	"strings"
)

var reader *bufio.Reader

// Reads a line from std input without default trim.
func ReadLineNoTrim() string {

	// Initializing reader on first usage
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	// Reading line and error
	text, err := reader.ReadString('\n')

	// in case of error, just panic
	if err != nil {
		panic(err)
	}

	// convert CRLF to LF
	text = strings.Replace(text, "\r", "", -1)

	// Remove the last character LF
	text = text[0 : len(text)-1]

	return text
}

// Reads a line from std input with default trim.
func ReadLine() string {

	text := ReadLineNoTrim()

	// Trim space around
	text = strings.TrimSpace(text)

	return text
}

// If no error opening the file exists, returns a
// function that reads line by line when called.
// The returned function returns io.EOF as error when successfully done with the reading,
// or another error if a different error existed.
func NewLineByLineFileReader(path string) (readerFunc func() (string, error), err error) {

	// opening file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// setting reader
	reada := bufio.NewReader(file)

	// returning reader function
	return func() (string, error) {

		// reading any length lines
		str, err := reada.ReadString('\n')

		// cross-platform cleanning of CRLF or LF
		str = strings.Replace(str, "\r", "", -1)
		str = strings.Replace(str, "\n", "", -1)

		// intended final return
		if err != nil {
			return str, err
		}

		// returning the desired read text
		return str, nil
	}, nil
}
