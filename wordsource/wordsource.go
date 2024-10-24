package wordsource

// Stream is a channel of strings
type Stream chan string

// Source is a source of words
type Source interface {
	GetWords() Stream
}

type SourceMaker func() Source
