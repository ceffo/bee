package spellbee

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"ceffo.com/bee/app/common"
	"ceffo.com/bee/app/palette"
	"ceffo.com/bee/app/prompt"
	"ceffo.com/bee/bee"
	"ceffo.com/bee/pkg/columntable"
	"ceffo.com/bee/pkg/slices"
	"ceffo.com/bee/wordsource"
)

const (
	wordWidth    = 15
	scoreWidth   = 3
	maxItemWidth = wordWidth + scoreWidth
	headerHeight = 7
)

type result struct {
	word  string
	score int
}

type state int

const (
	statePrompt state = iota
	stateRetrieving
	stateRetrieved
)

type Model struct {
	wordSource wordsource.Maker
	solver     *bee.Solver
	prompt     prompt.Model
	table      columntable.Model
	help       help.Model
	state      state
	input      *bee.Input
	results    []result
	width      int
	height     int
}

func NewModel(wsm wordsource.Maker) Model {
	log.Info("Creating new spellbee model")
	return Model{
		wordSource: wsm,
	}.reset()
}

func newColumnTable() columntable.Model {
	return columntable.New(
		columntable.WithDotPaginator(
			palette.Secondary.Render("●"),
			palette.Primary.Faint(true).Render("○"),
		),
		columntable.WithItemWidth(maxItemWidth),
	)
}

func (m Model) reset() Model {
	log.Info("Resetting spellbee model")

	m.state = statePrompt
	m.input = nil
	m.results = nil
	m.solver = bee.NewSolver(m.wordSource())
	m.prompt = prompt.New()
	m.table = newColumnTable()
	m.help = help.New()
	return m
}

type newResultMsg struct {
	stream wordsource.Stream
	input  bee.Input
	result result
}

type resultsDoneMsg struct{}

func listenToResults(stream wordsource.Stream, input bee.Input) tea.Cmd {
	log.Debug("Listening to results")
	return func() tea.Msg {
		word, ok := <-stream
		if !ok {
			return resultsDoneMsg{}
		}
		return newResultMsg{
			result: result{
				word:  word,
				score: input.Score(word),
			},
			input:  input,
			stream: stream,
		}
	}
}

func (m Model) Init() tea.Cmd {
	log.Info("Initializing spellbee model")
	return tea.Batch(m.prompt.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	msgs := common.NewMsgBatch()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			if msg.String() == "q" {
				return m, tea.Quit
			}
		case tea.KeyEsc:
			switch m.state {
			case statePrompt:
				return m, tea.Quit
			default:
				return m.reset(), nil
			}
		case tea.KeyBackspace:
			if m.state == stateRetrieved {
				// reset the model
				return m.reset(), nil
			}
		}
	case prompt.DoneMsg:
		log.Info("Received prompt done message")
		if msg.Valid {
			input := msg.BeeInput
			m.state = stateRetrieving
			m.input = &input
			stream := m.solver.SolveFor(input)
			msgs.Add(listenToResults(stream, input))
		}
	case newResultMsg:
		m.results = append(m.results, msg.result)
		msgs.Add(listenToResults(msg.stream, msg.input))
		renderedItems := renderResults(m.results, m.input)
		m.table.SetItems(renderedItems)
	case resultsDoneMsg:
		log.Info("Received results done message")
		m.state = stateRetrieved
	}

	switch m.state {
	case statePrompt:
		msgs.Add(m.updatePrompt(msg))
	case stateRetrieving, stateRetrieved:
		msgs.Add(m.updateColumnTable(msg))
	}

	return m, msgs.Cmd()
}

func (m *Model) updatePrompt(msg tea.Msg) tea.Cmd {
	newModel, cmd := m.prompt.Update(msg)
	m.prompt = newModel
	return cmd
}

func (m *Model) updateColumnTable(msg tea.Msg) tea.Cmd {
	m.table.SetSize(m.width, m.height-headerHeight)
	newModel, cmd := m.table.Update(msg)
	m.table = newModel
	return cmd
}

func renderResults(results []result, input *bee.Input) []string {
	return slices.Map(results, func(r result) string {
		return renderResult(r, input)
	})
}

var (
	titleStyle = palette.Secondary.
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		Bold(true)
)

func (m Model) View() string {
	elements := []string{}
	titleView := titleStyle.Width(m.width).Render("Bee Solver")
	elements = append(elements, titleView)

	promptStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width)
	if m.state == statePrompt {
		promptStyle = lipgloss.NewStyle().Align(lipgloss.Left).Margin(1, 1)
	}
	promptView := promptStyle.Render(m.prompt.View())
	elements = append(elements, promptView)

	headerView := ""
	switch m.state {
	case stateRetrieving:
		headerView += fmt.Sprintf("retrieving words... (%d found)", len(m.results))
	case stateRetrieved:
		totalScore := slices.FoldLeft(m.results, func(acc int, r result) int {
			return acc + r.score
		}, 0)
		headerView += fmt.Sprintf("words %d ▪︎ score ", len(m.results))
		headerView += palette.Prompt.Render(strconv.Itoa(totalScore))
	}
	if headerView != "" {
		headerView = lipgloss.NewStyle().Align(lipgloss.Left).MarginLeft(1).Render(headerView)
		elements = append(elements, headerView)
	}

	contentView := ""
	if m.state == stateRetrieving || m.state == stateRetrieved {
		headerHeight := lipgloss.Height(headerView)
		tableView := m.table.View()
		contentHeight := lipgloss.Height(tableView)
		if headerHeight+contentHeight > m.height {
			contentView = "Window too small"
		} else {
			contentView = tableView
		}
	}
	if contentView != "" {
		elements = append(elements, contentView)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		elements...,
	)
	return view
}

func renderResult(r result, input *bee.Input) string {
	styleWord := palette.Primary.Width(wordWidth).Align(lipgloss.Left)
	styleScore := palette.Prompt.Width(scoreWidth).Align(lipgloss.Right)
	word := renderWord(r.word, input)
	score := strconv.Itoa(r.score)
	return styleWord.Render(word) + styleScore.Render(score)
}

func renderWord(word string, input *bee.Input) string {
	letterStyle := prompt.NormalLetterStyle
	centerStyle := prompt.CenterLetterStyle
	if input.IsPangram(word) {
		letterStyle = prompt.PangramLetterStyle
	}
	center := unicode.ToUpper(input.Center())
	sb := strings.Builder{}
	for _, l := range word {
		if l == center {
			sb.WriteString(centerStyle.Render(string(l)))
		} else {
			sb.WriteString(letterStyle.Render(string(l)))
		}
	}
	return sb.String()
}
