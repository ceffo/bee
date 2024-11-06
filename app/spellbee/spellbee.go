package spellbee

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"ceffo.com/bee/app/common"
	"ceffo.com/bee/app/palette"
	"ceffo.com/bee/app/prompt"
	"ceffo.com/bee/bee"
	"ceffo.com/bee/pkg/columntable"
	"ceffo.com/bee/pkg/keymap"
	"ceffo.com/bee/pkg/slices"
	"ceffo.com/bee/wordsource"
)

const (
	wordWidth    = 15
	scoreWidth   = 3
	maxItemWidth = wordWidth + scoreWidth
	headerHeight = 7
	helpHeight   = 1
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
	spinner    spinner.Model
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

func allKeyMap() keyMap {
	return keyMap{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/⌃+c", "quit"),
		),
		reset: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("⎋", "reset"),
		),
	}
}

func (m Model) keyMap() keyMap {
	km := allKeyMap()
	km.quit.SetEnabled(m.state != statePrompt)
	km.reset.SetEnabled(m.state != statePrompt)
	return km
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
	m.spinner = spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(
			lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		),
	)
	return m
}

type keyMap struct {
	quit  key.Binding
	reset key.Binding
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
		km.reset,
	}
	return bindings
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
	return tea.Batch(m.prompt.Init(), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	msgs := common.NewMsgBatch()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowSizeMsg(msg)
	case tea.KeyMsg:
		m, cmd = m.handleKeyMsg(msg)
		msgs.Add(cmd)
	case prompt.DoneMsg:
		msgs.Add(m.handlePromptDoneMsg(msg))
	case newResultMsg:
		msgs.Add(m.handleNewResultMsg(msg))
	case resultsDoneMsg:
		log.Info("Received results done message")
		m.state = stateRetrieved
	}

	switch m.state {
	case statePrompt:
		msgs.Add(m.updatePrompt(msg))
	case stateRetrieving:
		m.spinner, cmd = m.spinner.Update(msg)
		msgs.Add(cmd)
	case stateRetrieved:
		msgs.Add(m.updateColumnTable(msg))
	}

	return m, msgs.Cmd()
}

func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.width = msg.Width - 2
	m.height = msg.Height
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (Model, tea.Cmd) {
	km := m.keyMap()
	if key.Matches(msg, km.quit) {
		return m, tea.Quit
	}
	if key.Matches(msg, km.reset) {
		return m.reset(), nil
	}
	return m, nil
}

func (m *Model) handlePromptDoneMsg(msg prompt.DoneMsg) tea.Cmd {
	log.Info("Received prompt done message")
	if msg.Valid {
		input := msg.BeeInput
		m.state = stateRetrieving
		m.input = &input
		stream := m.solver.SolveFor(input)
		return tea.Batch(listenToResults(stream, input), m.spinner.Tick)
	}
	return nil
}

func (m *Model) handleNewResultMsg(msg newResultMsg) tea.Cmd {
	m.results = append(m.results, msg.result)
	renderedItems := renderResults(m.results, m.input)
	m.table.SetItems(renderedItems)
	return listenToResults(msg.stream, msg.input)
}

func (m *Model) updatePrompt(msg tea.Msg) tea.Cmd {
	newModel, cmd := m.prompt.Update(msg)
	m.prompt = newModel
	return cmd
}

func (m *Model) updateColumnTable(msg tea.Msg) tea.Cmd {
	m.table.SetSize(m.width, m.height-headerHeight-helpHeight)
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
		promptStyle = lipgloss.NewStyle().Align(lipgloss.Left).Margin(1, 0, 0, 1)
	}
	promptView := promptStyle.Render(m.prompt.View())
	elements = append(elements, promptView)

	headerView := ""
	switch m.state {
	case stateRetrieving:
		spinnerView := palette.Prompt.Render(m.spinner.View())
		headerView += fmt.Sprintf("words %03d %s", len(m.results), spinnerView)
	case stateRetrieved:
		totalScore := slices.FoldLeft(m.results, 0, func(acc int, r result) int {
			return acc + r.score
		})
		headerView += fmt.Sprintf("words %03d ▪︎ score ", len(m.results))
		headerView += palette.Prompt.Render(strconv.Itoa(totalScore))
	}
	if headerView != "" {
		headerView = lipgloss.NewStyle().Align(lipgloss.Left).MarginLeft(1).Render(headerView)
		elements = append(elements, headerView)
	}

	contentView := ""
	if m.state == stateRetrieved {
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

	bag := keymap.NewBag(m)
	switch m.state {
	case statePrompt:
		bag = bag.Add(m.prompt)
	case stateRetrieved:
		bag = bag.Add(m.table)
	}
	helpView := lipgloss.NewStyle().Margin(0, 1).Render(m.help.View(bag))
	elements = append(elements, helpView)

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
