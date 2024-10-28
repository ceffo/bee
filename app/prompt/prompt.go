package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ceffo.com/bee/app/common"
	"ceffo.com/bee/app/palette"
	"ceffo.com/bee/bee"
)

const (
	errShowDuration = time.Millisecond * 1500
	timerTick       = time.Millisecond * 100
)

// DoneMsg is the message sent when the prompt is done
type DoneMsg struct {
	Valid    bool
	BeeInput bee.Input
}

// Model is the model for the prompt
type Model struct {
	letters []rune
	valid   bool

	err         error
	errMsgTimer timer.Model
	timerID     int
	done        bool
}

// New creates a new prompt model
func New() Model {
	return Model{
		// create a stopped timer
		errMsgTimer: timer.NewWithInterval(0, timerTick),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.errMsgTimer.Init()
}

func (m *Model) triggerError(err error) tea.Cmd {
	m.err = err
	m.errMsgTimer = timer.NewWithInterval(errShowDuration, timerTick)
	m.timerID = m.errMsgTimer.ID()
	return m.errMsgTimer.Start()
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := common.NewMsgBatch()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			if len(m.letters) == bee.NumLetters {
				return m, nil
			}
			err := validateNewRune(msg.Runes[0], m.letters)
			if err != nil {
				cmds.Add(m.triggerError(err))
			} else {
				m.letters = append(m.letters, msg.Runes...)
				m.valid = len(m.letters) == bee.NumLetters
			}

		case tea.KeyBackspace:
			if len(m.letters) > 0 {
				m.letters = m.letters[:len(m.letters)-1]
				m.valid = false
			}

		case tea.KeyEnter:
			if m.valid {
				m.done = true
				return m, m.promptValidated
			}
		}
	case timer.TimeoutMsg:
		// clear the error message if the timer matches the last one
		if msg.ID == m.timerID {
			m.err = nil
		}
	}

	// update the error message timer
	timerModel, timerCmd := m.errMsgTimer.Update(msg)
	cmds.Add(timerCmd)
	m.errMsgTimer = timerModel

	return m, cmds.Cmd()
}

func (m Model) promptValidated() tea.Msg {
	var input bee.Input
	if m.valid {
		input = bee.NewBeeInput(m.letters[0], m.letters[1:])
	}
	return DoneMsg{Valid: m.valid, BeeInput: input}
}

func validateNewRune(r rune, letters []rune) error {
	if r < 'a' || r > 'z' {
		return fmt.Errorf("not a letter between a and z: %c", r)
	}
	// must not be in the list of existing letters
	for _, l := range letters {
		if l == r {
			return fmt.Errorf("letter already entered: %c", r)
		}
	}
	return nil
}

var (
	baseStyle         = lipgloss.NewStyle().Align(lipgloss.Left)
	promptStyle       = baseStyle.Inherit(palette.Prompt)
	doneStyle         = baseStyle.Inherit(palette.Positive)
	errorStyle        = baseStyle.Inherit(palette.Error)
	CenterLetterStyle = baseStyle.Inherit(palette.Secondary).Bold(true)
	OtherLetterStyle  = baseStyle.Inherit(palette.Primary)
)

// View returns the view for the model
func (m Model) View() string {
	strLetters := ""
	for i, l := range m.letters {
		strL := strings.ToUpper(string(l))
		if i == 0 {
			strLetters += CenterLetterStyle.Render(strL)
		} else {
			strLetters += " " + OtherLetterStyle.Render(strL)
		}
	}
	if m.done {
		return strLetters
	}
	strError := ""
	strHeader := promptStyle.Render("Enter 7 letters")
	if m.valid {
		strHeader = lipgloss.JoinHorizontal(lipgloss.Left, strHeader,
			doneStyle.Margin(0, 2).Render("press Enter to start solving"))
	}
	if m.err != nil {
		strError = errorStyle.Render(m.err.Error())
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		strHeader,
		strLetters,
		strError,
	)
}
