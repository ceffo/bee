package common

import tea "github.com/charmbracelet/bubbletea"

// MsgBatch is a batch of messages
type MsgBatch struct {
	msgs []tea.Cmd
}

// NewMsgBatch creates a new message batch
func NewMsgBatch(cmds ...tea.Cmd) MsgBatch {
	return MsgBatch{msgs: cmds}
}

// Add adds commands to the batch
func (m *MsgBatch) Add(cmds ...tea.Cmd) *MsgBatch {
	// add non nil commands to the batch
	for _, cmd := range cmds {
		if cmd != nil {
			m.msgs = append(m.msgs, cmd)
		}
	}
	return m
}

// Cmd returns the batch of commands
func (m *MsgBatch) Cmd() tea.Cmd {
	return tea.Batch(m.msgs...)
}
