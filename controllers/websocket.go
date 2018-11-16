package controllers

import (
	"encoding/json"

	"github.com/subhdeep/campus-app/models"

	"github.com/kataras/golog"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/websocket"
)

func init() {
	models.WS.OnConnection(websocketConnectionHandler)
}

// Websocket is the context handler for websocket connections
func Websocket() context.Handler {
	return models.WS.Handler()
}

func websocketConnectionHandler(c websocket.Connection) {
	ctx := c.Context()
	logger := ctx.Application().Logger()
	userID := ctx.Values().Get("userID").(string)
	models.AddConnection(userID, c)
	c.OnMessage(websocketMessageHandler(userID, logger, c))
	c.OnError(func(err error) {
		logger.Warnf("Error occurred for %s: %v", userID, err)
	})
}

func websocketMessageHandler(userID string, logger *golog.Logger, userCon websocket.Connection) func([]byte) {
	return func(b []byte) {
		models.MarkOnline(userID)
		var msg models.ServerClientMessage
		if err := json.Unmarshal(b, &msg); err != nil {
			logger.Errorf("Invalid message: %v", err)
			return
		}
		switch msg.Type {
		case models.Chat:
			chatHandler(userID, logger, msg.Message, userCon)
		}
	}
}

func chatHandler(userID string, logger *golog.Logger, msg interface{}, userCon websocket.Connection) {
	// Get Client Chat Message from interface
	var clientChatMsgBytes []byte
	clientChatMsgBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	var clientChatMsg models.ClientChatMessage
	if err := json.Unmarshal(clientChatMsgBytes, &clientChatMsg); err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}

	// Save message to the DB
	chatMsg := models.CreateChatMessage(&clientChatMsg, userID)
	models.CreateRecentMessage(chatMsg, userID, chatMsg.To)
	logger.Infof("Message: %s from %s to %s", chatMsg.Body, userID, chatMsg.To)

	// Sending to ack back to sender
	var clientAckMsg = models.ClientAckMessage{
		ChatMessage: chatMsg,
		TID:         clientChatMsg.TID,
	}
	var clntSvrMsg = models.ServerClientMessage{
		Type:    models.ChatAck,
		Message: clientAckMsg,
	}
	marshalled, err := json.Marshal(clntSvrMsg)
	if err != nil {
		logger.Errorf("Unable to marshal message: %v", err)
	}
	err = userCon.EmitMessage(marshalled)
	if err != nil {
		logger.Errorf("Unable to send the message: %v", err)
	}

	models.PublishChatMessage(chatMsg, userCon.ID())
}
