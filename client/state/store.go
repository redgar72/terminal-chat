package state

import (
	"time"
)

// Message represents a chat message in the client
type Message struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	AuthorID  string    `json:"author_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Channel represents a chat channel/room
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChatState holds the entire application state
type ChatState struct {
	ActiveChannel string
	Channels      []Channel
	Messages      map[string][]Message // per-channel buffer
	TypingUsers   map[string][]string  // per-channel typing users
	Username      string
	Connected     bool
	InputBuffer   string
}

// NewChatState creates a new chat state
func NewChatState(username string) *ChatState {
	return &ChatState{
		ActiveChannel: "general",
		Channels: []Channel{
			{ID: "general", Name: "general"},
			{ID: "random", Name: "random"},
			{ID: "dev", Name: "dev"},
		},
		Messages:    make(map[string][]Message),
		TypingUsers: make(map[string][]string),
		Username:    username,
		Connected:   false,
	}
}

// AddMessage adds a message to a channel
func (s *ChatState) AddMessage(channelID string, msg Message) {
	if s.Messages[channelID] == nil {
		s.Messages[channelID] = []Message{}
	}
	s.Messages[channelID] = append(s.Messages[channelID], msg)

	// Keep only last 1000 messages per channel to prevent memory issues
	if len(s.Messages[channelID]) > 1000 {
		s.Messages[channelID] = s.Messages[channelID][1:]
	}
}

// GetMessages returns messages for a channel
func (s *ChatState) GetMessages(channelID string) []Message {
	return s.Messages[channelID]
}

// SetActiveChannel changes the active channel
func (s *ChatState) SetActiveChannel(channelID string) {
	s.ActiveChannel = channelID
}

// AddChannel adds a new channel
func (s *ChatState) AddChannel(channel Channel) {
	s.Channels = append(s.Channels, channel)
}

// SetTypingUsers sets the typing users for a channel
func (s *ChatState) SetTypingUsers(channelID string, users []string) {
	s.TypingUsers[channelID] = users
}

// GetTypingUsers returns typing users for a channel
func (s *ChatState) GetTypingUsers(channelID string) []string {
	return s.TypingUsers[channelID]
}
