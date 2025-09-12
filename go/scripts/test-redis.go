package main

import (
	"context"
	"log"
	"os"
	"time"

	"discord-anime-bot/internal/services/redis"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6380"
	}

	log.Println("Testing Redis connection and operations...")

	// Initialize Redis
	if err := redis.InitRedis(redisURL); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	ctx := context.Background()

	// Test 1: Basic key-value operations
	log.Println("1. Testing basic key-value operations...")
	if err := redis.Set(ctx, "test:key", "test-value", 0); err != nil {
		log.Printf("Error setting key: %v", err)
	} else {
		var value string
		if err := redis.Get(ctx, "test:key", &value); err != nil {
			log.Printf("Error getting key: %v", err)
		} else {
			log.Printf("Set and get: %s", value)
		}
	}

	// Test 2: JSON operations
	log.Println("2. Testing JSON operations...")
	testObj := map[string]interface{}{
		"name": "Test Anime",
		"id":   12345,
	}
	if err := redis.Set(ctx, "test:json", testObj, 0); err != nil {
		log.Printf("Error setting JSON: %v", err)
	} else {
		var retrievedObj map[string]interface{}
		if err := redis.Get(ctx, "test:json", &retrievedObj); err != nil {
			log.Printf("Error getting JSON: %v", err)
		} else {
			log.Printf("JSON set and get: %+v", retrievedObj)
		}
	}

	// Test 3: Set operations (for watchlists)
	log.Println("3. Testing set operations (watchlist simulation)...")
	if err := redis.SetAdd(ctx, "test:watchlist:user123", 123, 456, 789); err != nil {
		log.Printf("Error adding to set: %v", err)
	} else {
		members, err := redis.SetMembers(ctx, "test:watchlist:user123")
		if err != nil {
			log.Printf("Error getting set members: %v", err)
		} else {
			log.Printf("Watchlist items: %v", members)
		}

		isInWatchlist, err := redis.SetIsMember(ctx, "test:watchlist:user123", 456)
		if err != nil {
			log.Printf("Error checking set membership: %v", err)
		} else {
			log.Printf("Is 456 in watchlist: %t", isInWatchlist)
		}

		if err := redis.SetRemove(ctx, "test:watchlist:user123", 456); err != nil {
			log.Printf("Error removing from set: %v", err)
		} else {
			updatedMembers, err := redis.SetMembers(ctx, "test:watchlist:user123")
			if err != nil {
				log.Printf("Error getting updated set members: %v", err)
			} else {
				log.Printf("After removal: %v", updatedMembers)
			}
		}
	}

	// Test 4: TTL operations
	log.Println("4. Testing TTL operations...")
	if err := redis.Set(ctx, "test:ttl", "expires-soon", 5*time.Second); err != nil {
		log.Printf("Error setting key with TTL: %v", err)
	} else {
		ttl, err := redis.TTL(ctx, "test:ttl")
		if err != nil {
			log.Printf("Error getting TTL: %v", err)
		} else {
			log.Printf("TTL for test key: %v", ttl)
		}
	}

	// Test 5: Keys pattern matching
	log.Println("5. Testing keys pattern matching...")
	keys, err := redis.Keys(ctx, "test:*")
	if err != nil {
		log.Printf("Error getting keys: %v", err)
	} else {
		log.Printf("All test keys: %v", keys)
	}

	// Cleanup
	log.Println("6. Cleaning up test data...")
	cleanupKeys := []string{"test:key", "test:json", "test:watchlist:user123", "test:ttl"}
	for _, key := range cleanupKeys {
		if err := redis.Delete(ctx, key); err != nil {
			log.Printf("Error deleting key %s: %v", key, err)
		}
	}
	log.Println("Cleanup complete")

	log.Println("All Redis tests passed!")
}