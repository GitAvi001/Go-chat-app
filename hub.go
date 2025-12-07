//Central control for clients

package main

import (
	"log"
	"time"
	"encoding/json"
)

type Hub struct{
	clients		map[*Client] bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte

	//send raw json publishing to redis cache
	redisPublish chan []byte
}


func NewHub() *Hub{
	return &Hub{
		clients:	make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
		redisPublish: make(chan []byte, 256),
	}
}

//Helper functions for current time calculation and json marshaling
func now() time.Time{
	return time.Now()
}

func jsonMarshal(v interface{}) ([]byte, error){
	return json.Marshal(v)
}

func (h *Hub) Run(){
	for{
		select{
		case client :=<-h.register:
			h.clients[client] = true //checks the correct client
			log.Printf("client registered: %s (total: %d)", client.username, len(h.clients))

			//client joining notifying
			joinMsg := Message{
				Type: "join",
				Username : client.username,
				Body : client.username+" joined",
				Timestamp: now(),
			}

			b, _ := jsonMarshal(joinMsg)
			h.redisPublish<-b

		
		case client := <-h.unregister:
			if _, ok :=h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("client unregistered: %s (total: %d)", client.username, len(h.clients))

				leaveMsg := Message{
					Type: "leave",
					Username: client.username,
					Body: client.username+" left the chat",
					Timestamp: now(),
				}

				b, _ := jsonMarshal(leaveMsg)
				h.redisPublish<-b
			}

		case message := <-h.broadcast:
			// send message to all connected clients

			for c:= range h.clients{
				select{
				case c.send<-message:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		}
	}

}