package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"terminal-chat/server/hub"
	"terminal-chat/server/models"
)

// readPump handles incoming messages from the WebSocket connection
func readPump(client *models.Client, h *hub.Hub) {
	defer func() {
		h.Unregister <- client
		client.Conn.(*websocket.Conn).Close()
	}()

	conn := client.Conn.(*websocket.Conn)
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse the message to add sender info
		var event models.Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Add sender information
		event.From = client.Username
		event.Timestamp = time.Now().Unix()

		// Re-encode with sender info
		modifiedMessage, err := json.Marshal(event)
		if err != nil {
			log.Printf("Error marshaling modified message: %v", err)
			continue
		}

		// Send to hub for broadcasting
		h.Broadcast <- modifiedMessage
	}
}

// writePump handles outgoing messages to the WebSocket connection
func writePump(client *models.Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.(*websocket.Conn).Close()
	}()

	conn := client.Conn.(*websocket.Conn)

	for {
		select {
		case message, ok := <-client.Send:
			log.Printf("writePump: sending message to %s: %s", client.Username, string(message))
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error writing message to %s: %v", client.Username, err)
				return
			}
			log.Printf("writePump: message sent successfully to %s", client.Username)

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		// In production, you'd want to restrict this
		return true
	},
}

// HandleWebSocket upgrades HTTP connection to WebSocket and manages the client
func HandleWebSocket(h *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		// For now, we'll use a simple username from query params
		// In a real app, you'd authenticate the user properly
		username := r.URL.Query().Get("username")
		if username == "" {
			username = "Anonymous"
		}

		client := &models.Client{
			ID:       uuid.New().String(),
			UserID:   uuid.New().String(), // In real app, this would come from auth
			Username: username,
			Conn:     conn,
			Send:     make(chan []byte, 256),
		}

		// Register client with hub
		h.Register <- client

		// Start goroutines for reading and writing
		go readPump(client, h)
		go writePump(client)
	}
}
