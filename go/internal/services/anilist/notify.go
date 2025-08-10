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
		filePath:      "data/notifications.json",
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

	var persistedNotifications []types.PersistedNotification
	if err := json.Unmarshal(data, &persistedNotifications); err != nil {
		log.Printf("Error parsing notifications file: %v", err)
		return
	}

	// Recreate timers for future notifications
	for _, pn := range persistedNotifications {
		airingTime := time.Unix(pn.AiringAt, 0)
		if airingTime.After(time.Now()) {
			notification := types.NotificationEntry{
				AnimeID:   pn.AnimeID,
				ChannelID: pn.ChannelID,
				UserID:    pn.UserID,
				AiringAt:  pn.AiringAt,
				Episode:   pn.Episode,
			}

			ns.scheduleNotification(&notification)
		}
	}

	log.Printf("Loaded %d notifications from storage", len(persistedNotifications))
}

// saveNotifications saves current notifications to the JSON file
func (ns *NotificationService) saveNotifications() error {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	var persistedNotifications []types.PersistedNotification
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

	if err := os.WriteFile(ns.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notifications file: %w", err)
	}

	log.Printf("💾 Saved %d notifications to storage", len(persistedNotifications))
	return nil
}

// createNotificationKey creates a unique key for a notification
func createNotificationKey(animeID int, userID string) string {
	return fmt.Sprintf("%d:%s", animeID, userID)
}

// scheduleNotification sets up a timer for the notification
func (ns *NotificationService) scheduleNotification(entry *types.NotificationEntry) {
	notificationKey := createNotificationKey(entry.AnimeID, entry.UserID)

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
			ns.removeNotificationInternal(notificationKey)
		}
	})

	ns.mu.Lock()
	ns.notifications[notificationKey] = &notificationTimer{
		Entry:      entry,
		Timer:      timer,
		CancelFunc: cancel,
	}
	ns.mu.Unlock()

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

// removeNotificationInternal removes a notification without locking (internal use)
func (ns *NotificationService) removeNotificationInternal(key string) {
	timer, exists := ns.notifications[key]
	if !exists {
		return
	}

	timer.Timer.Stop()
	timer.CancelFunc()
	delete(ns.notifications, key)
}

// AddNotification adds a new episode notification
func (ns *NotificationService) AddNotification(animeID int, channelID, userID string, airingAt time.Time, episode int) error {
	notificationKey := createNotificationKey(animeID, userID)

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

	// Schedule the notification
	ns.scheduleNotification(entry)

	// Save to file
	return ns.saveNotifications()
}

// RemoveNotification removes a notification for a specific anime and user
func (ns *NotificationService) RemoveNotification(animeID int, userID string) error {
	notificationKey := createNotificationKey(animeID, userID)

	ns.mu.Lock()
	defer ns.mu.Unlock()

	timer, exists := ns.notifications[notificationKey]
	if !exists {
		return fmt.Errorf("notification not found for anime %d and user %s", animeID, userID)
	}

	timer.Timer.Stop()
	timer.CancelFunc()
	delete(ns.notifications, notificationKey)

	// Save to file
	return ns.saveNotifications()
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

	if err := ns.saveNotifications(); err != nil {
		log.Printf("Error saving notifications during cleanup: %v", err)
	}

	log.Println("✅ Notification service cleanup complete")
}
