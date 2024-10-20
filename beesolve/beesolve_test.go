package beesolve

import (
	"strings"
	"testing"

	"ceffo.com/bee/wordsource"
	fakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

type TestWordSource struct {
	words []string
}

func NewFakeTestWordSource(numWords, seed int) *TestWordSource {
	fakeit.Seed(seed)
	words := make([]string, numWords)
	for i := 0; i < numWords; i++ {
		fakeWord := strings.ToLower(fakeit.Word())
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
	words := []string{
		"manual",
		"mature",
		"manually",
		"maturely",
		"null",
		"amateur",
		"runny",
	}
	wordSource := NewFixedTestWordSource(words)
	tr := NewBeeSolve(wordSource)

	tests := []struct {
		input Input
		want  []string
	}{
		{
			input: NewBeeInput('n', []rune{'m', 'a', 'u', 'l'}),
			want:  []string{"MANUAL", "NULL"},
		},
		{
			input: NewBeeInput('n', []rune{'m', 'a', 'r', 'u', 'l', 'y', 't'}),
			want: []string{
				"MANUAL",
				"MANUALLY",
				"NULL",
				"RUNNY",
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
