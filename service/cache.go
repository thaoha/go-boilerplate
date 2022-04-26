package main

import (
	"os"
	"log"
	"time"
	"context"
	"strings"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct{
	prefix string
	client *redis.Client
}

func (c *RedisCache) Init() *RedisCache {
	addr := strings.Split(os.Getenv("REDIS_HOST"), ":")
	if len(addr) == 1 {
		redisPort := os.Getenv("REDIS_PORT")
		if redisPort == "" {
			redisPort = "6379"
		}
		addr = append(addr, redisPort)
	}
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     strings.Join(addr, ":"),
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) {
	key = c.prefix + key
	err := c.client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *RedisCache) Get(key string) *redis.StringCmd {
	key = c.prefix + key
	value := c.client.Get(context.Background(), key)
	err := value.Err()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		log.Fatal(err)
	}
	return value
}