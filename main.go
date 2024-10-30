package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"ceffo.com/bee/app/spellbee"
	"ceffo.com/bee/wordsource"
	"ceffo.com/bee/wordsource/reader"
)

type cleanupFunc func()

type config struct {
	WordListFileName string
	LogFileName      string
}

func defaultConfig() *config {
	return &config{
		WordListFileName: "data/words.txt",
		LogFileName:      "bee.log",
	}
}

func main() {
	err := runApp(defaultConfig())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func runApp(config *config) error {
	logCleanup, err := setupLogging(config)
	if err != nil {
		return err
	}
	log.Info("Starting app")
	defer logCleanup()
	fileCleanup, sourceMaker, err := sourceMaker(config.WordListFileName)
	if err != nil {
		return err
	}
	defer fileCleanup()

	model := spellbee.NewModel(sourceMaker)
	opts := []tea.ProgramOption{}

	opts = append(opts, tea.WithAltScreen())

	app := tea.NewProgram(model, opts...)
	_, err = app.Run()
	if err != nil {
		log.Errorf("Error running app: %v", err)
	} else {
		log.Info("App exited without error")
	}
	return err
}

func setupLogging(config *config) (cleanupFunc, error) {
	f, err := os.OpenFile(config.LogFileName, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	log.SetFormatter(log.JSONFormatter)
	log.SetTimeFormat(time.RFC3339Nano)
	return func() { f.Close() }, nil
}

func sourceMaker(wordlistFileName string) (cleanupFunc, wordsource.Maker, error) {
	log.Infof("Opening wordlist file %s", wordlistFileName)
	wordFile, err := os.Open(wordlistFileName)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { wordFile.Close() }
	maker := func() wordsource.Source {
		log.Infof("Creating new word source from %s", wordlistFileName)
		wordFile.Seek(0, 0)
		return reader.NewReaderSource(wordFile)
	}
	return cleanup, maker, nil
}
