package models

import (
	"time"
)

// RecentMessage Model
type RecentMessage struct {
	User1         string    `gorm:"primary_key"`
	User2         string    `gorm:"primary_key"`
	CreatedAt     time.Time `gorm:"primary_key"`
	ChatMessageID string    `gorm:"type:uuid"`
}

// RecentMessagePayload is sent to frontend
type RecentMessagePayload struct {
	UserID       string      `json:"userId"`
	FirstMessage ChatMessage `json:"firstMessage"`
}

// CreateRecentMessage function adds a DB entry of a recent message
func CreateRecentMessage(chatMsg ChatMessage, username1 string, username2 string) {
	var msg RecentMessage
	var res RecentMessage
	if username1 < username2 {
		msg = RecentMessage{
			User1: username1,
			User2: username2,
		}
	} else {
		msg = RecentMessage{
			User1: username2,
			User2: username1,
		}
	}
	db.Where(msg).Assign(RecentMessage{ChatMessageID: chatMsg.ID}).FirstOrCreate(&res)
}

// GetRecents function retrieves the messages from a given timestamp
func GetRecents(username string, offset time.Time, limit int) []ChatMessage {
	var msgs []ChatMessage
	db.Joins("JOIN recent_messages ON recent_messages.chat_message_id = chat_messages.id").Where("(((recent_messages.user1 = ?) OR (recent_messages.user2 = ?)) AND chat_messages.created_at < ?)", username, username, offset).Order("created_at desc").Limit(limit).Find(&msgs)
	return msgs
}
