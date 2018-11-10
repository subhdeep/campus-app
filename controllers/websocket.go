package controllers

import (
	"encoding/json"

	"github.com/subhdeep/campus-app/models"

	"github.com/kataras/golog"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/websocket"
)

var ws *websocket.Server

var connections map[string][]websocket.Connection

// Type casts the chat models
// type (
// 	MessageType         string
// 	ServerClientMessage models.ServerClientMessage
// 	ServerChatMessage   models.ServerChatMessage
// 	ClientChatMessage   models.ClientChatMessage
// )

const (
	Chat models.MessageType = "chat"
)

func init() {
	ws = websocket.New(websocket.Config{})
	connections = make(map[string][]websocket.Connection)

	ws.OnConnection(websocketConnectionHandler)

}

// Websocket is the context handler for websocket connections
func Websocket() context.Handler {
	return ws.Handler()
}

func websocketConnectionHandler(c websocket.Connection) {
	ctx := c.Context()
	logger := ctx.Application().Logger()
	userID := ctx.Values().Get("userID").(string)
	connections[userID] = append(connections[userID], c)
	c.OnMessage(websocketMessageHandler(userID, logger))
	c.OnError(func(err error) {
		logger.Warnf("Error occurred for %s: %v", userID, err)
	})
	c.OnDisconnect(func() {
		logger.Infof("Disconnected from %s", userID)
		c1, ok := connections[userID]
		if !ok || len(c1) == 0 {
			logger.Infof("%s is not online. Unable to disconnect", userID)
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

func websocketMessageHandler(userID string, logger *golog.Logger) func([]byte) {
	return func(b []byte) {
		var msg models.ServerClientMessage
		if err := json.Unmarshal(b, &msg); err != nil {
			logger.Errorf("Invalid message: %v", err)
			return
		}
		switch msg.Type {
		case Chat:
			chatHandler(userID, logger, msg.Message)
		}
	}
}

func chatHandler(userID string, logger *golog.Logger, msg []byte) {
	var clientChatMsg models.ClientChatMessage
	if err := json.Unmarshal(msg, &clientChatMsg); err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	// TODO need to save the msg to the database.
	chatMsg := models.CreateChatMessage(&clientChatMsg, userID)
	logger.Infof("Message: %s from %s to %s", chatMsg.Body, userID, chatMsg.To)
	c1, ok := connections[clientChatMsg.To]
	if !ok || len(c1) == 0 {
		logger.Infof("%s is not online. Unable to send message", clientChatMsg.To)
		return
	}
	var serverMsg = models.ServerChatMessage{
		From: userID,
		Body: clientChatMsg.Body,
	}
	marshalled, err := json.Marshal(&serverMsg)
	if err != nil {
		logger.Errorf("Unable to marshal message: %v", err)
		return
	}
	for _, con := range c1 {
		err = con.EmitMessage(marshalled)
		if err != nil {
			logger.Errorf("Unable to send message: %v", err)
		}
	}
}
