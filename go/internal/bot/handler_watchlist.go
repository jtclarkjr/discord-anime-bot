package bot

import (
	"fmt"
	"log"
	"strings"

	"discord-anime-bot/internal/services/anilist"

	"github.com/bwmarrin/discordgo"
)

// handleWatchlistCommand handles the anime watchlist command
func (b *Bot) handleWatchlistCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
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
		b.handleWatchlistListCommand(s, i)
		return
	}

	// Validate that ID is provided for actions that require it
	if (action == "add" || action == "remove") && animeID == 0 {
		msg := "Please provide an anime ID for this action."
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}

	switch action {
	case "add":
		b.handleWatchlistAddCommand(s, i, animeID)
	case "remove":
		b.handleWatchlistRemoveCommand(s, i, animeID)
	default:
		b.respondWithError(s, i, "Unknown watchlist action")
	}
}

func (b *Bot) handleWatchlistAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate, animeID int) {
	userID := i.Member.User.ID
	msg, err := anilist.AddToWatchlist(userID, animeID)
	if err == nil && msg == "Anime added to your watchlist." {
		// Fetch anime name for confirmation
		anime, err := anilist.GetAnimeByID(animeID)
		var title string
		if err == nil && anime != nil {
			if anime.Title.English != nil && *anime.Title.English != "" {
				title = *anime.Title.English
			} else {
				title = anime.Title.Romaji
			}
		} else {
			title = fmt.Sprintf("Anime ID %d", animeID)
		}
		msg = fmt.Sprintf("Added **%s** (ID: %d) to your watchlist.", title, animeID)
	} else if err != nil {
		msg = "Failed to add to watchlist"
	}
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}

func (b *Bot) handleWatchlistListCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	ids, err := anilist.GetUserWatchlist(userID)
	if err != nil {
		msg := "Failed to fetch your watchlist"
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}
	if len(ids) == 0 {
		msg := "Your watchlist is empty."
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
		return
	}
	var sb strings.Builder
	sb.WriteString("Your Anime Watchlist:\n")
	for _, id := range ids {
		sb.WriteString(fmt.Sprintf("• AniList ID: %d\n", id))
	}
	msg := sb.String()
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}

func (b *Bot) handleWatchlistRemoveCommand(s *discordgo.Session, i *discordgo.InteractionCreate, animeID int) {
	userID := i.Member.User.ID
	msg, err := anilist.RemoveFromWatchlist(userID, animeID)
	if err == nil && msg == "Anime removed from your watchlist." {
		// Fetch anime name for confirmation
		anime, err := anilist.GetAnimeByID(animeID)
		var title string
		if err == nil && anime != nil {
			if anime.Title.English != nil && *anime.Title.English != "" {
				title = *anime.Title.English
			} else {
				title = anime.Title.Romaji
			}
		} else {
			title = fmt.Sprintf("Anime ID %d", animeID)
		}
		msg = fmt.Sprintf("Removed **%s** (ID: %d) from your watchlist.", title, animeID)
	} else if err != nil {
		msg = "Failed to remove from watchlist"
	}
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg}); err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}
