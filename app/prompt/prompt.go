package prompt

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	done    bool

	err         error
	errMsgTimer timer.Model
	timerID     int

	keyMap promptKeyMap
	help   help.Model
}

// New creates a new prompt model
func New() Model {
	return Model{
		// create a stopped timer
		errMsgTimer: timer.NewWithInterval(0, timerTick),
		keyMap:      defaultKeyMap(),
		help:        help.New(),
	}
}

func defaultKeyMap() promptKeyMap {
	return promptKeyMap{
		remove: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("⌫", "Remove last letter"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↩", "Start solving"),
		),
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
		return fmt.Errorf("not a letter between A and Z: %c", r)
	}
	// must not be in the list of existing letters
	for _, l := range letters {
		if l == r {
			return fmt.Errorf("letter already entered: %c", unicode.ToUpper(r))
		}
	}
	return nil
}

type promptKeyMap struct {
	remove key.Binding
	enter  key.Binding
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
	}
}

func (m Model) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		m.keyMap.remove,
	}
	if m.valid {
		bindings = append(bindings, m.keyMap.enter)
	}
	return bindings
}

var (
	// helpStyle          = baseStyle.Inherit(palette.Help)
	baseStyle          = lipgloss.NewStyle()
	promptStyle        = baseStyle.Inherit(palette.Prompt)
	errorStyle         = baseStyle.Inherit(palette.Error)
	CenterLetterStyle  = baseStyle.Inherit(palette.Secondary).Bold(true)
	NormalLetterStyle  = baseStyle.Inherit(palette.Primary)
	PangramLetterStyle = baseStyle.Inherit(palette.Tertiary)
)

// View returns the view for the model
func (m Model) View() string {
	strLetters := ""
	for i := range 7 {
		if i >= len(m.letters) {
			if i > 0 {
				strLetters += " "
			}
			strLetters += "_"
			continue
		}
		l := m.letters[i]
		strL := strings.ToUpper(string(l))
		if i == 0 {
			strLetters += CenterLetterStyle.Render(strL)
		} else {
			strLetters += " " + NormalLetterStyle.Render(strL)
		}
	}
	if m.done {
		return strLetters
	}
	strInput := strLetters
	if m.err != nil {
		strError := errorStyle.Margin(0, 2).Render(m.err.Error())
		strInput = lipgloss.JoinHorizontal(lipgloss.Top, strInput, strError)
	}

	strPrompt := promptStyle.Render("Enter 7 letters")

	content := lipgloss.JoinHorizontal(lipgloss.Top,
		strPrompt,
		lipgloss.NewStyle().Margin(0, 2).Render(strInput),
	)
	content = lipgloss.JoinVertical(lipgloss.Left, content, m.help.View(m))
	return lipgloss.NewStyle().Margin(1, 0).Render(content)
}
