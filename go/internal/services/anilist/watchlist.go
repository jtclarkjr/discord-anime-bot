package anilist

import (
	"context"
	"strconv"
	"time"

	"discord-anime-bot/internal/services/redis"
)

const watchlistKeyPrefix = "watchlist:user:"


// AddToWatchlist adds an anime to a user's watchlist
func AddToWatchlist(userID string, animeID int) (string, error) {
	ctx := context.Background()
	redisKey := watchlistKeyPrefix + userID

	// Check if anime is already in watchlist
	isInWatchlist, err := redis.SetIsMember(ctx, redisKey, animeID)
	if err != nil {
		return "", err
	}
	if isInWatchlist {
		return "Anime already in your watchlist.", nil
	}

	// Add anime to watchlist
	if err := redis.SetAdd(ctx, redisKey, animeID); err != nil {
		return "", err
	}

	// Set TTL (30 days)
	if err := redis.Expire(ctx, redisKey, 30*24*time.Hour); err != nil {
		return "", err
	}

	return "Anime added to your watchlist.", nil
}

// RemoveFromWatchlist removes an anime from a user's watchlist
func RemoveFromWatchlist(userID string, animeID int) (string, error) {
	ctx := context.Background()
	redisKey := watchlistKeyPrefix + userID

	// Check if anime is in watchlist
	isInWatchlist, err := redis.SetIsMember(ctx, redisKey, animeID)
	if err != nil {
		return "", err
	}
	if !isInWatchlist {
		return "Anime not found in your watchlist.", nil
	}

	// Remove anime from watchlist
	if err := redis.SetRemove(ctx, redisKey, animeID); err != nil {
		return "", err
	}

	return "Anime removed from your watchlist.", nil
}

// GetUserWatchlist returns a user's watchlist
func GetUserWatchlist(userID string) ([]int, error) {
	ctx := context.Background()
	redisKey := watchlistKeyPrefix + userID

	// Get all members of the set
	members, err := redis.SetMembers(ctx, redisKey)
	if err != nil {
		return []int{}, nil // Return empty slice on error
	}

	// Convert string members to integers
	var animeIDs []int
	for _, member := range members {
		if animeID, err := strconv.Atoi(member); err == nil {
			animeIDs = append(animeIDs, animeID)
		}
	}

	return animeIDs, nil
}
