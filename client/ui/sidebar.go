package ui

import (
	"fmt"
	"strings"

	"terminal-chat/client/state"
)

// Sidebar represents the channel sidebar
type Sidebar struct {
	chatState *state.ChatState
	width     int
	height    int
	selected  int
}

// NewSidebar creates a new sidebar
func NewSidebar(chatState *state.ChatState) *Sidebar {
	return &Sidebar{
		chatState: chatState,
		selected:  0,
	}
}

// SetSize sets the sidebar dimensions
func (s *Sidebar) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// View renders the sidebar
func (s *Sidebar) View() string {
	var lines []string

	lines = append(lines, "Channels")

	// Ensure width is positive before using strings.Repeat
	if s.width > 0 {
		lines = append(lines, strings.Repeat("─", s.width))
	} else {
		lines = append(lines, "─────────") // fallback
	}

	for i, channel := range s.chatState.Channels {
		prefix := "  "
		if i == s.selected {
			prefix = "> "
		}

		if channel.ID == s.chatState.ActiveChannel {
			prefix = "* "
		}

		line := fmt.Sprintf("%s#%s", prefix, channel.Name)
		lines = append(lines, line)
	}

	// Fill remaining space
	for len(lines) < s.height {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// Update updates the sidebar state
func (s *Sidebar) Update() {
	// Update selected based on active channel
	for i, channel := range s.chatState.Channels {
		if channel.ID == s.chatState.ActiveChannel {
			s.selected = i
			break
		}
	}
}

// NextChannel selects the next channel
func (s *Sidebar) NextChannel() {
	s.selected = (s.selected + 1) % len(s.chatState.Channels)
	channel := s.chatState.Channels[s.selected]
	s.chatState.SetActiveChannel(channel.ID)
}

// PrevChannel selects the previous channel
func (s *Sidebar) PrevChannel() {
	s.selected--
	if s.selected < 0 {
		s.selected = len(s.chatState.Channels) - 1
	}
	channel := s.chatState.Channels[s.selected]
	s.chatState.SetActiveChannel(channel.ID)
}

// GetSelectedChannel returns the currently selected channel
func (s *Sidebar) GetSelectedChannel() string {
	if s.selected >= 0 && s.selected < len(s.chatState.Channels) {
		return s.chatState.Channels[s.selected].ID
	}
	return ""
}
