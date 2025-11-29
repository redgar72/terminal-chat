package network

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// WSClient handles WebSocket communication with the server
type WSClient struct {
	conn      *websocket.Conn
	serverURL string
	username  string
	incoming  chan []byte
	outgoing  chan []byte
	done      chan struct{}
}

// NewWSClient creates a new WebSocket client
func NewWSClient(serverURL, username string) *WSClient {
	return &WSClient{
		serverURL: serverURL,
		username:  username,
		incoming:  make(chan []byte, 256),
		outgoing:  make(chan []byte, 256),
		done:      make(chan struct{}),
	}
}

// Connect establishes the WebSocket connection
func (c *WSClient) Connect() error {
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return err
	}

	// Add username as query parameter
	q := u.Query()
	q.Set("username", c.username)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn

	// Join the default channel
	joinMsg := map[string]interface{}{
		"type":    "join_channel",
		"channel": "general",
	}
	if data, err := json.Marshal(joinMsg); err == nil {
		c.outgoing <- data
	}

	// Start goroutines
	go c.readPump()
	go c.writePump()

	return nil
}

// Close closes the WebSocket connection
func (c *WSClient) Close() {
	close(c.done)
	if c.conn != nil {
		c.conn.Close()
	}
}

// SendMessage sends a chat message
func (c *WSClient) SendMessage(channel, content string) {
	msg := map[string]interface{}{
		"type":    "send_message",
		"channel": channel,
		"content": content,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	c.outgoing <- data
}

// JoinChannel joins a channel
func (c *WSClient) JoinChannel(channel string) {
	msg := map[string]interface{}{
		"type":    "join_channel",
		"channel": channel,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling join message: %v", err)
		return
	}

	c.outgoing <- data
}

// readPump reads messages from the WebSocket connection
func (c *WSClient) readPump() {
	defer func() {
		close(c.incoming)
	}()

	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			select {
			case c.incoming <- message:
			case <-c.done:
				return
			}
		}
	}
}

// writePump writes messages to the WebSocket connection
func (c *WSClient) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message := <-c.outgoing:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case <-c.done:
			return
		}
	}
}

// GetIncomingChannel returns the channel for incoming messages
func (c *WSClient) GetIncomingChannel() <-chan []byte {
	return c.incoming
}

// GetUsername returns the client username
func (c *WSClient) GetUsername() string {
	return c.username
}
