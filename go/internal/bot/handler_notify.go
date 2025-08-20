package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"discord-anime-bot/internal/services/anilist"
	"discord-anime-bot/internal/utils"

	"github.com/bwmarrin/discordgo"
)

// handleNotifyCommand handles the anime notify command
func (b *Bot) handleNotifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	var action string
	var animeID int

	// Parse options
	for _, option := range options {
		switch option.Name {
		case "action":
			action = option.StringValue()
		case "id":
			animeID = int(option.IntValue())
		}
	}

	// If no action is specified, show the list
	if action == "" {
		b.handleNotifyListCommand(s, i)
		return
	}

	// Validate that ID is provided for actions that require it
	if (action == "add" || action == "cancel") && animeID == 0 {
		message := "Please provide an anime ID for this action."
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	switch action {
	case "add":
		b.handleNotifyAddCommand(s, i, animeID)
	case "cancel":
		b.handleNotifyCancelCommand(s, i, animeID)
	default:
		b.respondWithError(s, i, "Unknown notify action")
	}
}

// handleNotifyAddCommand handles adding a notification
func (b *Bot) handleNotifyAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate, animeID int) {
	userID := i.Member.User.ID
	channelID := i.ChannelID

	// Get next episode data
	nextEpisode, err := anilist.GetNextEpisode(animeID)
	if err != nil {
		log.Printf("Error getting next episode for anime %d: %v", animeID, err)
		message := "❌ Failed to get anime information"
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	if nextEpisode == nil {
		message := "❌ No upcoming episodes found for this anime"
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	// Add notification
	err = b.notificationService.AddNotification(animeID, channelID, userID, time.Unix(int64(nextEpisode.AiringAt), 0), nextEpisode.Episode)
	if err != nil {
		log.Printf("Error adding notification: %v", err)
		message := "❌ Failed to add notification"
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	// Get anime details for the response
	anime, err := anilist.GetAnimeByID(animeID)
	if err != nil {
		log.Printf("Error getting anime details: %v", err)
		message := "❌ Failed to get anime information"
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	// Build title
	title := anime.Title.Romaji
	if anime.Title.English != nil && *anime.Title.English != "" {
		title = *anime.Title.English
	}

	airingTime := time.Unix(int64(nextEpisode.AiringAt), 0)
	relativeTime := utils.FormatRelativeTimestamp(airingTime)

	embed := &discordgo.MessageEmbed{
		Title:       "🔔 Notification Added",
		Description: fmt.Sprintf("You'll be notified when **Episode %d** of **%s** airs %s", nextEpisode.Episode, title, relativeTime),
		Color:       0x00FF00,
		Timestamp:   airingTime.Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}

// handleNotifyListCommand handles listing user's notifications
func (b *Bot) handleNotifyListCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	notifications := b.notificationService.GetUserNotifications(userID)

	if len(notifications) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "📋 Your Notifications",
			Description: "You have no active episode notifications.",
			Color:       0x808080,
		}

		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	var description strings.Builder
	for _, notification := range notifications {
		// Get anime details
		anime, err := anilist.GetAnimeByID(notification.AnimeID)
		if err != nil {
			log.Printf("Error getting anime details for %d: %v", notification.AnimeID, err)
			continue
		}

		title := anime.Title.Romaji
		if anime.Title.English != nil && *anime.Title.English != "" {
			title = *anime.Title.English
		}

		airingTime := time.Unix(notification.AiringAt, 0)
		relativeTime := utils.FormatRelativeTimestamp(airingTime)
		description.WriteString(fmt.Sprintf("• **%s** - Episode %d airs %s (ID: %d)\n", title, notification.Episode, relativeTime, notification.AnimeID))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📋 Your Notifications",
		Description: description.String(),
		Color:       0x0099FF,
	}

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}

// handleNotifyCancelCommand handles cancelling a notification
func (b *Bot) handleNotifyCancelCommand(s *discordgo.Session, i *discordgo.InteractionCreate, animeID int) {
	userID := i.Member.User.ID
	channelID := i.ChannelID

	err := b.notificationService.RemoveNotification(animeID, channelID, userID)
	if err != nil {
		log.Printf("Error removing notification: %v", err)
		message := "❌ Failed to cancel notification"
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &message,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🔕 Notification Cancelled",
		Description: fmt.Sprintf("Notification for anime ID %d has been cancelled.", animeID),
		Color:       0xFF6600,
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}
