package database

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var RedisCacheClient chan string

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		// Docker container
		// Addr: "redis:6379",
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	log.Println("Redis connected successfullly")
}

func SetupCacheChannel() {
	RedisCacheClient = make(chan string)
	go func(ch chan string) {
		for {
			time.Sleep(10 * time.Second)
			// to receive the string value from goroutine use <-
			key := <-ch
			RedisClient.Del(context.TODO(), key)
			log.Println("Cache cleared " + key)
		}
	}(RedisCacheClient)
}

func ClearCache(keys ...string) {
	for _, key := range keys {
		RedisCacheClient <- key
	}
}
