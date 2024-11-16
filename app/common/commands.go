package common

import tea "github.com/charmbracelet/bubbletea"

// ToCmd converts a message to a command
func ToCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
