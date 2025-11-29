package models

import "time"

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Channel represents a chat channel/room
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// User represents a chat user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID       string
	UserID   string
	Username string
	Conn     interface{} // WebSocket connection
	Send     chan []byte
}

// Event represents messages exchanged between client and server
type Event struct {
	Type      string                 `json:"type"`
	Channel   string                 `json:"channel,omitempty"`
	From      string                 `json:"from,omitempty"`
	Timestamp int64                  `json:"timestamp,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}
