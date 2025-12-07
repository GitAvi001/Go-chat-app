package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	client *redis.Client
	pubsub *redis.PubSub
	hub    *Hub
	ctx    context.Context
	stop   chan struct{}
}

func NewRedisBroker(hub *Hub) *RedisBroker {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	pubsub := client.Subscribe(ctx, "chat_messages")

	broker := &RedisBroker{
		client: client,
		pubsub: pubsub,
		hub:    hub,
		ctx:    ctx,
		stop:   make(chan struct{}),
	}

	go broker.listenRedis()

	return broker
}

// Publish outgoing messages from Hub → Redis
func (b *RedisBroker) Publish(message string) {
	err := b.client.Publish(b.ctx, "chat_messages", message).Err()
	if err != nil {
		log.Println("Redis publish error:", err)
	}
}

// Listen incoming messages Redis → Hub
func (b *RedisBroker) listenRedis() {
	ch := b.pubsub.Channel()

	for {
		select {
		case msg := <-ch:
			b.hub.broadcast <- []byte(msg.Payload)
		case <-b.stop:
			return
		}
	}
}

func (b *RedisBroker) Close() {
	close(b.stop)
	b.pubsub.Close()
	b.client.Close()
}
