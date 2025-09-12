package redis

import (
	"context"
	"encoding/json"
	"time"
)

// Set stores a key-value pair with optional TTL
func Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	client := GetClient()
	
	// Serialize value to JSON if it's not a string
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data = string(jsonData)
	}
	
	if ttl > 0 {
		return client.SetEx(ctx, key, data, ttl).Err()
	}
	return client.Set(ctx, key, data, 0).Err()
}

// Get retrieves a value by key and optionally unmarshals JSON
func Get(ctx context.Context, key string, dest any) error {
	client := GetClient()
	
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	
	// If dest is a string pointer, set directly
	if strPtr, ok := dest.(*string); ok {
		*strPtr = result
		return nil
	}
	
	// Otherwise unmarshal JSON
	return json.Unmarshal([]byte(result), dest)
}

// GetString retrieves a string value by key
func GetString(ctx context.Context, key string) (string, error) {
	client := GetClient()
	return client.Get(ctx, key).Result()
}

// Delete removes a key
func Delete(ctx context.Context, key string) error {
	client := GetClient()
	return client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func Exists(ctx context.Context, key string) (bool, error) {
	client := GetClient()
	result, err := client.Exists(ctx, key).Result()
	return result == 1, err
}

// Keys returns all keys matching a pattern
func Keys(ctx context.Context, pattern string) ([]string, error) {
	client := GetClient()
	return client.Keys(ctx, pattern).Result()
}

// SetAdd adds an item to a set
func SetAdd(ctx context.Context, key string, members ...any) error {
	client := GetClient()
	return client.SAdd(ctx, key, members...).Err()
}

// SetRemove removes an item from a set
func SetRemove(ctx context.Context, key string, members ...any) error {
	client := GetClient()
	return client.SRem(ctx, key, members...).Err()
}

// SetMembers returns all members of a set
func SetMembers(ctx context.Context, key string) ([]string, error) {
	client := GetClient()
	return client.SMembers(ctx, key).Result()
}

// SetIsMember checks if a value is in a set
func SetIsMember(ctx context.Context, key string, member any) (bool, error) {
	client := GetClient()
	return client.SIsMember(ctx, key, member).Result()
}

// Expire sets TTL on a key
func Expire(ctx context.Context, key string, ttl time.Duration) error {
	client := GetClient()
	return client.Expire(ctx, key, ttl).Err()
}

// TTL gets the time to live for a key
func TTL(ctx context.Context, key string) (time.Duration, error) {
	client := GetClient()
	return client.TTL(ctx, key).Result()
}