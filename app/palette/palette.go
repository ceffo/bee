package palette

import "github.com/charmbracelet/lipgloss"

// Palette colors
var (
	Prompt    = lipgloss.NewStyle().Foreground(lipgloss.Color("#931e93"))
	Details   = lipgloss.NewStyle().Foreground(lipgloss.Color("#934e13"))
	Positive  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	Error     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	Primary   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ede2c4"))
	Secondary = lipgloss.NewStyle().Foreground(lipgloss.Color("#c8970e"))
	Tertiary  = lipgloss.NewStyle().Foreground(lipgloss.Color("#1665dd"))
	Help      = lipgloss.NewStyle().Foreground(lipgloss.Color("#343434"))
)
