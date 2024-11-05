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
				letters: []rune{'A', 'B', 'C'},
			},
			args: args{
				word: "",
			},
			want: false,
		},
		{
			name: "not pangram",
			fields: fields{
				letters: []rune{'A', 'B', 'C'},
			},
			args: args{
				word: "BCBCBC",
			},
			want: false,
		},
		{
			name: "pangram same letters",
			fields: fields{
				letters: []rune{'A', 'B', 'C'},
			},
			args: args{
				word: "ABC",
			},
			want: true,
		},
		{
			name: "pangram more letters",
			fields: fields{
				letters: []rune{'A', 'B', 'C'},
			},
			args: args{
				word: "ABCCBA",
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
				letters: []rune{'A', 'B', 'C'},
			},
			args: args{
				word: "",
			},
			want: 0,
		},
		{
			name: "less than 4",
			fields: fields{
				letters: []rune{'A', 'B', 'C'},
			},
			args: args{
				word: "ABC",
			},
			want: 0,
		},
		{
			name: "4 letters",
			fields: fields{
				letters: []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'},
			},
			args: args{
				word: "ABCA",
			},
			want: 1,
		},
		{
			name: "5 letters",
			fields: fields{
				letters: []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'},
			},
			args: args{
				word: "ABCAB",
			},
			want: 5,
		},
		{
			name: "7 letters - not pangram",
			fields: fields{
				letters: []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'},
			},
			args: args{
				word: "ABCDEFA",
			},
			want: 7,
		},
		{
			name: "7 letters - pangram",
			fields: fields{
				letters: []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'},
			},
			args: args{
				word: "ABCDEFG",
			},
			want: 14,
		},
		{
			name: "7 letters - longer pangram",
			fields: fields{
				letters: []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G'},
			},
			args: args{
				word: "ABCDEFGAA",
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
