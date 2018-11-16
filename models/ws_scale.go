package models

import (
	"encoding/json"
	"log"

	"github.com/kataras/iris/websocket"
)

// Connections map of the different client connected to the server
var connections map[string][]websocket.Connection
var WS *websocket.Server

type publishChatPayload struct {
	ChatMessage ChatMessage
	ID          string
}

func init() {
	WS = websocket.New(websocket.Config{})
	connections = make(map[string][]websocket.Connection)
}

func AddConnection(userID string, c websocket.Connection) {
	connections[userID] = append(connections[userID], c)
	c.OnDisconnect(func() {
		c1, ok := connections[userID]
		if !ok || len(c1) == 0 {
			log.Printf("%s is not online. Unable to disconnect", userID)
			return
		}
		for i, con := range c1 {
			if con.ID() == c.ID() {
				c1 = append(c1[:i], c1[i+1:]...)
				break
			}
		}
		connections[userID] = c1
	})
}

func PublishChatMessage(chatMsg ChatMessage, conID string) {
	payload := publishChatPayload{
		ChatMessage: chatMsg,
		ID:          conID,
	}
	marshalled, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
		return
	}
	client.Publish(ChatChannel, marshalled)
}
