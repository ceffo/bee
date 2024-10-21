package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"ceffo.com/bee/app/spellbee"
	"ceffo.com/bee/wordsource/reader"
)

func main() {
	err := runApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runApp() error {
	wordFile, err := os.Open("data/words.txt")
	if err != nil {
		panic(err)
	}
	defer wordFile.Close()

	model := spellbee.NewModel(reader.NewReaderSource(wordFile))
	opts := []tea.ProgramOption{}

	// are we debugging? don't go alt screen
	if os.Getenv("DEBUG") == "" {
		opts = append(opts, tea.WithAltScreen())
	}

	app := tea.NewProgram(model, opts...)
	_, err = app.Run()
	return err
}
