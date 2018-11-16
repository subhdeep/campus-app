package models

// MessageType identifies the type of a message and how to handle it
type MessageType string

// Username type is for type checking
type Username string

// ConnID type is for type checking
type ConnID string

const (
	Chat       MessageType = "chat"
	ChatAck                = "chat-ack"
	WebRTC                 = "webrtc"
	WebRTCAck              = "webrtc-ack"
	WebRTCInit             = "webrtc-init"
)

// ServerClientMessage is the generic message exchanged between
// client and server.
type ServerClientMessage struct {
	Type    MessageType `json:"type"`
	Message interface{} `json:"message"`
}

// ClientChatMessage is the chat message sent from a client to the
// server.
type ClientChatMessage struct {
	To   Username `json:"to"`
	Body string   `json:"body"`
	TID  int      `json:"tid"`
}

// ClientAckMessage is the acknowledment message sent from the server to the client
type ClientAckMessage struct {
	ChatMessage
	TID int `json:"tid"`
}
