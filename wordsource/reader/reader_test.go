package reader

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderSource_GetWords(t *testing.T) {
	type fields struct {
		reader io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "empty",
			fields: fields{
				reader: strings.NewReader(""),
			},
			want: []string{},
		},
		{
			name: "one word",
			fields: fields{
				reader: strings.NewReader("hello\n"),
			},
			want: []string{"hello"},
		},
		{
			name: "two words",
			fields: fields{
				reader: strings.NewReader("hello\nworld\n"),
			},
			want: []string{"hello", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := ReaderSource{
				reader: tt.fields.reader,
			}
			got := rs.GetWords()
			var words []string
			for word := range got {
				words = append(words, word)
			}
			assert.ElementsMatch(t, tt.want, words)
		})
	}
}
