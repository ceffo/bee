package bee

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"

	"ceffo.com/bee/wordsource"
)

type Solver struct {
	wordSource wordsource.Source
}

func NewSolver(wordSource wordsource.Source) *Solver {
	return &Solver{wordSource: wordSource}
}

func (t *Solver) SolveFor(input Input) wordsource.Stream {
	result := make(chan string)
	go func() {
		defer close(result)
		for word := range t.wordSource.GetWords() {
			word = strings.ToUpper(word)
			if satisfies(word, input) {
				result <- word
			}
		}
	}()
	return result
}

// satisfies returns true if the word satisfies the input
func satisfies(word string, input Input) bool {
	if len(word) < minWordLength {
		return false
	}

	wordRunes := []rune(word)
	wordSet := mapset.NewSet(wordRunes...)
	letters := mapset.NewSet(input.letters...)
	return wordSet.Contains(input.center) && wordSet.IsSubset(letters)
}
