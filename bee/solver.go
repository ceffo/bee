package bee

import (
	"strings"

	"github.com/charmbracelet/log"
	mapset "github.com/deckarep/golang-set/v2"

	"ceffo.com/bee/wordsource"
)

const (
	solverChannelSize = 100
)

// Solver solves the Bee game
type Solver struct {
	source wordsource.Maker
}

// NewSolver creates a new Solver
func NewSolver(maker wordsource.Maker) *Solver {
	return &Solver{source: maker}
}

// SolveFor solves the Bee game for the given input
func (t *Solver) SolveFor(input *Input) wordsource.Stream {
	result := make(chan string, solverChannelSize)
	go func() {
		log.Infof("Solving for '%s'", input)
		defer close(result)
		numfound := 0
		for word := range t.source().GetWords() {
			word = strings.ToUpper(word)
			if satisfies(word, input) {
				result <- word
				numfound++
			}
		}
		log.Infof("Found %d words", numfound)
	}()
	return result
}

// satisfies returns true if the word satisfies the input
func satisfies(word string, input *Input) bool {
	if len(word) < minWordLength {
		return false
	}
	wordRunes := []rune(word)
	wordSet := mapset.NewSet(wordRunes...)
	return wordSet.Contains(input.center) && wordSet.IsSubset(input.lettersSet)
}
