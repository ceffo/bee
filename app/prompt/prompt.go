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
	mapset "github.com/deckarep/golang-set/v2"

	"ceffo.com/bee/app/common"
	"ceffo.com/bee/app/palette"
	"ceffo.com/bee/bee"
	"ceffo.com/bee/pkg/slices"
)

const (
	errShowDuration = time.Millisecond * 1500
	timerTick       = time.Millisecond * 100
)

// DoneMsg is the message sent when the prompt is done
type DoneMsg struct {
	BeeInput *bee.Input
}

type promptState int

const (
	promptStatePrompt promptState = iota
	promptStateDone
)

// Model is the model for the prompt
type Model struct {
	letters     []rune
	input       *bee.Input
	state       promptState
	err         error
	errMsgTimer timer.Model
	timerID     int
}

// New creates a new prompt model
func New(letters ...rune) Model {
	m := Model{
		// create a stopped timer
		errMsgTimer: timer.NewWithInterval(0, timerTick),
		state:       promptStatePrompt,
	}
	if len(letters) <= bee.NumLetters {
		m.SetLetters(letters)
	}
	return m
}

type keyMap struct {
	remove key.Binding
	enter  key.Binding
	reset  key.Binding
	quit   key.Binding
}

// FullHelp returns the full help for the model
func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
	}
}

// ShortHelp returns the short help for the model
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

// IsInputValid returns whether the input is valid
func (m Model) IsInputValid() bool {
	return m.input != nil
}

func (m Model) keyMap() keyMap {
	km := keyMap{
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("^c", "quit"),
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
			key.WithHelp("␛", "reset"),
		),
	}
	hasLetters := len(m.letters) > 0
	km.reset.SetEnabled(hasLetters)
	km.remove.SetEnabled(hasLetters)
	km.enter.SetEnabled(m.IsInputValid())
	return km
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	cmds := common.NewMsgBatch(m.errMsgTimer.Init())
	return cmds.Cmd()
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
		m.SetLetters(append(m.letters, newRune))
	}
	if key.Matches(msg, km.remove) {
		if len(m.letters) > 0 {
			m.SetLetters(m.letters[:len(m.letters)-1])
		}
	}
	if key.Matches(msg, km.reset) {
		m.SetLetters(nil)
	}
	if key.Matches(msg, km.enter) {
		return m.onPromptValidated()
	}
	return m, nil
}

// SetLetters sets the letters for the prompt model, validates them, and updates the input field.
// If the letters are invalid, it logs an error and does not update the model.
func (m *Model) SetLetters(letters []rune) {
	validatedLetters, err := validateLetters(letters)
	if err != nil {
		log.Errorf("SetLetters: invalid letters %v: %v", letters, err)
		return
	}
	m.letters = validatedLetters
	m.input = nil
	if len(letters) == bee.NumLetters {
		input, err := bee.NewInput(letters...)
		if err != nil {
			log.Errorf("SetLetters: failed to create input for letters %v: %v", letters, err)
		} else {
			m.input = input
		}
	}
}

func (m *Model) onPromptValidated() (Model, tea.Cmd) {
	if m.input == nil {
		return *m, nil
	}
	m.state = promptStateDone
	return *m, common.ToCmd(DoneMsg{BeeInput: m.input})
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

func validateLetters(letters []rune) ([]rune, error) {
	letters = slices.Map(letters, unicode.ToUpper)
	lettersSet := mapset.NewSet(letters...)
	if lettersSet.Cardinality() != len(letters) {
		return nil, fmt.Errorf("duplicate letters")
	}
	for _, l := range letters {
		if l < 'A' || l > 'Z' {
			return nil, fmt.Errorf("not a letter between A and Z: %c", l)
		}
	}
	return letters, nil
}

var (
	baseStyle   = lipgloss.NewStyle()
	promptStyle = baseStyle.Inherit(palette.Prompt)
	errorStyle  = baseStyle.Inherit(palette.Error)
	// CenterLetterStyle is the style for the center letter
	CenterLetterStyle = baseStyle.Inherit(palette.Secondary).Bold(true)
	// NormalLetterStyle is the style for the normal letters
	NormalLetterStyle = baseStyle.Inherit(palette.Primary)
	// PangramLetterStyle is the style for the pangram letters
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
	if m.state == promptStateDone {
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
