package bee

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInput_IsPangram(t *testing.T) {
	type fields struct {
		center  rune
		letters []rune
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "empty",
			fields: fields{
				letters: []rune{'a', 'b', 'c'},
			},
			args: args{
				word: "",
			},
			want: false,
		},
		{
			name: "not pangram",
			fields: fields{
				letters: []rune{'a', 'b', 'c'},
			},
			args: args{
				word: "bcbcbc",
			},
			want: false,
		},
		{
			name: "pangram same letters",
			fields: fields{
				letters: []rune{'a', 'b', 'c'},
			},
			args: args{
				word: "abc",
			},
			want: true,
		},
		{
			name: "pangram more letters",
			fields: fields{
				letters: []rune{'a', 'b', 'c'},
			},
			args: args{
				word: "abccba",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Input{
				center:  tt.fields.center,
				letters: tt.fields.letters,
			}
			got := i.IsPangram(tt.args.word)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInput_Score(t *testing.T) {
	type fields struct {
		center  rune
		letters []rune
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "empty",
			fields: fields{
				letters: []rune{'a', 'b', 'c'},
			},
			args: args{
				word: "",
			},
			want: 0,
		},
		{
			name: "less than 4",
			fields: fields{
				letters: []rune{'a', 'b', 'c'},
			},
			args: args{
				word: "abc",
			},
			want: 0,
		},
		{
			name: "4 letters",
			fields: fields{
				letters: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
			},
			args: args{
				word: "abca",
			},
			want: 1,
		},
		{
			name: "5 letters",
			fields: fields{
				letters: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
			},
			args: args{
				word: "abcab",
			},
			want: 5,
		},
		{
			name: "7 letters - not pangram",
			fields: fields{
				letters: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
			},
			args: args{
				word: "abcdefa",
			},
			want: 7,
		},
		{
			name: "7 letters - pangram",
			fields: fields{
				letters: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
			},
			args: args{
				word: "abcdefg",
			},
			want: 14,
		},
		{
			name: "7 letters - longer pangram",
			fields: fields{
				letters: []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g'},
			},
			args: args{
				word: "abcdefgaa",
			},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Input{
				center:  tt.fields.center,
				letters: tt.fields.letters,
			}
			got := i.Score(tt.args.word)
			assert.Equal(t, tt.want, got)
		})
	}
}
