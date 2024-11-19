package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"ceffo.com/bee/app/spellbee"
	"ceffo.com/bee/bee"
	"ceffo.com/bee/core"
)

// App is the main application struct
type App struct {
	core *core.Core
}

// New creates a new App
func New(c *core.Core) *App {
	return &App{c}
}

// Config is the configuration for the App
type Config struct {
	Input string
}

// Run runs the App
func (a *App) Run(config *Config) error {
	source := a.core.Source()
	solver := bee.NewSolver(source)

	modelOpts := []spellbee.Option{}
	if config.Input != "" {
		modelOpts = append(modelOpts, spellbee.WithInput(config.Input))
	}
	model, err := spellbee.New(solver, modelOpts...)
	if err != nil {
		return err
	}
	programOpts := []tea.ProgramOption{}

	programOpts = append(programOpts, tea.WithAltScreen())

	app := tea.NewProgram(model, programOpts...)
	_, err = app.Run()
	if err != nil {
		log.Errorf("Error running app: %v", err)
	}
	return err
}
