package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"terminal-chat/client/state"
)

// Input represents the message input field
type Input struct {
	chatState *state.ChatState
	textinput textinput.Model
}

// NewInput creates a new input field
func NewInput(chatState *state.ChatState) *Input {
	ti := textinput.New()
	ti.Placeholder = "Type your message here..."
	ti.CharLimit = 500
	ti.Width = 80

	return &Input{
		chatState: chatState,
		textinput: ti,
	}
}

// SetSize sets the input dimensions
func (i *Input) SetSize(width, height int) {
	i.textinput.Width = width
}

// View renders the input field
func (i *Input) View() string {
	return i.textinput.View()
}

// Update updates the input field
func (i *Input) Update(msg tea.Msg) (interface{}, tea.Cmd) {
	var cmd tea.Cmd
	i.textinput, cmd = i.textinput.Update(msg)
	return i, cmd
}

// Focus focuses the input field
func (i *Input) Focus() {
	i.textinput.Focus()
}

// Blur blurs the input field
func (i *Input) Blur() {
	i.textinput.Blur()
}

// GetValue returns the current input value
func (i *Input) GetValue() string {
	return i.textinput.Value()
}

// Clear clears the input field
func (i *Input) Clear() {
	i.textinput.Reset()
}

// SetValue sets the input value
func (i *Input) SetValue(value string) {
	i.textinput.SetValue(value)
}
