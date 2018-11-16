package models

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

// MarkOnline is used to mark a user as online
func MarkOnline(userID Username) {
	var score = (float64)(time.Now().Unix())
	var member = redis.Z{
		Score:  score,
		Member: string(userID),
	}
	cmd := client.ZAdd("online-users", member)
	if cmd.Err() != nil {
		log.Printf("[warn] An error occurred while interacting with redis: %v", cmd.Err())
	}
}

// IsOnline checks if user is online
func IsOnline(userID Username) (bool, error) {
	var currentscore = (float64)(time.Now().Unix())
	cmd := client.ZScore("online-users", string(userID))
	value, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		log.Printf("[warn] An error occurred while interacting with redis: %v", cmd.Err())
		return false, cmd.Err()
	}
	if currentscore-value < (float64)(5*time.Minute/1000) {
		return true, nil
	}
	return false, nil
}
