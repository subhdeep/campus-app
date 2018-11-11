package models

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

func MarkOnline(userID string) {
	var score = (float64)(time.Now().Unix())
	var member = redis.Z{
		Score:  score,
		Member: userID,
	}
	cmd := client.ZAdd("online-users", member)
	if cmd.Err() != nil {
		log.Printf("[warn] An error occurred while interacting with redis: %v", cmd.Err())
	}
}

func IsOnline(userID string) (bool, error) {
	var currentscore = (float64)(time.Now().Unix())
	cmd := client.ZScore("online-users", userID)
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
