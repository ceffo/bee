package wordsource

// Stream is a channel of strings
type Stream chan string

// Source is a source of words
type Source interface {
	GetWords() Stream
}

// Maker is a function that returns a Source
type Maker func() Source
