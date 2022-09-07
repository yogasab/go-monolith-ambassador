package database

import (
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		// Docker container
		// Addr: "redis:6379",
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	log.Println("Redis connected successfullly")
}
