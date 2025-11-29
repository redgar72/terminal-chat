package hub

import (
	"encoding/json"
	"log"
	"sync"

	"terminal-chat/server/models"
)

// Hub maintains the set of active clients and broadcasts messages to channels
type Hub struct {
	// Registered clients
	clients map[*models.Client]bool

	// Inbound messages from the clients
	Broadcast chan []byte

	// Register requests from the clients
	Register chan *models.Client

	// Unregister requests from clients
	Unregister chan *models.Client

	// Channel subscriptions: channelID -> clients
	channels map[string]map[*models.Client]bool

	// Mutex for thread safety
	mu sync.RWMutex
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
		clients:    make(map[*models.Client]bool),
		channels:   make(map[string]map[*models.Client]bool),
	}
}

// Run starts the hub and handles client registration, unregistration, and message broadcasting
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			log.Printf("Client %s connected", client.Username)

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)

				// Remove from all channels
				h.mu.Lock()
				for channelID, clients := range h.channels {
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.channels, channelID)
					}
				}
				h.mu.Unlock()

				log.Printf("Client %s disconnected", client.Username)
			}

		case message := <-h.Broadcast:
			var event models.Event
			if err := json.Unmarshal(message, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			h.handleEvent(event, message)
		}
	}
}

// handleEvent processes different types of events
func (h *Hub) handleEvent(event models.Event, rawMessage []byte) {
	switch event.Type {
	case "join_channel":
		h.joinChannel(event.Channel, event.From)
	case "leave_channel":
		h.leaveChannel(event.Channel, event.From)
	case "send_message":
		h.broadcastToChannel(event.Channel, rawMessage)
	case "typing_start", "typing_stop":
		h.broadcastToChannel(event.Channel, rawMessage)
	case "user_joined", "user_left":
		h.broadcastToChannel(event.Channel, rawMessage)
	}
}

// joinChannel adds a client to a channel
func (h *Hub) joinChannel(channelID, username string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.channels[channelID] == nil {
		h.channels[channelID] = make(map[*models.Client]bool)
	}

	// Find client by username (in a real app, you'd have a better lookup)
	for client := range h.clients {
		if client.Username == username {
			h.channels[channelID][client] = true
			log.Printf("User %s joined channel %s", username, channelID)
			break
		}
	}
}

// leaveChannel removes a client from a channel
func (h *Hub) leaveChannel(channelID, username string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, exists := h.channels[channelID]; exists {
		for client := range clients {
			if client.Username == username {
				delete(clients, client)
				log.Printf("User %s left channel %s", username, channelID)
				break
			}
		}

		if len(clients) == 0 {
			delete(h.channels, channelID)
		}
	}
}

// broadcastToChannel sends a message to all clients in a specific channel
func (h *Hub) broadcastToChannel(channelID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, exists := h.channels[channelID]; exists {
		log.Printf("Broadcasting to %d clients in channel %s", len(clients), channelID)
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
				delete(clients, client)
			}
		}
	} else {
		log.Printf("No clients in channel %s", channelID)
	}
}

// GetClientsInChannel returns all clients in a channel (for debugging/admin purposes)
func (h *Hub) GetClientsInChannel(channelID string) []*models.Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []*models.Client
	if channelClients, exists := h.channels[channelID]; exists {
		for client := range channelClients {
			clients = append(clients, client)
		}
	}
	return clients
}
