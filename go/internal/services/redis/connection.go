package redis

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// InitRedis initializes the Redis connection
func InitRedis(redisURL string) error {
	var initErr error

	once.Do(func() {
		// Parse Redis URL
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			initErr = err
			return
		}

		// Create client
		client = redis.NewClient(opt)

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = client.Ping(ctx).Result()
		if err != nil {
			initErr = err
			return
		}

		log.Println("Redis connection established")
	})

	return initErr
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	if client == nil {
		log.Fatal("Redis client not initialized. Call InitRedis() first.")
	}
	return client
}

// IsConnected checks if Redis is connected
func IsConnected() bool {
	if client == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	return err == nil
}

// Close closes the Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
