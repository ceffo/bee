// colomntable is a package that provides a simple way to print a table with columns
package columntable

import (
	"flag"
	"testing"

	"ceffo.com/bee/pkg/columntable/testdata"
	"ceffo.com/bee/pkg/slices"
	"ceffo.com/bee/pkg/testutils"
	"github.com/charmbracelet/lipgloss"
)

var update = flag.Bool("update", false, "update .golden files")

var (
	test_items = []string{
		"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty",
	}
)

func pad_items(items []string, width int) []string {
	return slices.Map(items, func(s string) string {
		return lipgloss.NewStyle().Width(width).Align(lipgloss.Left).Render(s)
	})
}

func TestModel_renderTable(t *testing.T) {
	type fields struct {
		width      int
		height     int
		items      []string
		itemWidth  int
		transposed bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "empty",
			fields: fields{
				width:     10,
				height:    10,
				items:     []string{},
				itemWidth: 10,
			},
		},
		{
			name: "one item",
			fields: fields{
				width:     25,
				height:    3,
				items:     []string{"one"},
				itemWidth: 10,
			},
		},
		{
			name: "all items",
			fields: fields{
				width:     48,
				height:    9,
				items:     test_items,
				itemWidth: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.width = tt.fields.width
			m.height = tt.fields.height
			m.items = pad_items(tt.fields.items, tt.fields.itemWidth)
			m.itemWidth = tt.fields.itemWidth
			m.updatePaginator()
			got := m.renderTable()
			goldenFile := testdata.Golden("render_table/" + tt.name)
			testutils.SaveOrAssertEqual(t, got, goldenFile, *update)
		})
	}
}
