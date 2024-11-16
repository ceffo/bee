package bee

import (
	"strings"
	"testing"

	fakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	"ceffo.com/bee/pkg/must"
	"ceffo.com/bee/wordsource"
)

type TestWordSource struct {
	words []string
}

func NewFakeTestWordSource(numWords, seed int) *TestWordSource {
	err := fakeit.Seed(seed)
	if err != nil {
		panic(err)
	}
	words := make([]string, numWords)
	for i := range numWords {
		fakeWord := strings.ToUpper(fakeit.Word())
		words[i] = fakeWord
	}
	return &TestWordSource{words: words}
}

func NewFixedTestWordSource(words []string) *TestWordSource {
	return &TestWordSource{words: words}
}

func (tws TestWordSource) GetWords() wordsource.Stream {
	result := make(chan string)
	go func() {
		defer close(result)
		for _, word := range tws.words {
			result <- word
		}
	}()
	return result
}

func TestBeesolve_SolveFor(t *testing.T) {
	wordList := []string{
		"MANUAL",
		"MATURE",
		"MANUALLY",
		"MATURELY",
		"NULL",
		"AMATEUR",
		"RUNNY",
	}
	wordSource := NewFixedTestWordSource(wordList)
	tr := NewSolver(func() wordsource.Source {
		return wordSource
	})

	tests := []struct {
		input *Input
		want  []string
	}{
		{
			input: must.NoError(NewInput('N', 'M', 'A', 'U', 'L', 'Y', 'T')),
			want: []string{
				"MANUAL",
				"MANUALLY",
				"NULL",
			},
		},
		{
			input: must.NoError(NewInput('M', 'N', 'A', 'U', 'L', 'Y', 'T')),
			want: []string{
				"MANUAL",
				"MANUALLY",
			},
		},
	}
	for _, tt := range tests {
		name := tt.input.String()
		t.Run(name, func(t *testing.T) {
			result := tr.SolveFor(tt.input)
			words := make([]string, 0, len(tt.want))
			for word := range result {
				words = append(words, word)
			}
			assert.ElementsMatch(t, tt.want, words)
		})
	}
}
