package main

import (
	"fmt"
	"os"

	"ceffo.com/bee/bee"
	"ceffo.com/bee/wordsource/reader"
)

func main() {
	input := bee.NewBeeInput('e', []rune{'c', 'h', 'o', 'u', 'l', 'n'})

	wordFile, err := os.Open("data/words.txt")
	if err != nil {
		panic(err)
	}
	defer wordFile.Close()
	beesolve := bee.NewSolver(reader.NewReaderSource(wordFile))
	for word := range beesolve.SolveFor(input) {
		fmt.Printf("%s %d\n", word, input.Score(word))
	}
}
