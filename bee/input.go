package bee

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

const (
	// minWordLength is the minimum length of a word
	minWordLength = 4
	// pangramBonus is the bonus for using all the letters
	pangramBonus = 7
	// NumLetters is the number of letters in a game
	NumLetters = 7
)

// Input represents the input to the Bee solver
type Input struct {
	center  rune
	letters []rune
}

func (i Input) Score(word string) int {
	l := len(word)
	if l < minWordLength {
		return 0
	}
	// 4 letters words are worth 1 point
	// longer words earn 1 point for each letter
	// using all letters earns a 7 point bonus
	score := 1
	if l > minWordLength {
		score = l
	}
	if l >= len(i.letters) && i.IsPangram(word) {
		score += pangramBonus
	}
	return score
}

func (i Input) Center() rune {
	return i.center
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
