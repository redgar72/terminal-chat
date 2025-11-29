FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o client main.go

FROM golang:1.24-alpine

WORKDIR /app

# Copy source code for development
COPY . .

# Run with go run for development
CMD ["go", "run", "main.go"]
