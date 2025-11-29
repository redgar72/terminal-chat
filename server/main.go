package main

import (
	"log"
	"net/http"

	"terminal-chat/server/hub"
	"terminal-chat/server/ws"
)

func main() {
	// Create the chat hub
	h := hub.NewHub()

	// Start the hub in a goroutine
	go h.Run()

	// Set up WebSocket endpoint
	http.HandleFunc("/ws", ws.HandleWebSocket(h))

	// Serve static files (for development)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Chat server starting on :8080")
	log.Println("WebSocket endpoint: ws://localhost:8080/ws?username=yourname")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
