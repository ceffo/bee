package bee

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	mapset "github.com/deckarep/golang-set/v2"

	"ceffo.com/bee/pkg/slices"
)

var (
	ErrInvalidInputSize = errors.New("invalid input size")
	ErrCenterNotInInput = errors.New("center letter not in input")
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
	center  rune
	letters []rune
}

func (i *Input) Validate() error {
	if len(i.letters) != NumLetters {
		return ErrInvalidInputSize
	}
	set := mapset.NewSet(i.letters...)
	if !set.Contains(i.center) {
		return ErrCenterNotInInput
	}
	if set.Cardinality() != NumLetters {
		return ErrDuplicateLetters
	}
	return nil
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
	if l >= len(i.letters) && i.IsPangram(word) {
		score += pangramBonus
	}
	return score
}

func (i *Input) Center() rune {
	return i.center
}

func (i *Input) IsPangram(word string) bool {
	letterSet := mapset.NewSet(i.letters...)
	wordRunes := []rune(strings.ToUpper(word))
	wordSet := mapset.NewSet(wordRunes...)
	return letterSet.Equal(wordSet)
}

func (i *Input) IsExactPangram(word string) bool {
	return i.IsPangram(word) && len(word) == NumLetters
}

func NewInput(center rune, letters []rune) (*Input, error) {
	center = unicode.ToUpper(center)
	letters = slices.Map(letters, unicode.ToUpper)
	i := &Input{center: center, letters: append(letters, center)}
	if err := i.Validate(); err != nil {
		return nil, err
	}
	return i, nil
}

func NewFrom(str string) (*Input, error) {
	if len(str) != NumLetters {
		return nil, ErrInvalidInputSize
	}
	center, _ := utf8.DecodeRuneInString(str)
	letters := []rune(str)[1:]
	return NewInput(center, letters)
}

func (i *Input) String() string {
	return string(i.letters)
}
