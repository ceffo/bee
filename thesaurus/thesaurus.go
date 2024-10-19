package thesaurus

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type WordStream chan string

type WordSource interface {
	GetWordSource() WordStream
}

type Thesaurus struct {
	wordSource WordSource
}

func NewThesaurus(wordSource WordSource) *Thesaurus {
	return &Thesaurus{wordSource: wordSource}
}

func (t *Thesaurus) FindWordsWithJust(runes ...rune) WordStream {
	result := make(chan string)
	letterSet := mapset.NewSet(runes...)

	go func() {
		defer close(result)
		for word := range t.wordSource.GetWordSource() {
			if wordsHasOnlyLetters(word, letterSet) {
				result <- word
			}
		}
	}()
	return result
}

func wordsHasOnlyLetters(word string, letters mapset.Set[rune]) bool {
	wordRunes := []rune(word)
	wordSet := mapset.NewSet(wordRunes...)
	return wordSet.IsSubset(letters)
}
