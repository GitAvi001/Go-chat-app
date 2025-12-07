// main.go
package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var redisAddr = flag.String("redis", "localhost:6379", "redis address")
var redisChannel = flag.String("channel", "chat_messages", "redis pubsub channel")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for demo; restrict in production
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Get username from query param (e.g., /ws?username=alice)
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	ctx, cancel := contextWithCancel()
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Register client
	hub.register <- client

	// Start pumps
	go client.writePump()
	go client.readPump()
}

func main() {
	flag.Parse()
	hub := NewHub()
	go hub.Run()

	// Setup Redis broker
	redisBroker := NewRedisBroker(hub)
	defer redisBroker.Close()

	// Serve static HTML
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Printf("Starting server on %s (redis: %s channel: %s)", *addr, *redisAddr, *redisChannel)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// small helper functions
func contextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
