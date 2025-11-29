package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"terminal-chat/client/state"
)

// MessageView represents the message display area
type MessageView struct {
	chatState *state.ChatState
	viewport  viewport.Model
	width     int
	height    int
}

// NewMessageView creates a new message view
func NewMessageView(chatState *state.ChatState) *MessageView {
	vp := viewport.New(60, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	return &MessageView{
		chatState: chatState,
		viewport:  vp,
	}
}

// SetSize sets the message view dimensions
func (m *MessageView) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
}

// View renders the message view
func (m *MessageView) View() string {
	messages := m.chatState.GetMessages(m.chatState.ActiveChannel)

	var content strings.Builder

	// Add channel header
	content.WriteString(fmt.Sprintf("#%s\n", m.chatState.ActiveChannel))

	// Ensure width is positive before using strings.Repeat
	if m.width > 0 {
		content.WriteString(strings.Repeat("─", m.width) + "\n")
	} else {
		content.WriteString("─────────\n") // fallback
	}

	// Add messages
	for _, msg := range messages {
		timestamp := msg.Timestamp.Format("15:04")
		line := fmt.Sprintf("[%s] %s: %s\n", timestamp, msg.Username, msg.Content)
		content.WriteString(line)
	}

	// Add typing indicator
	typingUsers := m.chatState.GetTypingUsers(m.chatState.ActiveChannel)
	if len(typingUsers) > 0 {
		content.WriteString(fmt.Sprintf("\n%s is typing...", strings.Join(typingUsers, ", ")))
	}

	m.viewport.SetContent(content.String())

	// Auto-scroll to bottom for new messages
	m.viewport.GotoBottom()

	return m.viewport.View()
}

// Update updates the message view content
func (m *MessageView) Update() {
	// Content will be updated in View()
}

// ScrollUp scrolls the viewport up
func (m *MessageView) ScrollUp() {
	m.viewport.LineUp(1)
}

// ScrollDown scrolls the viewport down
func (m *MessageView) ScrollDown() {
	m.viewport.LineDown(1)
}

// AddMessage adds a message to the view (called when new messages arrive)
func (m *MessageView) AddMessage(msg state.Message) {
	// The view will update automatically when View() is called
}
