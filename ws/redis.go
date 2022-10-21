package demo

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	log.Println("connecting to Redis...")
	redisClient = redis.NewClient(&redis.Options{Addr: "redis:6379"})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("failed to connect to redis", err)
	}
	log.Println("connected to redis")
}
