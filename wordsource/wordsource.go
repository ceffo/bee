package wordsource

type Stream chan string

type Source interface {
	GetWords() Stream
}
