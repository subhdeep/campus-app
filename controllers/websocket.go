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
	userID := ctx.Values().Get("userID").(models.Username)
	models.AddConnection(userID, c)
	c.OnMessage(websocketMessageHandler(userID, logger, c))
	c.OnError(func(err error) {
		logger.Warnf("Error occurred for %s: %v", userID, err)
	})
}

func websocketMessageHandler(userID models.Username, logger *golog.Logger, userCon websocket.Connection) func([]byte) {
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
		case models.WebRTC:
			webRTCHandler(userID, logger, msg.Message, userCon)
		case models.WebRTCAck:
			webRTCAckHandler(userID, logger, msg.Message, userCon)
		case models.WebRTCInit:
			webRTCInitHandler(userID, logger, msg.Message, userCon)
		}
	}
}

func chatHandler(userID models.Username, logger *golog.Logger, msg interface{}, userCon websocket.Connection) {
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
	chatMsg := models.CreateChatMessage(&clientChatMsg, string(userID))
	models.CreateRecentMessage(chatMsg, string(userID), chatMsg.To)
	logger.Infof("Message: %s from %s to %s", chatMsg.Body, userID, chatMsg.To)

	// Sending ack back to sender
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

	// Publishing message to all interested peeps
	models.PublishChatMessage(chatMsg, models.ConnID(userCon.ID()))
}

func webRTCHandler(userID models.Username, logger *golog.Logger, msg interface{}, userCon websocket.Connection) {
	var webRTCBytes []byte
	webRTCBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	var webRTC models.WebRTCMessage
	if err := json.Unmarshal(webRTCBytes, &webRTC); err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	webRTC.FromID = models.ConnID(userCon.ID())
	models.PublishWebRTCMessage(webRTC)
}

func webRTCInitHandler(userID models.Username, logger *golog.Logger, msg interface{}, userCon websocket.Connection) {
	var webRTCInitBytes []byte
	webRTCInitBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	var webRTCInit models.WebRTCInitMessage
	if err := json.Unmarshal(webRTCInitBytes, &webRTCInit); err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	webRTCInit.From = userID
	webRTCInit.FromID = models.ConnID(userCon.ID())
	models.PublishWebRTCInitMessage(webRTCInit)
}

func webRTCAckHandler(userID models.Username, logger *golog.Logger, msg interface{}, userCon websocket.Connection) {
	var webRTCAckBytes []byte
	webRTCAckBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	var webRTCAck models.WebRTCAckMessage
	if err := json.Unmarshal(webRTCAckBytes, &webRTCAck); err != nil {
		logger.Errorf("Invalid message: %v", err)
		return
	}
	webRTCAck.From = userID
	webRTCAck.FromID = models.ConnID(userCon.ID())
	models.PublishWebRTCAckMessage(webRTCAck)
}
