package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/subhdeep/campus-app/config"
)

var client *redis.Client

const ChatChannel string = "chat-channel"

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Error while connecting to redis: %v", err)
	}
	pubsub := client.Subscribe(ChatChannel)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive()
	if err != nil {
		log.Fatalf("Error while connecting to redis: %v", err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	go processChatChannel(ch)
}

func processChatChannel(ch <-chan *redis.Message) {
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		var payload publishChatPayload
		err := json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			log.Fatalf("Invalid Message %v", err)
			continue
		}
		processChatMessage(payload.ChatMessage, payload.ID)
	}
}
