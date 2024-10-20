package beesolve

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

// Input represents the input to the Bee solver
type Input struct {
	center  rune
	letters []rune
}

func (i Input) Score(word string) int {
	l := len(word)
	if l < 4 {
		return 0
	}
	// 4 letters words are worth 1 point
	// longer words earn 1 point for each letter
	// using all 7 letters earns a 7 point bonus
	score := l - 3
	if l >= len(i.letters) && i.IsPangram(word) {
		score += 7
	}
	return score
}

func (i Input) IsPangram(word string) bool {
	letterSet := mapset.NewSet(i.letters...)
	wordRunes := []rune(strings.ToLower(word))
	wordSet := mapset.NewSet(wordRunes...)
	return letterSet.Equal(wordSet)
}

func NewBeeInput(center rune, letters []rune) Input {
	return Input{center: center, letters: append(letters, center)}
}

func (i Input) String() string {
	return string(i.letters)
}
