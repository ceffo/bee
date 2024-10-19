package thesaurus

import (
	"strings"
	"testing"

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

func (tws TestWordSource) GetWordSource() WordStream {
	result := make(chan string)
	go func() {
		defer close(result)
		for _, word := range tws.words {
			result <- word
		}
	}()
	return result
}

func TestThesaurus_FindWordsWithJust(t *testing.T) {
	words := []string{
		"manual",
		"mature",
		"manually",
		"maturely",
		"null",
		"amateur",
	}
	wordSource := NewFixedTestWordSource(words)
	tr := NewThesaurus(wordSource)

	tests := []struct {
		runes []rune
		want  []string
	}{
		{
			runes: []rune{'m', 'a', 'n', 'u', 'l'},
			want:  []string{"manual", "null"},
		},
		{
			runes: []rune{'m', 'a', 't', 'u', 'r', 'e', 'l', 'y', 'n'},
			want:  words,
		},
	}
	for _, tt := range tests {
		name := string(tt.runes)
		t.Run(name, func(t *testing.T) {
			result := tr.FindWordsWithJust(tt.runes...)
			words := make([]string, 0, len(tt.want))
			for word := range result {
				words = append(words, word)
			}
			assert.ElementsMatch(t, tt.want, words)

		})
	}
}
