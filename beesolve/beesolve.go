package beesolve

import (
	"strings"

	"ceffo.com/bee/wordsource"
	mapset "github.com/deckarep/golang-set/v2"
)

type BeeSolve struct {
	wordSource wordsource.Source
}

func NewBeeSolve(wordSource wordsource.Source) *BeeSolve {
	return &BeeSolve{wordSource: wordSource}
}

func (t *BeeSolve) SolveFor(input Input) wordsource.Stream {
	result := make(chan string)

	go func() {
		defer close(result)
		for word := range t.wordSource.GetWords() {
			if meetsBeeRequirements(word, input) {
				result <- strings.ToUpper(word)
			}
		}
	}()
	return result
}

func meetsBeeRequirements(word string, input Input) bool {
	if len(word) < 4 {
		return false
	}
	wordRunes := []rune(strings.ToLower(word))
	wordSet := mapset.NewSet(wordRunes...)
	letters := mapset.NewSet(input.letters...)
	return wordSet.Contains(input.center) && wordSet.IsSubset(letters)
}
