package models

import (
	"time"
)

// ChatMessage Model
type ChatMessage struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	From      string    `json:"from" gorm:"index:msg"`
	To        string    `json:"to" gorm:"index:msg"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at" gorm:"index:msg"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateChatMessage function adds a DB entry of a chat message
func CreateChatMessage(chatMsg *ClientChatMessage, userID string) ChatMessage {
	msg := ChatMessage{
		From: userID,
		To:   string(chatMsg.To),
		Body: chatMsg.Body,
	}
	db.Create(&msg)
	return msg
}

// GetMessages function retrieves the messages from a given timestamp
func GetMessages(username1 string, username2 string, offset time.Time, limit int) []ChatMessage {
	var msgs []ChatMessage
	db.Where("((chat_messages.to = ? AND chat_messages.from = ?) OR (chat_messages.to = ? AND chat_messages.from = ?)) AND created_at < ?", username1, username2, username2, username1, offset).Order("created_at desc").Limit(limit).Find(&msgs)
	return msgs
}
