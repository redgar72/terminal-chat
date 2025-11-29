package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"terminal-chat/client/state"
)

var (
	// Styles
	appStyle = lipgloss.NewStyle().
			Margin(1, 2)

	sidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderRight(true).
			Width(20).
			Height(20)

	messageViewStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				Width(60).
				Height(20)

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			Width(80)
)

// Layout represents the main TUI layout
type Layout struct {
	sidebar     *Sidebar
	messageView *MessageView
	input       *Input
	width       int
	height      int
}

// NewLayout creates a new layout
func NewLayout(chatState *state.ChatState) *Layout {
	return &Layout{
		sidebar:     NewSidebar(chatState),
		messageView: NewMessageView(chatState),
		input:       NewInput(chatState),
	}
}

// SetSize sets the layout dimensions
func (l *Layout) SetSize(width, height int) {
	l.width = width
	l.height = height

	// Update component sizes
	l.sidebar.SetSize(20, height-4)
	l.messageView.SetSize(width-25, height-4)
	l.input.SetSize(width-4, 3)
}

// View renders the layout
func (l *Layout) View() string {
	sidebarView := l.sidebar.View()
	messageView := l.messageView.View()
	inputView := l.input.View()

	// Create the top row (sidebar + message view)
	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarStyle.Render(sidebarView),
		messageViewStyle.Render(messageView),
	)

	// Combine with input at bottom
	layout := lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		inputStyle.Render(inputView),
	)

	return appStyle.Render(layout)
}

// UpdateInput updates the input component
func (l *Layout) UpdateInput(msg tea.Msg) tea.Cmd {
	_, cmd := l.input.Update(msg)
	return cmd
}

// FocusInput focuses the input field
func (l *Layout) FocusInput() {
	l.input.Focus()
}

// BlurInput blurs the input field
func (l *Layout) BlurInput() {
	l.input.Blur()
}

// GetInputValue returns the current input value
func (l *Layout) GetInputValue() string {
	return l.input.GetValue()
}

// ClearInput clears the input field
func (l *Layout) ClearInput() {
	l.input.Clear()
}

// GetSelectedChannel returns the currently selected channel
func (l *Layout) GetSelectedChannel() string {
	return l.sidebar.GetSelectedChannel()
}

// UpdateSidebar updates the sidebar
func (l *Layout) UpdateSidebar() {
	l.sidebar.Update()
}

// UpdateMessageView updates the message view
func (l *Layout) UpdateMessageView() {
	l.messageView.Update()
}

// ScrollMessageViewUp scrolls the message view up
func (l *Layout) ScrollMessageViewUp() {
	l.messageView.ScrollUp()
}

// ScrollMessageViewDown scrolls the message view down
func (l *Layout) ScrollMessageViewDown() {
	l.messageView.ScrollDown()
}

// NextChannel switches to the next channel
func (l *Layout) NextChannel() {
	l.sidebar.NextChannel()
	l.UpdateSidebar()
}

// PrevChannel switches to the previous channel
func (l *Layout) PrevChannel() {
	l.sidebar.PrevChannel()
	l.UpdateSidebar()
}
