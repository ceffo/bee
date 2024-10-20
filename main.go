package main

import (
	"fmt"
	"os"

	"ceffo.com/bee/app/prompt"
	"ceffo.com/bee/bee"
	"ceffo.com/bee/wordsource/reader"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	inputModel := prompt.NewModel()
	app := tea.NewProgram(inputModel)
	m, err := app.Run()
	if err != nil {
		panic(err)
	}
	promptModel := m.(prompt.Model)
	if promptModel.IsAborted() {
		fmt.Printf("Aborted\n")
		return
	}

	input, err := promptModel.ToBeeInput()
	if err != nil {
		panic(err)
	}

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
