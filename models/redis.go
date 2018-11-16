package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/subhdeep/campus-app/config"
)

var client *redis.Client

// Channel constants
const (
	ChatChannel       string = "chat-channel"
	WebRTCChannel            = "webrtc-channel"
	WebRTCAckChannel         = "webrtc-ack-channel"
	WebRTCInitChannel        = "webrtc-init-channel"
)

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

	createChannel(ChatChannel, processChatChannel)
	createChannel(WebRTCChannel, processWebRTCChannel)
	createChannel(WebRTCAckChannel, processWebRTCAckChannel)
	createChannel(WebRTCInitChannel, processWebRTCInitChannel)
}

func createChannel(name string, f func(<-chan *redis.Message)) {
	pubsub := client.Subscribe(name)
	// Wait for confirmation that subscription is created before publishing anything.
	if _, err := pubsub.Receive(); err != nil {
		log.Fatalf("Error while connecting to redis: %v", err)
	}
	// Go channel which receives messages.
	ch := pubsub.Channel()
	go f(ch)
}

func processChatChannel(ch <-chan *redis.Message) {
	for msg := range ch {
		var payload publishChatPayload
		err := json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			log.Fatalf("Invalid Message %v", err)
			continue
		}
		processChatMessage(payload.ChatMessage, payload.ID)
	}
}

func processWebRTCChannel(ch <-chan *redis.Message) {
	for msg := range ch {
		var payload WebRTCMessage
		err := json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			log.Fatalf("Invalid Message %v", err)
			continue
		}
		processWebRTCMessage(payload)
	}
}

func processWebRTCInitChannel(ch <-chan *redis.Message) {
	for msg := range ch {
		var payload WebRTCInitMessage
		err := json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			log.Fatalf("Invalid Message %v", err)
			continue
		}
		processWebRTCInitMessage(payload)
	}
}

func processWebRTCAckChannel(ch <-chan *redis.Message) {
	for msg := range ch {
		var payload WebRTCAckMessage
		err := json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			log.Fatalf("Invalid Message %v", err)
			continue
		}
		processWebRTCAckMessage(payload)
	}
}
