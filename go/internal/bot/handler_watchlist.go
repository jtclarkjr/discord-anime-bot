package bot

import (
	"fmt"
	"strings"

	"discord-anime-bot/internal/services/anilist"

	"github.com/bwmarrin/discordgo"
)

// handleWatchlistCommand handles the anime watchlist subcommand group
func (b *Bot) handleWatchlistCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondWithError(s, i, "No watchlist subcommand provided")
		return
	}

	subcommand := options[0]
	switch subcommand.Name {
	case "add":
		b.handleWatchlistAddCommand(s, i, subcommand.Options)
	case "list":
		b.handleWatchlistListCommand(s, i)
	case "remove":
		b.handleWatchlistRemoveCommand(s, i, subcommand.Options)
	default:
		b.respondWithError(s, i, "Unknown watchlist subcommand")
	}
}

func (b *Bot) handleWatchlistAddCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		msg := "Please provide an anime ID"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		return
	}
	   animeID := int(options[0].IntValue())
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
		   msg = fmt.Sprintf("✅ Added **%s** (ID: %d) to your watchlist.", title, animeID)
	   } else if err != nil {
		   msg = "❌ Failed to add to watchlist"
	   }
	   s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
}

func (b *Bot) handleWatchlistListCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	ids, err := anilist.GetUserWatchlist(userID)
	if err != nil {
		msg := "❌ Failed to fetch your watchlist"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		return
	}
	if len(ids) == 0 {
		msg := "📋 Your watchlist is empty."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		return
	}
	var sb strings.Builder
	sb.WriteString("Your Anime Watchlist:\n")
	for _, id := range ids {
		sb.WriteString(fmt.Sprintf("• AniList ID: %d\n", id))
	}
	msg := sb.String()
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
}

func (b *Bot) handleWatchlistRemoveCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		msg := "Please provide an anime ID"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		return
	}
	animeID := int(options[0].IntValue())
	userID := i.Member.User.ID
	msg, err := anilist.RemoveFromWatchlist(userID, animeID)
	if err != nil {
		msg = "❌ Failed to remove from watchlist"
	}
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
}
