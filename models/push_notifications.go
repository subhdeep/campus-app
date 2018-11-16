package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	webpush "github.com/sherclockholmes/webpush-go"
	"github.com/subhdeep/campus-app/config"
)

// PushNotification represents a notification subscription
type PushNotification struct {
	gorm.Model
	User  string `gorm:"index"`
	Value string
}

// CreatePushNotification function adds a DB entry of a push notification
func CreatePushNotification(userID string, value string) {
	push := PushNotification{
		User:  userID,
		Value: value,
	}
	db.Create(&push)
}

var pushOptions = webpush.Options{
	Subscriber:      "mailto:example@yourdomain.org",
	VAPIDPrivateKey: config.VAPIDKey,
}

type notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon"`
}

func sendPushNotification(to Username, chatMsg ChatMessage) {
	var subs []PushNotification
	db.Where(PushNotification{User: string(to)}).Find(&subs)
	var push = map[string]interface{}{
		"notification": notification{
			Title: "Campus App",
			Body:  fmt.Sprintf("%s: %s", chatMsg.From, chatMsg.Body),
			Icon:  "http://home.iitk.ac.in/~yashsriv/dp",
		},
	}
	marshalled, err := json.Marshal(push)
	if err != nil {
		log.Printf("[warn] unable to marshal push messsage: %v", err)
		return
	}
	for _, sub := range subs {
		s := webpush.Subscription{}
		if err := json.NewDecoder(bytes.NewBufferString(sub.Value)).Decode(&s); err != nil {
			log.Printf("[warn] unable to decode push subs: %v", err)
			continue
		}

		// Send Notification
		if _, err := webpush.SendNotification(marshalled, &s, &pushOptions); err != nil {
			log.Printf("[warn] unable to push notification: %v", err)
		}
	}
}
