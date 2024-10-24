package reader

import (
	"bufio"
	"io"

	"ceffo.com/bee/wordsource"
	"github.com/charmbracelet/log"
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
		log.Info("Start scanning words")
		defer close(result)
		scanner := bufio.NewScanner(rs.reader)
		for scanner.Scan() {
			result <- scanner.Text()
		}
		log.Info("Finished scanning words")
	}()
	return result
}
