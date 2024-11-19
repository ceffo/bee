// colomntable is a package that provides a simple way to print a table with columns
package columntable

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ceffo.com/bee/app/common"
)

const (
	extraHeight      = 1 // extra height for the paginator
	defaultSeparator = " │ "
)

// Model is the model for the columntable
type Model struct {
	width      int
	height     int
	items      []string
	itemWidth  int
	numColumns int
	separator  string
	paginator  paginator.Model
}

func (m Model) keyMap() paginator.KeyMap {
	km := m.paginator.KeyMap
	km.PrevPage.SetHelp("←/h", "previous")
	km.NextPage.SetHelp("→/l", "next")
	multiPage := m.paginator.TotalPages > 1
	km.PrevPage.SetEnabled(multiPage && !m.paginator.OnFirstPage())
	km.NextPage.SetEnabled(multiPage && !m.paginator.OnLastPage())
	return km
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
		km.PrevPage,
		km.NextPage,
	}
	return bindings
}

// Option is a functional option for the Model
type Option func(*Model)

// WithSeparator sets the separator for the columntable
func WithSeparator(s string) Option {
	return func(m *Model) {
		m.separator = s
	}
}

// WithDotPaginator sets the paginator to use dots
func WithDotPaginator(activeDot, inactiveDot string) Option {
	return func(m *Model) {
		m.paginator.Type = paginator.Dots
		m.paginator.ActiveDot = activeDot
		m.paginator.InactiveDot = inactiveDot
	}
}

// WithItemWidth sets the width of the items in the table
func WithItemWidth(width int) Option {
	return func(m *Model) {
		m.itemWidth = width
	}
}

// New creates a new columntable model
func New(opts ...Option) Model {
	pg := paginator.New()

	model := Model{
		paginator: pg,
		separator: defaultSeparator,
	}
	for _, opt := range opts {
		opt(&model)
	}
	return model
}

// SetItemWidth sets the width of the items in the table
func (m *Model) SetItemWidth(width int) {
	m.itemWidth = width
}

// SetItems sets the items in the table
func (m *Model) SetItems(items []string) {
	m.items = items
}

// Init initializes the model
func (Model) Init() tea.Cmd {
	return nil
}

// Update updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	msgs := common.NewMsgBatch()
	msgs.Add(m.updatePaginator(msg))
	return m, nil
}

// SetSize sets the size of the table
func (m *Model) SetSize(width, height int) {
	m.width = max(width, 0)
	m.height = max(height, 0)
}

func (m *Model) updatePaginator(msg tea.Msg) tea.Cmd {
	m.calcPaginator()
	paginatorModel, paginatorCmd := m.paginator.Update(msg)
	m.paginator = paginatorModel
	return paginatorCmd
}

func (m *Model) calcPaginator() {
	sepWidth := lipgloss.Width(m.separator)
	m.numColumns = (m.width + sepWidth) / (m.itemWidth + sepWidth)

	if m.numColumns == 0 {
		m.paginator.PerPage = 0
		m.paginator.SetTotalPages(0)
		return
	}

	totalRows := len(m.items) / m.numColumns
	m.paginator.PerPage = m.height - extraHeight
	m.paginator.SetTotalPages(totalRows)
}

// View renders the model
func (m Model) View() string {
	tableView := lipgloss.NewStyle().Width(m.width).Border(lipgloss.RoundedBorder()).Render(m.renderTable())
	if m.paginator.TotalPages <= 1 {
		return tableView
	}
	return lipgloss.JoinVertical(lipgloss.Right,
		tableView,
		lipgloss.NewStyle().Padding(0, 2).Render(m.paginator.View()),
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
	blankCell := strings.Repeat(" ", m.itemWidth)
	for i, cell := range padIterIdx(items, blankCell, nbCells) {
		if i > 0 && i%m.numColumns == 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(cell)
		x := i % m.numColumns
		if x < m.numColumns-1 {
			sb.WriteString(defaultSeparator)
		}
	}
	return sb.String()
}

type iteratorIdx[T any] func(yield func(int, T) bool)

// padIterIdx returns an iterator that yields indexed items from a slice, and then continues
// padding with a given value if the slice is shorter than the number of items requested.
func padIterIdx[T any](items []T, padding T, numItems int) iteratorIdx[T] {
	return func(yield func(int, T) bool) {
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
