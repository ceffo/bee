package prompt

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

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
	BeeInput *bee.Input
}

// Model is the model for the prompt
type Model struct {
	letters []rune
	valid   bool
	done    bool

	err         error
	errMsgTimer timer.Model
	timerID     int
}

// New creates a new prompt model
func New() Model {
	return Model{
		// create a stopped timer
		errMsgTimer: timer.NewWithInterval(0, timerTick),
	}
}

type keyMap struct {
	remove key.Binding
	enter  key.Binding
	reset  key.Binding
	quit   key.Binding
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
	}
}

func (m Model) ShortHelp() []key.Binding {
	km := m.keyMap()
	bindings := []key.Binding{
		km.quit,
		km.remove,
		km.reset,
		km.enter,
	}
	return bindings
}

func (m Model) keyMap() keyMap {
	km := keyMap{
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("⌃+c", "quit"),
		),
		remove: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("⌫", "delete"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↩", "solve"),
		),
		reset: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("⎋", "reset"),
		),
	}
	hasLetters := len(m.letters) > 0
	km.reset.SetEnabled(hasLetters)
	km.remove.SetEnabled(hasLetters)
	km.enter.SetEnabled(m.valid)
	return km
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.errMsgTimer.Init()
}

func (m Model) triggerError(err error) (Model, tea.Cmd) {
	m.err = err
	m.errMsgTimer = timer.NewWithInterval(errShowDuration, timerTick)
	m.timerID = m.errMsgTimer.ID()
	cmd := m.errMsgTimer.Start()
	return m, cmd
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := common.NewMsgBatch()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd = m.handleKeyMsg(msg)
		cmds.Add(cmd)
	case timer.TimeoutMsg:
		// clear the error message if the timer matches the last one
		if msg.ID == m.timerID {
			m.err = nil
		}
	}

	// update the error message timer
	m.errMsgTimer, cmd = m.errMsgTimer.Update(msg)
	cmds.Add(cmd)

	return m, cmds.Cmd()
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	km := m.keyMap()
	if key.Matches(msg, km.quit) {
		return m, tea.Quit
	}
	if msg.Type == tea.KeyRunes {
		if len(m.letters) == bee.NumLetters {
			return m, nil
		}
		newRune := unicode.ToUpper(msg.Runes[0])
		err := validateNewRune(newRune, m.letters)
		if err != nil {
			return m.triggerError(err)
		}
		m.letters = append(m.letters, newRune)
		m.valid = len(m.letters) == bee.NumLetters
	}
	if key.Matches(msg, km.remove) {
		if len(m.letters) > 0 {
			m.letters = m.letters[:len(m.letters)-1]
			m.valid = false
		}
	}
	if key.Matches(msg, km.reset) {
		m.letters = nil
		m.valid = false
		m.done = false
	}
	if key.Matches(msg, km.enter) {
		if m.valid {
			m.done = true
			return m, m.promptValidated
		}
	}
	return m, nil
}

func (m Model) promptValidated() tea.Msg {
	var input *bee.Input
	if m.valid {
		var err error
		input, err = bee.NewInput(m.letters[0], m.letters[1:])
		if err != nil {
			log.Errorf("error creating input: %v", err)
			return DoneMsg{Valid: false}
		}
	}
	return DoneMsg{Valid: m.valid, BeeInput: input}
}

func validateNewRune(r rune, letters []rune) error {
	if r < 'A' || r > 'Z' {
		return fmt.Errorf("not a letter between A and Z: %c", r)
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
	for i := range bee.NumLetters {
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

	strPrompt := promptStyle.Render(fmt.Sprintf("Enter %d letters", bee.NumLetters))

	content := lipgloss.JoinHorizontal(lipgloss.Top,
		strPrompt,
		lipgloss.NewStyle().Margin(0, 2).Render(strInput),
	)
	return content
}
