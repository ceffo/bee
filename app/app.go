package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"

	"ceffo.com/bee/app/spellbee"
	"ceffo.com/bee/bee"
	"ceffo.com/bee/core"
)

type App struct {
	core *core.Core
}

func New(c *core.Core) *App {
	return &App{c}
}

func (a *App) Run() error {
	source := a.core.Source()
	solver := bee.NewSolver(source)
	model := spellbee.NewModel(solver)
	opts := []tea.ProgramOption{}

	opts = append(opts, tea.WithAltScreen())

	app := tea.NewProgram(model, opts...)
	_, err := app.Run()
	if err != nil {
		log.Errorf("Error running app: %v", err)
	}
	return err
}
