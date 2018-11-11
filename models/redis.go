package models

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/subhdeep/campus-app/config"
)

var client *redis.Client

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
}
