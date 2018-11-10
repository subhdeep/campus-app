package models

import (
	"time"
)

type MessageType string

// ServerClientMessage is the generic message exchanged between
// client and server.
type ServerClientMessage struct {
	Type    MessageType `json:"type"`
	Message []byte      `json:"message"`
}

// ClientChatMessage is the chat message sent from a client to the
// server.
type ClientChatMessage struct {
	To   string `json:"to"`
	Body string `json:"body"`
	TID  string `json:"tid"`
}

// ClientAckMessage is the acknowledment messaage sent from the server to the client
type ClientAckMessage struct {
	ID   string `json:"id"`
	To   string `json:"to"`
	Body string `json:"body"`
	TID  string `json:"tid"`
}

// ServerChatMessage is the chat message sent from the server to the client
type ServerChatMessage struct {
	From string `json:"from"`
	Body string `json:"body"`
	ID   string `json:"id"`
}

// ChatMessage Model
type ChatMessage struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateChatMessage function adds a DB entry of a chat message
func CreateChatMessage(chatMsg *ClientChatMessage, userID string) *ChatMessage {
	msg := ChatMessage{
		From: userID,
		To:   chatMsg.To,
		Body: chatMsg.Body,
	}
	db.Create(&msg)
	return &msg
}

// GetMessages function retrieves the messages from a given timestamp
func GetMessages(username string, offset time.Time, limit int) []ChatMessage {
	var msgs []ChatMessage
	db.Where(ChatMessage{To: username}).Or(ChatMessage{From: username}).Where("created_at < ?", offset).Limit(limit).Find(&msgs)
	return msgs
}
