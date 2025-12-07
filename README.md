# Go Chat Application

A real-time chat application built with Go, WebSockets, and Redis for message broadcasting.

## Features

- **Real-time Messaging**: Instant message delivery using WebSockets
- **Multiple Users**: Support for multiple concurrent users
- **User Presence**: Notifications when users join or leave the chat
- **Redis Integration**: Message broadcasting using Redis Pub/Sub
- **Simple Interface**: Clean and intuitive web interface

## Prerequisites

- Go 1.16 or higher
- Redis server (for message broadcasting)
- Modern web browser with WebSocket support

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-chat-app.git
   cd go-chat-app
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start the Redis server (if not already running):
   ```bash
   redis-server
   ```

## Configuration

The application supports the following command-line flags:

- `-addr`: HTTP service address (default: `:8080`)
- `-redis`: Redis server address (default: `localhost:6379`)
- `-channel`: Redis Pub/Sub channel name (default: `chat_messages`)

## Running the Application

1. Start the chat server:
   ```bash
   go run .
   ```

2. Open your web browser and navigate to:
   ```
   http://localhost:8080
   ```

3. Enter a username when prompted and start chatting!

## Project Structure

- `main.go`: Application entry point and WebSocket server setup
- `hub.go`: Central hub for managing WebSocket connections
- `client.go`: Client connection handling and message processing
- `message.go`: Message structure and serialization
- `redis_broker.go`: Redis Pub/Sub integration
- `static/`: Frontend assets (HTML, CSS, JavaScript)

## API Endpoints

- `GET /`: Serves the chat interface
- `GET /ws`: WebSocket endpoint for real-time communication

## Message Format

Messages are exchanged in JSON format with the following structure:

```json
{
  "type": "message",
  "username": "username",
  "body": "Hello, World!",
  "timestamp": "2023-12-07T08:23:00Z"
}
```

## Message Types

- `message`: Regular chat message
- `join`: Notification when a user joins the chat
- `leave`: Notification when a user leaves the chat

## License

This project is open source and available under the [MIT License](LICENSE).