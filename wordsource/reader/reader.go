package reader

import (
	"bufio"
	"io"

	"github.com/charmbracelet/log"

	"ceffo.com/bee/wordsource"
)

// ReaderSource is a source that reads words from an io.Reader
type ReaderSource struct {
	reader io.Reader
}

// NewReaderSource creates a new ReaderSource
func NewReaderSource(reader io.Reader) *ReaderSource {
	if reader == nil {
		panic("reader cannot be nil")
	}
	return &ReaderSource{reader: reader}
}

// GetWords reads words from the reader
func (rs ReaderSource) GetWords() wordsource.Stream {
	log.Info("Reading words from reader")
	result := make(chan string)
	go func() {
		defer close(result)
		scanner := bufio.NewScanner(rs.reader)
		for scanner.Scan() {
			result <- scanner.Text()
		}
	}()
	return result
}
