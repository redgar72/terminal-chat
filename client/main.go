package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"terminal-chat/client/network"
	"terminal-chat/client/state"
	"terminal-chat/client/ui"
)

// Model represents the main application model
type Model struct {
	chatState *state.ChatState
	layout    *ui.Layout
	wsClient  *network.WSClient
	ready     bool
	msgs      chan []byte
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Focus the input field so it can receive keyboard input
	m.layout.FocusInput()

	return tea.Batch(
		tea.EnterAltScreen,
		connectToServer(m.wsClient),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "\r", "\n", "return", "ctrl+m":
			// Send message
			content := m.layout.GetInputValue()
			if content != "" && m.chatState.Connected {
				// Add message to local chat state immediately for instant UI feedback
				localMsg := state.Message{
					ChannelID: m.chatState.ActiveChannel,
					Username:  m.chatState.Username,
					Content:   content,
					Timestamp: time.Now(),
				}
				m.chatState.AddMessage(m.chatState.ActiveChannel, localMsg)
				m.layout.UpdateMessageView()

				// Send to server
				m.wsClient.SendMessage(m.chatState.ActiveChannel, content)
				m.layout.ClearInput()
			}
		case "tab":
			// Switch to next channel
			m.layout.NextChannel()
			channel := m.chatState.ActiveChannel
			m.wsClient.JoinChannel(channel)
		case "shift+tab":
			// Switch to previous channel
			m.layout.PrevChannel()
			channel := m.chatState.ActiveChannel
			m.wsClient.JoinChannel(channel)
		case "up":
			m.layout.ScrollMessageViewUp()
		case "down":
			m.layout.ScrollMessageViewDown()
		default:
			// Update input field
			return m, m.layout.UpdateInput(msg)
		}

	case []byte:
		// Handle incoming WebSocket message
		var event map[string]interface{}
		if err := json.Unmarshal(msg, &event); err != nil {
			log.Printf("Error parsing incoming message: %v", err)
			break
		}
		m.handleEvent(event)
		// Restart message polling and trigger re-render
		return m, tea.Batch(
			listenForMessages(m.wsClient),
			tea.Tick(time.Millisecond*10, func(time.Time) tea.Msg {
				return nil // This should trigger another update
			}),
		)

	case tea.WindowSizeMsg:
		if !m.ready {
			m.layout.SetSize(msg.Width, msg.Height)
			m.ready = true
		}

	case connectMsg:
		if msg.err != nil {
			log.Printf("Failed to connect: %v", msg.err)
		} else {
			m.chatState.Connected = true
			cmds = append(cmds, listenForMessages(m.wsClient))
		}

	case checkForMessagesMsg:
		// Trigger another check for messages
		return m, listenForMessages(msg.client)

	case errMsg:
		log.Printf("Error: %v", msg.err)
	}

	return m, tea.Batch(cmds...)
}

// handleEvent processes incoming events from the server
func (m *Model) handleEvent(event map[string]interface{}) {
	eventType, ok := event["type"].(string)
	if !ok {
		return
	}

	switch eventType {
	case "send_message":
		channel, _ := event["channel"].(string)
		from, _ := event["from"].(string)
		content, _ := event["content"].(string)
		timestamp, _ := event["timestamp"].(float64)

		// Skip messages from ourselves to avoid duplicates (we add them locally when sending)
		if from == m.chatState.Username {
			return
		}

		msg := state.Message{
			ChannelID: channel,
			Username:  from,
			Content:   content,
		}

		if timestamp > 0 {
			// Convert Unix timestamp to time.Time
			msg.Timestamp = time.Unix(int64(timestamp), 0)
		}

		m.chatState.AddMessage(channel, msg)
		m.layout.UpdateMessageView()
	}
}

// View renders the UI
func (m Model) View() string {
	return m.layout.View()
}

// Messages for handling async operations
type connectMsg struct {
	err error
}

type errMsg struct {
	err error
}

// Commands
func connectToServer(client *network.WSClient) tea.Cmd {
	return func() tea.Msg {
		err := client.Connect()
		return connectMsg{err: err}
	}
}

// listenForMessages creates a command that checks for WebSocket messages
func listenForMessages(client *network.WSClient) tea.Cmd {
	return tea.Tick(time.Millisecond*10, func(t time.Time) tea.Msg {
		select {
		case data := <-client.GetIncomingChannel():
			return data
		default:
			// No data available, trigger another check
			return checkForMessagesMsg{client: client}
		}
	})
}

// Message to trigger another message check
type checkForMessagesMsg struct {
	client *network.WSClient
}

func main() {
	username := "user"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	chatState := state.NewChatState(username)
	layout := ui.NewLayout(chatState)
	wsClient := network.NewWSClient("ws://localhost:8080/ws", username)

	model := Model{
		chatState: chatState,
		layout:    layout,
		wsClient:  wsClient,
		ready:     false,
		msgs:      make(chan []byte, 100),
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
