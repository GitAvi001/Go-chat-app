package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const(

	writeWait = 10 * time.Second //time for write a message by client 


	pongWait = 60 * time.Second //time taken for reply message

	pingPeriod = (pongWait * 9)/10

	//maximum message size 
	maxMessageSize=512

)

type Client struct{

	hub		 *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	ctx		 context.Context
	cancel   context.CancelFunc

}


func (c *Client) readPump(){
	defer func(){
		c.hub.unregister <- c
		c.conn.Close()

	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})


	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("readPump error: %v", err)
			}
			break
		}

	//Parse message to attach timestamp 
	var msg Message

	if err:= json.Unmarshal(message, &msg); err!= nil { //optional error message sending if goes wrong
		continue
	}

	msg.Timestamp = time.Now()
	msg.Username= c.username
	msg.Type = "message"

	b, _ := json.Marshal(msg)

	//Publish to redis
	c.hub.redisPublish <-b

	}
}


func (c *Client) writePump(){
	ticker := time.NewTicker(pingPeriod)

	defer func(){
		ticker.Stop()
		c.conn.Close()
	}()


	for {
		select{
		case message, ok:=<-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			
			if !ok{
				//closing the channel in the hub
				_=c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return 
			}

			//write messages
			w, err :=c.conn.NextWriter(websocket.TextMessage)

			if err!=nil{
				return 
			}

			_, _ = w.Write(message)

			//write queued messages
			n := len(c.send)

			for i:=0; i<n;i++{
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err!=nil{
				return
			}

			case <-ticker.C:

			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {

				return

			}

			case <-c.ctx.Done():

			return
		}
	}
}





