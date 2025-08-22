package anilist

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"discord-anime-bot/internal/config"
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
	filePath      string
}

// NewNotificationService creates a new notification service
func NewNotificationService(session *discordgo.Session) *NotificationService {
	service := &NotificationService{
		notifications: make(map[string]*notificationTimer),
		session:       session,
		filePath:      config.LoadConfig().StorageFile,
	}

	// Load existing notifications
	service.loadNotifications()

	return service
}

// loadNotifications loads notifications from the JSON file
func (ns *NotificationService) loadNotifications() {
	// Check if file exists
	if _, err := os.Stat(ns.filePath); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(ns.filePath)
	if err != nil {
		log.Printf("Error reading notifications file: %v", err)
		return
	}

	// Handle edge case where file contains "null" instead of empty array
	if string(data) == "null" || string(data) == "" {
		log.Printf("Notifications file is null or empty, initializing with empty array")
		return
	}

	var persistedNotifications []types.PersistedNotification
	if err := json.Unmarshal(data, &persistedNotifications); err != nil {
		log.Printf("Error parsing notifications file: %v", err)
		// Reset file to empty array if parsing fails
		if err := ns.saveNotificationsInternal(); err != nil {
			log.Printf("Failed to reset notifications file: %v", err)
		}
		return
	}

	log.Printf("Successfully parsed %d notifications from file", len(persistedNotifications))

	// Recreate timers for future notifications
	for _, pn := range persistedNotifications {
		// Convert from milliseconds (TypeScript format) to seconds for Go
		airingTime := time.Unix(pn.AiringAt/1000, 0)
		if airingTime.After(time.Now()) {
			notification := types.NotificationEntry{
				AnimeID:   pn.AnimeID,
				ChannelID: pn.ChannelID,
				UserID:    pn.UserID,
				AiringAt:  pn.AiringAt / 1000, // Convert to seconds for storage
				Episode:   pn.Episode,
			}

			ns.scheduleNotification(&notification)
		}
	}

	log.Printf("Loaded %d notifications from storage", len(persistedNotifications))
}

// saveNotificationsInternal saves notifications without acquiring locks (internal use)
func (ns *NotificationService) saveNotificationsInternal() error {
	// Initialize as empty slice to ensure JSON marshals as [] not null
	persistedNotifications := make([]types.PersistedNotification, 0)

	for _, timer := range ns.notifications {
		persistedNotifications = append(persistedNotifications, types.PersistedNotification{
			AnimeID:   timer.Entry.AnimeID,
			ChannelID: timer.Entry.ChannelID,
			UserID:    timer.Entry.UserID,
			AiringAt:  timer.Entry.AiringAt,
			Episode:   timer.Entry.Episode,
		})
	}

	if err := os.MkdirAll(filepath.Dir(ns.filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(persistedNotifications, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notifications: %w", err)
	}

	log.Printf("💾 Saved %d notifications to storage", len(persistedNotifications))

	if err := os.WriteFile(ns.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notifications file: %w", err)
	}

	log.Printf("💾 Saved %d notifications to storage", len(persistedNotifications))
	return nil
}

// createNotificationKey creates a unique key for a notification
func createNotificationKey(animeID int, channelID, userID string) string {
	return fmt.Sprintf("%d-%s-%s", animeID, channelID, userID)
}

// scheduleNotification sets up a timer for the notification (with locking)
func (ns *NotificationService) scheduleNotification(entry *types.NotificationEntry) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.scheduleNotificationInternal(entry)
}

// scheduleNotificationInternal sets up a timer for the notification (without locking)
func (ns *NotificationService) scheduleNotificationInternal(entry *types.NotificationEntry) {
	notificationKey := createNotificationKey(entry.AnimeID, entry.ChannelID, entry.UserID)

	// Calculate delay until airing time
	airingTime := time.Unix(entry.AiringAt, 0)
	delay := time.Until(airingTime)
	if delay <= 0 {
		log.Printf("⚠️ Notification for anime %d is in the past, skipping", entry.AnimeID)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	timer := time.AfterFunc(delay, func() {
		select {
		case <-ctx.Done():
			return
		default:
			ns.sendNotification(entry)
			// Remove the notification after sending (no saving needed since it's fired)
			ns.mu.Lock()
			delete(ns.notifications, notificationKey)
			ns.mu.Unlock()
		}
	})

	ns.notifications[notificationKey] = &notificationTimer{
		Entry:      entry,
		Timer:      timer,
		CancelFunc: cancel,
	}

	log.Printf("⏰ Scheduled notification for anime %d in %v", entry.AnimeID, delay)
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
		Title:       "🔔 Episode Alert!",
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
		log.Printf("📢 Sent notification for anime %s episode %d to user %s", title, entry.Episode, entry.UserID)
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

	log.Printf("🔔 Adding notification for anime %d, episode %d", animeID, episode)

	// Schedule the notification (internal method that doesn't acquire locks)
	ns.scheduleNotificationInternal(entry)

	// Save to file (using internal method to avoid deadlock)
	return ns.saveNotificationsInternal()
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

	// Save to file (using internal method to avoid deadlock)
	return ns.saveNotificationsInternal()
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

// Cleanup stops all timers and saves notifications
func (ns *NotificationService) Cleanup() {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	log.Println("🧹 Cleaning up notification service...")

	for _, timer := range ns.notifications {
		timer.Timer.Stop()
		timer.CancelFunc()
	}

	if err := ns.saveNotificationsInternal(); err != nil {
		log.Printf("Error saving notifications during cleanup: %v", err)
	}

	log.Println("✅ Notification service cleanup complete")
}
