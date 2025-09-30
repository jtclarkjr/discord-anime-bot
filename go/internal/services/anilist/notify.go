package anilist

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"discord-anime-bot/internal/services/redis"
	"discord-anime-bot/internal/types"

	"github.com/bwmarrin/discordgo"
)

// GetNextEpisode gets the next airing episode for an anime
func GetNextEpisode(animeID int) (*types.NextAiringEpisode, error) {
	anime, err := GetAnimeByID(animeID)
	if err != nil {
		return nil, err
	}

	return anime.NextAiringEpisode, nil
}

// notificationTimer holds a notification with its timer
type notificationTimer struct {
	Entry      *types.NotificationEntry
	Timer      *time.Timer
	CancelFunc context.CancelFunc
}

// NotificationService handles episode notifications
type NotificationService struct {
	notifications map[string]*notificationTimer
	session       *discordgo.Session
	mu            sync.RWMutex
}

// NewNotificationService creates a new notification service
func NewNotificationService(session *discordgo.Session) *NotificationService {
	service := &NotificationService{
		notifications: make(map[string]*notificationTimer),
		session:       session,
	}

	// Load existing notifications
	service.loadNotifications()

	return service
}

// loadNotifications loads notifications from Redis and reschedules them
func (ns *NotificationService) loadNotifications() {
	ctx := context.Background()

	// Get all notification keys from Redis
	notificationKeys, err := redis.Keys(ctx, "notification:*")
	if err != nil {
		log.Printf("Error fetching notification keys from Redis: %v", err)
		return
	}

	if len(notificationKeys) == 0 {
		log.Println("No notifications found in Redis")
		return
	}

	now := time.Now()
	loadedCount := 0

	for _, redisKey := range notificationKeys {
		var persistedNotification types.PersistedNotification
		err := redis.Get(ctx, redisKey, &persistedNotification)
		if err != nil {
			log.Printf("Error getting notification %s from Redis: %v", redisKey, err)
			continue
		}

		// Skip expired notifications
		airingTime := time.Unix(persistedNotification.AiringAt/1000, 0)
		if airingTime.Before(now) {
			// Remove expired notification
			if err := redis.Delete(ctx, redisKey); err != nil {
				log.Printf("Error deleting expired notification %s from Redis: %v", redisKey, err)
			}
			continue
		}

		notification := types.NotificationEntry{
			AnimeID:   persistedNotification.AnimeID,
			ChannelID: persistedNotification.ChannelID,
			UserID:    persistedNotification.UserID,
			AiringAt:  persistedNotification.AiringAt / 1000, // Convert to seconds
			Episode:   persistedNotification.Episode,
		}

		ns.scheduleNotificationInternal(&notification)
		loadedCount++
	}

	log.Printf("Loaded %d active notifications from Redis", loadedCount)
}

// saveNotificationToRedis saves a single notification to Redis
func (ns *NotificationService) saveNotificationToRedis(notificationKey string, entry *types.NotificationEntry) error {
	ctx := context.Background()
	redisKey := "notification:" + notificationKey

	persistedEntry := types.PersistedNotification{
		AnimeID:   entry.AnimeID,
		ChannelID: entry.ChannelID,
		UserID:    entry.UserID,
		AiringAt:  entry.AiringAt * 1000, // Convert to milliseconds for consistency
		Episode:   entry.Episode,
	}

	// Calculate TTL based on airing time (with buffer)
	airingTime := time.Unix(entry.AiringAt, 0)
	ttl := time.Until(airingTime) + time.Hour // 1 hour buffer
	if ttl <= 0 {
		ttl = time.Minute // Minimum 1 minute TTL
	}

	err := redis.Set(ctx, redisKey, persistedEntry, ttl)
	if err != nil {
		return fmt.Errorf("failed to save notification to Redis: %w", err)
	}

	log.Printf("Saved notification %s to Redis", notificationKey)
	return nil
}

// removeNotificationFromRedis removes a notification from Redis
func (ns *NotificationService) removeNotificationFromRedis(notificationKey string) error {
	ctx := context.Background()
	redisKey := "notification:" + notificationKey

	err := redis.Delete(ctx, redisKey)
	if err != nil {
		return fmt.Errorf("failed to remove notification from Redis: %w", err)
	}

	log.Printf("Removed notification %s from Redis", notificationKey)
	return nil
}

// createNotificationKey creates a unique key for a notification
func createNotificationKey(animeID int, channelID, userID string) string {
	return fmt.Sprintf("%d-%s-%s", animeID, channelID, userID)
}

