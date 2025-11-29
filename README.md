# Terminal Chat

A Discord-like chat service that runs entirely in the terminal using Bubble Tea and Go.

## Architecture

- **Frontend**: Terminal UI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Backend**: WebSocket server built with Go
- **Communication**: JSON messages over WebSocket

## Features

- Multiple chat channels
- Real-time messaging
- Terminal-based UI
- WebSocket communication

## Project Structure

```
terminal-chat/
├── server/          # WebSocket server
│   ├── hub/         # Message routing hub
│   ├── models/      # Data models
│   ├── ws/          # WebSocket handlers
│   └── main.go      # Server entry point
├── client/          # Terminal UI client
│   ├── ui/          # UI components
│   ├── network/     # WebSocket client
│   ├── state/       # State management
│   └── main.go      # Client entry point
├── docker-compose.yml
└── README.md
```

## Quick Start

### Recommended Setup

1. Start the server using Docker:
```bash
docker-compose up --build server
```

2. In separate terminals, start clients:
```bash
# Terminal 1
cd client && go run main.go user1

# Terminal 2
cd client && go run main.go user2
```

### Manual Setup (Alternative)

1. Start the server:
```bash
cd server
go run main.go
```

2. In another terminal, start the client:
```bash
cd client
go run main.go [username]
```

## Controls

- **`Tab`**: Switch to next channel (cycles through: general → random → dev)
- **`Shift+Tab`**: Switch to previous channel
- **`↑ (Up Arrow)`**: Scroll messages up (view older messages)
- **`↓ (Down Arrow)`**: Scroll messages down (view newer messages)
- **`Enter`**: Send message
- **`Ctrl+C`** or **`q`**: Quit application

## Development

### Prerequisites

- Go 1.21+
- Docker (optional)

### Building

```bash
# Server
cd server && go build -o server main.go

# Client
cd client && go build -o client main.go
```

### Testing

1. Start the server with Docker: `docker-compose up server`
2. Open multiple terminals and run: `cd client && go run main.go [username]`
3. Send messages between clients

## Architecture Details

### Backend Components

- **WebSocket Gateway**: Handles client connections and authentication
- **Chat Hub**: Pub/sub system for message routing between channels
- **Data Models**: Message, Channel, and User structs
- **Event System**: JSON-based messaging protocol

### Frontend Components

- **Layout**: Main UI layout with sidebar, message view, and input
- **Sidebar**: Channel navigation
- **Message View**: Scrollable message display with viewport
- **Input**: Message composition field
- **State Management**: Centralized chat state
- **WebSocket Client**: Server communication

## Future Enhancements

- User authentication
- Message history persistence
- File uploads
- Typing indicators
- User presence
- Channel management
- Message reactions
- Search functionality
