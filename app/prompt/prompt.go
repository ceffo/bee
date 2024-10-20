package prompt

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ceffo.com/bee/bee"
)

// Model is the model for the prompt
type Model struct {
	letters []rune
	errMsg  string
	valid   bool
	abort   bool
}

func NewModel() Model {
	return Model{}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) IsAborted() bool {
	return m.abort
}

func (m Model) ToBeeInput() (bee.Input, error) {
	if m.valid {
		return bee.NewBeeInput(m.letters[0], m.letters[1:]), nil
	}
	return bee.Input{}, fmt.Errorf("invalid input: %s", m.errMsg)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.abort = true
			return m, tea.Quit
		case tea.KeyRunes:
			if len(m.letters) == 7 {
				m.errMsg = "only 7 letters allowed"
				return m, nil
			}
			err := validateNewRune(msg.Runes[0], m.letters)
			if err != nil {
				m.errMsg = err.Error()
				return m, nil
			}
			m.errMsg = ""
			m.letters = append(m.letters, msg.Runes...)
			m.valid = len(m.letters) == 7 // only valid if 7 letters
		case tea.KeyBackspace:
			if len(m.letters) > 0 {
				m.letters = m.letters[:len(m.letters)-1]
			}
		case tea.KeyEnter:
			return m, tea.Quit
		}
	}

	return m, nil
}

func validateNewRune(r rune, letters []rune) error {
	if r < 'a' || r > 'z' {
		return fmt.Errorf("invalid rune: %c", r)
	}
	// must not be in the list of existing letters
	for _, l := range letters {
		if l == r {
			return fmt.Errorf("letter already entered: %c", r)
		}
	}
	return nil
}

const width = 20

var (
	promptStyle       = lipgloss.NewStyle().Width(width).Foreground(lipgloss.Color("#931e93"))
	centerLetterStyle = lipgloss.NewStyle().Width(width).Foreground(lipgloss.Color("#c8970e")).Bold(true)
	otherLetterStyle  = lipgloss.NewStyle().Width(width).Foreground(lipgloss.Color("#ede2c4")).Bold(false)
)

func (m Model) View() string {
	strLetters := ""
	for i, l := range m.letters {
		strL := strings.ToUpper(string(l))
		if i == 0 {
			strLetters += centerLetterStyle.Render(strL)
		} else {
			strLetters += " " + otherLetterStyle.Render(strL)
		}
	}
	strPrompt := promptStyle.Render("Enter 7 letters:")
	if m.valid {
		strPrompt = lipgloss.JoinHorizontal(lipgloss.Left, strPrompt,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render("valid input - press Enter to continue"))
	}
	if m.errMsg != "" {
		strPrompt = lipgloss.JoinHorizontal(lipgloss.Left, strPrompt,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(m.errMsg))
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		strPrompt,
		strLetters,
	)
}
