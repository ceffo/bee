package bee_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"ceffo.com/bee/bee"
	"ceffo.com/bee/pkg/must"
)

func TestInput_New(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "no letters",
			wantErr: bee.ErrInvalidInputSize,
		},
		{
			name:    "not enough letters",
			input:   "A",
			wantErr: bee.ErrInvalidInputSize,
		},
		{
			name:    "duplicate letters",
			input:   "ABCDEFA",
			wantErr: bee.ErrDuplicateLetters,
		},
		{
			name:    "valid",
			input:   "ABCDEFG",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bee.NewInputFrom(tt.input)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestInput_IsPangram(t *testing.T) {
	tests := []struct {
		name  string
		input string
		word  string
		want  bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name:  "not pangram",
			input: "abcdefg",
			word:  "bcbcbc",
			want:  false,
		},
		{
			name:  "pangram same letters",
			input: "abcdefg",
			word:  "abcdefg",
			want:  true,
		},
		{
			name:  "pangram more letters",
			input: "abcdefg",
			word:  "abcdefgab",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := must.NoError(bee.NewInputFrom("abcdefg"))
			got := i.IsPangram(tt.word)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInput_IsExactPangram(t *testing.T) {
	tests := []struct {
		name  string
		input string
		word  string
		want  bool
	}{
		{
			name: "empty",
			want: false,
		},
		{
			name:  "not pangram",
			input: "abcdefg",
			word:  "bcbcbc",
			want:  false,
		},
		{
			name:  "pangram more letters",
			input: "abcdefg",
			word:  "abcdefgab",
			want:  false,
		},
		{
			name:  "exact pangram",
			input: "abcdefg",
			word:  "abcdefg",
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := must.NoError(bee.NewInputFrom("abcdefg"))
			got := i.IsExactPangram(tt.word)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInput_Score(t *testing.T) {
	tests := []struct {
		name string
		word string
		want int
	}{
		{
			name: "empty word",
			word: "",
			want: 0,
		},
		{
			name: "less than 4",
			word: "abc",
			want: 0,
		},
		{
			name: "4 letters",
			word: "abca",
			want: 1,
		},
		{
			name: "5 letters",
			word: "abcab",
			want: 5,
		},
		{
			name: "7 letters - not pangram",
			word: "ABCDEFA",
			want: 7,
		},
		{
			name: "7 letters - pangram",
			word: "ABCDEFG",
			want: 14,
		},
		{
			name: "7 letters - longer pangram",
			word: "ABCDEFGAA",
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := must.NoError(bee.NewInputFrom("abcdefg"))
			got := i.Score(tt.word)
			assert.Equal(t, tt.want, got)
		})
	}
}
