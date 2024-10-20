package main

import (
	"fmt"
	"os"

	"ceffo.com/bee/beesolve"
	"ceffo.com/bee/wordsource/reader"
	mapset "github.com/deckarep/golang-set/v2"
)

func main() {
	input := beesolve.NewBeeInput('e', []rune{'c', 'h', 'o', 'u', 'l', 'n'})

	wordFile, err := os.Open("data/words.txt")
	if err != nil {
		panic(err)
	}
	defer wordFile.Close()
	beesolve := beesolve.NewBeeSolve(reader.NewReaderSource(wordFile))
	seen := mapset.NewSet[string]()
	for word := range beesolve.SolveFor(input) {
		if !seen.Contains(word) {
			fmt.Printf("%s %d\n", word, input.Score(word))
			seen.Add(word)
		}
	}
}
