package models

import (
	"encoding/json"
	"log"

	"github.com/kataras/iris/websocket"
)

// Connections map of the different client connected to the server
var connections map[Username][]websocket.Connection

// WS is the main websocket server maintaining connections between various
// clients
var WS *websocket.Server

type publishChatPayload struct {
	ChatMessage ChatMessage
	ID          ConnID
}

func init() {
	WS = websocket.New(websocket.Config{MaxMessageSize: 4096})
	connections = make(map[Username][]websocket.Connection)
}

func sendToUsername(msg []byte, username Username, ignoring ConnID) {
	if c1, ok := connections[username]; ok {
		for _, con := range c1 {
			if con.ID() != string(ignoring) {
				if err := con.EmitMessage(msg); err != nil {
					log.Printf("[warn] Unable to send message: %v", err)
				}
			}
		}
	}
}

func sendToConnID(msg []byte, connID ConnID) {
	if conn := WS.GetConnection(string(connID)); conn != nil {
		if err := conn.EmitMessage(msg); err != nil {
			log.Printf("[warn] Unable to send message: %v", err)
		}
	}
}

// AddConnection allows adding a connection to our connections map
func AddConnection(userID Username, c websocket.Connection) {
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

// PublishChatMessage allows publishing a chat message to a channel
func PublishChatMessage(chatMsg ChatMessage, connID ConnID) {
	payload := publishChatPayload{
		ChatMessage: chatMsg,
		ID:          connID,
	}
	marshalled, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
		return
	}
	client.Publish(ChatChannel, marshalled)
}

// PublishWebRTCMessage allows publishing a WebRTCMessage to a channel
func PublishWebRTCMessage(msg WebRTCMessage) {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
		return
	}
	client.Publish(WebRTCChannel, marshalled)
}

// PublishWebRTCAckMessage allows publishing a WebRTCAckMessage to a channel
func PublishWebRTCAckMessage(msg WebRTCAckMessage) {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
		return
	}
	client.Publish(WebRTCAckChannel, marshalled)
}

// PublishWebRTCInitMessage allows publishing a WebRTCInitMessage to a channel
func PublishWebRTCInitMessage(msg WebRTCInitMessage) {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
		return
	}
	client.Publish(WebRTCInitChannel, marshalled)
}

// PublishWebRTCRejectMessage allows publishing a WebRTCRejectMessage to a channel
func PublishWebRTCRejectMessage(msg WebRTCRejectMessage) {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
		return
	}
	client.Publish(WebRTCRejectChannel, marshalled)
}

// PublishWebRTCCancelMessage allows publishing a WebRTCCancelMessage to a channel
func PublishWebRTCCancelMessage(msg WebRTCCancelMessage) {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("unexpected error %v", err)
		return
	}
	client.Publish(WebRTCCancelChannel, marshalled)
}

func processChatMessage(chatMessage ChatMessage, connID ConnID) {
	clntSvrMsg := ServerClientMessage{
		Type:    Chat,
		Message: chatMessage,
	}
	marshalled, err := json.Marshal(clntSvrMsg)
	if err != nil {
		log.Printf("Unable to marshal message: %v", err)
		return
	}

	// Sending message to sender's other clients (ignoring connID)
	sendToUsername(marshalled, Username(chatMessage.From), connID)

	// Sending to recipient user's online clients
	if chatMessage.From != chatMessage.To {
		sendToUsername(marshalled, Username(chatMessage.To), "")
		sendPushNotification(Username(chatMessage.To), chatMessage)
	}
}

func processWebRTCInitMessage(msg WebRTCInitMessage) {
	clntSvrMsg := ServerClientMessage{
		Type:    WebRTCInit,
		Message: msg,
	}
	marshalled, err := json.Marshal(clntSvrMsg)
	if err != nil {
		log.Printf("Unable to marshal message: %v", err)
		return
	}

	// Forward this message to all clients of webRTCInit.To
	if msg.From != msg.To {
		sendToUsername(marshalled, msg.To, "")
	}
}

func processWebRTCAckMessage(msg WebRTCAckMessage) {
	clntSvrMsg := ServerClientMessage{
		Type:    WebRTCAck,
		Message: msg,
	}
	marshalled, err := json.Marshal(clntSvrMsg)
	if err != nil {
		log.Printf("Unable to marshal message: %v", err)
		return
	}

	sendToConnID(marshalled, msg.ToID)
	sendToUsername(marshalled, msg.From, msg.FromID)
}

func processWebRTCMessage(msg WebRTCMessage) {
	clntSvrMsg := ServerClientMessage{
		Type:    WebRTC,
		Message: msg,
	}
	marshalled, err := json.Marshal(clntSvrMsg)
	if err != nil {
		log.Printf("Unable to marshal message: %v", err)
		return
	}

	sendToConnID(marshalled, msg.ToID)
}

func processWebRTCRejectMessage(msg WebRTCRejectMessage) {
	clntSvrMsg := ServerClientMessage{
		Type: WebRTCReject,
	}
	marshalled, err := json.Marshal(clntSvrMsg)
	if err != nil {
		log.Printf("Unable to marshal message: %v", err)
		return
	}

	sendToConnID(marshalled, msg.ToID)
	sendToUsername(marshalled, msg.From, msg.FromID)
}

func processWebRTCCancelMessage(msg WebRTCCancelMessage) {
	clntSvrMsg := ServerClientMessage{
		Type:    WebRTCCancel,
		Message: msg,
	}
	marshalled, err := json.Marshal(clntSvrMsg)
	if err != nil {
		log.Printf("Unable to marshal message: %v", err)
		return
	}

	// Forward this message to all clients of webRTCInit.To
	sendToUsername(marshalled, msg.To, "")
}