// scheduleNotificationInternal sets up a timer for the notification (without locking)
func (ns *NotificationService) scheduleNotificationInternal(entry *types.NotificationEntry) {
	notificationKey := createNotificationKey(entry.AnimeID, entry.ChannelID, entry.UserID)

	// Calculate delay until airing time
	airingTime := time.Unix(entry.AiringAt, 0)
	delay := time.Until(airingTime)
	if delay <= 0 {
		log.Printf("Warning: Notification for anime %d is in the past, skipping", entry.AnimeID)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	timer := time.AfterFunc(delay, func() {
		select {
		case <-ctx.Done():
			return
		default:
			ns.sendNotification(entry)
			// Remove the notification after sending
			ns.mu.Lock()
			delete(ns.notifications, notificationKey)
			ns.mu.Unlock()
			// Remove from Redis as well
			if err := ns.removeNotificationFromRedis(notificationKey); err != nil {
				log.Printf("Error removing notification from Redis: %v", err)
			}
		}
	})

	ns.notifications[notificationKey] = &notificationTimer{
		Entry:      entry,
		Timer:      timer,
		CancelFunc: cancel,
	}

	log.Printf("Scheduled notification for anime %d in %v", entry.AnimeID, delay)
}

// sendNotification sends the Discord notification
func (ns *NotificationService) sendNotification(entry *types.NotificationEntry) {
	// Get anime details for the notification
	anime, err := GetAnimeByID(entry.AnimeID)
	if err != nil {
		log.Printf("Error getting anime details for notification: %v", err)
		return
	}

	title := anime.Title.Romaji
	if anime.Title.English != nil && *anime.Title.English != "" {
		title = *anime.Title.English
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Episode Alert!",
		Description: fmt.Sprintf("**Episode %d** of **%s** is now airing!", entry.Episode, title),
		Color:       0x00FF00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: anime.CoverImage.Large,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Set more notifications with /anime notify add",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	content := fmt.Sprintf("<@%s>", entry.UserID)

	_, err = ns.session.ChannelMessageSendComplex(entry.ChannelID, &discordgo.MessageSend{
		Content: content,
		Embeds:  []*discordgo.MessageEmbed{embed},
	})

	if err != nil {
		log.Printf("Error sending notification: %v", err)
	} else {
		log.Printf("Sent notification for anime %s episode %d to user %s", title, entry.Episode, entry.UserID)
	}
}

// AddNotification adds a new episode notification
func (ns *NotificationService) AddNotification(animeID int, channelID, userID string, airingAt time.Time, episode int) error {
	notificationKey := createNotificationKey(animeID, channelID, userID)

	ns.mu.Lock()
	defer ns.mu.Unlock()

	// Check if notification already exists
	if existing, exists := ns.notifications[notificationKey]; exists {
		existing.Timer.Stop()
		existing.CancelFunc()
		delete(ns.notifications, notificationKey)
	}

	entry := &types.NotificationEntry{
		AnimeID:   animeID,
		ChannelID: channelID,
		UserID:    userID,
		AiringAt:  airingAt.Unix(),
		Episode:   episode,
	}

	log.Printf("Adding notification for anime %d, episode %d", animeID, episode)

	// Schedule the notification (internal method that doesn't acquire locks)
	ns.scheduleNotificationInternal(entry)

	// Save to Redis
	return ns.saveNotificationToRedis(notificationKey, entry)
}

// RemoveNotification removes a notification for a specific anime and user
func (ns *NotificationService) RemoveNotification(animeID int, channelID, userID string) error {
	notificationKey := createNotificationKey(animeID, channelID, userID)

	ns.mu.Lock()
	defer ns.mu.Unlock()

	timer, exists := ns.notifications[notificationKey]
	if !exists {
		return fmt.Errorf("notification not found for anime %d and user %s", animeID, userID)
	}

	timer.Timer.Stop()
	timer.CancelFunc()
	delete(ns.notifications, notificationKey)

	// Remove from Redis
	return ns.removeNotificationFromRedis(notificationKey)
}

// GetUserNotifications returns all notifications for a specific user
func (ns *NotificationService) GetUserNotifications(userID string) []*types.NotificationEntry {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	var notifications []*types.NotificationEntry
	for _, timer := range ns.notifications {
		if timer.Entry.UserID == userID {
			notifications = append(notifications, timer.Entry)
		}
	}

	return notifications
}

// Cleanup stops all timers
func (ns *NotificationService) Cleanup() {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	log.Println("Cleaning up notification service...")

	for _, timer := range ns.notifications {
		timer.Timer.Stop()
		timer.CancelFunc()
	}

	log.Println("Notification service cleanup complete")
}
