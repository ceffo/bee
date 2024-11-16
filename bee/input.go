package bee

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	mapset "github.com/deckarep/golang-set/v2"

	"ceffo.com/bee/pkg/slices"
)

var (
	ErrInvalidInputSize = errors.New("invalid input size")
	ErrDuplicateLetters = errors.New("duplicate letters in input")
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
	center     rune
	lettersSet mapset.Set[rune]
}

func (i *Input) Score(word string) int {
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
	if l >= i.lettersSet.Cardinality() && i.IsPangram(word) {
		score += pangramBonus
	}
	return score
}

func (i *Input) Center() rune {
	return i.center
}

func (i *Input) IsPangram(word string) bool {
	wordRunes := []rune(strings.ToUpper(word))
	wordSet := mapset.NewSet(wordRunes...)
	return i.lettersSet.Equal(wordSet)
}

func (i *Input) IsExactPangram(word string) bool {
	return i.IsPangram(word) && len(word) == NumLetters
}

// NewInput creates a new input from a set of letters, considering the first letter as the center
func NewInput(letters ...rune) (*Input, error) {
	if len(letters) != NumLetters {
		return nil, fmt.Errorf("%w: expected %d letters, got %d", ErrInvalidInputSize, NumLetters, len(letters))
	}
	letters = slices.Map(letters, unicode.ToUpper)
	lettersSet := mapset.NewSet(letters...)
	if lettersSet.Cardinality() != NumLetters {
		return nil, ErrDuplicateLetters
	}
	return &Input{center: letters[0], lettersSet: lettersSet}, nil
}

func NewInputFrom(str string) (*Input, error) {
	return NewInput([]rune(str)...)
}

func (i *Input) String() string {
	withoutCenter := i.lettersSet.Clone()
	withoutCenter.Remove(i.center)
	return fmt.Sprintf("%c|%s", i.center, string(withoutCenter.ToSlice()))
}
