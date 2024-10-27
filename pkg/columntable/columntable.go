// colomntable is a package that provides a simple way to print a table with columns
package columntable

import (
	"strings"
	"unicode/utf8"

	"ceffo.com/bee/app/common"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	extraHeight      = 1 // extra height for the paginator
	defaultSeparator = ".│."
)

var (
	tableStyle = lipgloss.NewStyle()
)

type Model struct {
	width      int
	height     int
	items      []string
	itemWidth  int
	numColumns int
	separator  string

	paginator paginator.Model
}

// events

type SetSize struct {
	Width  int
	Height int
}

type Option func(*Model)

func WithSeparator(s string) Option {
	return func(m *Model) {
		m.separator = s
	}
}

func WithDotPaginator(activeDot, inactiveDot string) Option {
	return func(m *Model) {
		m.paginator.Type = paginator.Dots
		m.paginator.ActiveDot = activeDot
		m.paginator.InactiveDot = inactiveDot
	}
}

// New creates a new columntable model
func New(opts ...Option) *Model {
	pg := paginator.New()

	model := &Model{
		paginator: pg,
		separator: defaultSeparator,
	}
	for _, opt := range opts {
		opt(model)
	}
	return model
}

func (m *Model) SetItems(items []string, itemsWidth int) {
	m.items = items
	m.itemWidth = itemsWidth
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	msgs := common.NewMsgBatch()

	switch msg := msg.(type) {
	case SetSize:
		m.width = max(msg.Width, 0)
		m.height = max((msg.Height - extraHeight), 0)
		if m.width != 0 || m.height != 0 {
			m.updatePaginator()
		}
	}

	// update paginator
	paginatorModel, paginatorCmd := m.paginator.Update(msg)
	m.paginator = paginatorModel
	msgs.Add(paginatorCmd)

	return m, nil
}

func (m *Model) updatePaginator() {
	s := utf8.RuneCountInString(defaultSeparator)
	m.numColumns = (m.width + s) / (m.itemWidth + s)

	if m.numColumns == 0 {
		m.paginator.PerPage = 0
		m.paginator.SetTotalPages(0)
		return
	}

	totalRows := len(m.items) / m.numColumns
	// setup paginator
	m.paginator.PerPage = m.height - extraHeight
	m.paginator.SetTotalPages(totalRows)
}

/*
item1 │ item2 │ item3
item4 │ item5 │ item6
item7 │       │
*/

// View renders the model
func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.renderTable(),
		m.paginator.View(),
	)
}

func (m Model) renderTable() string {
	if len(m.items) == 0 || m.width == 0 || m.height == 0 || m.numColumns == 0 {
		return ""
	}
	totalRows := len(m.items) / m.numColumns
	// add one row if there are remaining items
	if len(m.items)%m.numColumns > 0 {
		totalRows++
	}
	rowStart, rowEnd := m.paginator.GetSliceBounds(totalRows)
	itemStart := rowStart * m.numColumns
	itemEnd := min((1+rowEnd)*m.numColumns, len(m.items))
	items := m.items[itemStart:itemEnd]
	nbCells := m.numColumns * m.paginator.PerPage

	sb := strings.Builder{}
	blank := strings.Repeat(" ", m.itemWidth)
	for i, cell := range iterateItems(items, blank, nbCells) {
		if i > 0 && i%m.numColumns == 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(cell)
		x := i % m.numColumns
		if x < m.numColumns-1 {
			sb.WriteString(defaultSeparator)
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

func iterateItems(items []string, padding string, numItems int) iterator {
	return pad_iter(items, padding, numItems)
}

type iterator func(yield func(int, string) bool)

func pad_iter(items []string, padding string, numItems int) iterator {
	return func(yield func(int, string) bool) {
		i := 0
		numFromItems := min(len(items), numItems)
		for i = range numFromItems {
			if !yield(i, items[i]) {
				return
			}
		}
		if numFromItems == numItems {
			return
		}
		for i = numFromItems; i < numItems; i++ {
			if !yield(i, padding) {
				return
			}
		}
	}
}
