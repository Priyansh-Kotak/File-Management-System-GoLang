package utils

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()	

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // e.g., "localhost:6379"
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Redis connection established")
}
