package spellbee

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"ceffo.com/bee/app/common"
	"ceffo.com/bee/app/palette"
	"ceffo.com/bee/app/prompt"
	"ceffo.com/bee/bee"
	"ceffo.com/bee/wordsource"
)

const (
	maxWordLength   = 15
	scoreWidth      = 4
	columInterspace = " │ "
	resultHeight    = 20
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
	wordSourceMaker wordsource.SourceMaker
	solver          *bee.Solver

	prompt    prompt.Model
	paginator paginator.Model

	state   state
	input   *bee.Input
	results []result

	// window size
	width  int
	height int
}

func NewModel(wordSourceMaker wordsource.SourceMaker) Model {
	log.Info("Creating new spellbee model")
	return Model{
		state:           statePrompt,
		wordSourceMaker: wordSourceMaker,
		solver:          bee.NewSolver(wordSourceMaker()),
		prompt:          prompt.New(),
		paginator:       newPaginator(),
	}
}

func newPaginator() paginator.Model {
	pg := paginator.New()
	pg.Type = paginator.Dots
	pg.ActiveDot = palette.Secondary.Render("●")
	pg.InactiveDot = palette.Primary.Faint(true).Render("○")
	return pg
}

func (m Model) reset() Model {
	log.Info("Resetting spellbee model")

	m.state = statePrompt
	m.input = nil
	m.results = nil
	m.solver = bee.NewSolver(m.wordSourceMaker())
	m.prompt = prompt.New()
	m.paginator = newPaginator()
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
		log.Info("Updating window size")
		m.width = msg.Width - 2
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.state == stateRetrieved {
				// reset the model
				return m.reset(), nil
			}
			return m, tea.Quit
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
	case resultsDoneMsg:
		log.Info("Received results done message")
		m.state = stateRetrieved
	}

	switch m.state {
	case statePrompt:
		promptModel, promptCmd := m.prompt.Update(msg)
		m.prompt = promptModel.(prompt.Model)
		msgs.Add(promptCmd)
	case stateRetrieving, stateRetrieved:
		// update items per page
		columns := m.width / (maxWordLength + scoreWidth + len(columInterspace))
		numLines := len(m.results) / columns
		maxLines := resultHeight
		m.paginator.PerPage = maxLines
		m.paginator.SetTotalPages(numLines)

		// update paginator
		var paginatorCmd tea.Cmd
		m.paginator, paginatorCmd = m.paginator.Update(msg)
		msgs.Add(paginatorCmd)
	}

	return m, msgs.Cmd()
}

var (
	titleStyle = palette.Secondary.
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		Bold(true)
)

func (m Model) View() string {
	headerView := titleStyle.Width(m.width).Render("Spelling Bee")
	headerView = lipgloss.JoinVertical(lipgloss.Left, headerView, m.prompt.View())
	var contentView string
	contentWidth := m.width
	switch m.state {
	case stateRetrieving:
		headerView += fmt.Sprintf("\n\nRetrieving words... (%d found)", len(m.results))
	case stateRetrieved:
		headerView += fmt.Sprintf("\n\nFound %d words", len(m.results))
	}
	if m.state == stateRetrieving || m.state == stateRetrieved {
		headerHeight := strings.Count(headerView, "\n")
		contentHeight := resultHeight
		if headerHeight+contentHeight > m.height {
			contentView = "Window too small"
		} else {
			contentView = m.renderResults(contentWidth, contentHeight)
		}
	}
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Align(lipgloss.Left).Render(headerView),
		contentView,
	)
	return view
}

func (m Model) renderResults(width, height int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(width).
		MaxHeight(height + 1)

	if len(m.results) == 0 {
		return style.Render("No words found")
	}

	columns := m.width / (maxWordLength + scoreWidth + len(columInterspace))
	numLines := len(m.results) / columns

	startLine, endLine := m.paginator.GetSliceBounds(numLines)
	startIdx := startLine * columns
	endIdx := endLine*columns + columns
	if endIdx > len(m.results) {
		endIdx = len(m.results)
	}
	results := m.results[startIdx:endIdx]

	styleWord := palette.Primary.Width(maxWordLength).Align(lipgloss.Left)
	styleScore := palette.Prompt.Width(scoreWidth).Align(lipgloss.Right)
	xIdx := 0
	sb := strings.Builder{}
	for _, r := range results {
		word := renderWord(r.word, m.input.Center())
		score := strconv.Itoa(r.score)
		sb.WriteString(styleWord.Render(word))
		sb.WriteString(styleScore.Render(score))
		xIdx++
		if xIdx == columns {
			sb.WriteString("\n")
			xIdx = 0
		} else {
			sb.WriteString(columInterspace)
		}
	}
	render := style.Render(sb.String())
	if m.paginator.TotalPages > 1 {
		render = lipgloss.JoinVertical(
			lipgloss.Center,
			render,
			lipgloss.NewStyle().MarginTop(1).Align(lipgloss.Center).Render(m.paginator.View()))
	}
	return render
}

func renderWord(word string, highlightedChar rune) string {
	highlightedChar = unicode.ToUpper(highlightedChar)
	sb := strings.Builder{}
	for _, l := range word {
		if l == highlightedChar {
			sb.WriteString(prompt.CenterLetterStyle.Bold(true).Render(string(l)))
		} else {
			sb.WriteString(prompt.OtherLetterStyle.Render(string(l)))
		}
	}
	return sb.String()
}
