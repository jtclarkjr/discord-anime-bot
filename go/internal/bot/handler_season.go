package bot

import (
	"fmt"
	"strings"
	"time"

	"discord-anime-bot/internal/services/anilist"
	"discord-anime-bot/internal/types"

	"github.com/bwmarrin/discordgo"
)

// handleSeasonCommand handles the /anime season command
func (b *Bot) handleSeasonCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	// Get season parameter
	if len(options) == 0 {
		b.respondWithError(s, i, "❌ Season parameter is required.")
		return
	}

	season := options[0].StringValue()

	// Validate season
	validSeasons := []string{"winter", "spring", "summer", "fall"}
	isValidSeason := false
	for _, validSeason := range validSeasons {
		if strings.ToLower(season) == validSeason {
			isValidSeason = true
			break
		}
	}

	if !isValidSeason {
		b.respondWithError(s, i, "❌ Invalid season. Please use: winter, spring, summer, or fall.")
		return
	}

	// Get year parameter (default to current year)
	year := time.Now().Year()
	if len(options) > 1 {
		year = int(options[1].IntValue())
	}

	// Fetch seasonal anime
	seasonAnime, err := anilist.GetSeasonAnime(season, year, 1, 50)
	if err != nil {
		b.respondWithError(s, i, "❌ An error occurred while fetching seasonal anime.")
		return
	}

	if len(seasonAnime.Data.Page.Media) == 0 {
		b.respondWithError(s, i, fmt.Sprintf("❌ No anime found for %s %d.", season, year))
		return
	}

	// Create embeds for the seasonal anime
	embeds := b.createSeasonEmbeds(seasonAnime.Data.Page.Media, season, year)

	// Discord allows up to 10 embeds per message
	if len(embeds) <= 10 {
		embedsSlice := embeds
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &embedsSlice,
		})
	} else {
		// Send first 10 embeds initially
		firstEmbeds := embeds[:10]
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &firstEmbeds,
		})
		if err == nil {
			// Send remaining embeds as follow-up messages
			for j := 10; j < len(embeds); j += 10 {
				end := j + 10
				if end > len(embeds) {
					end = len(embeds)
				}
				embedSlice := embeds[j:end]
				_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
					Embeds: embedSlice,
				})
				if err != nil {
					break
				}
			}
		}
	}

	if err != nil {
		b.respondWithError(s, i, "❌ Failed to send seasonal anime information.")
	}
}

// createSeasonEmbeds creates Discord embeds for seasonal anime
func (b *Bot) createSeasonEmbeds(media []types.SeasonAnime, season string, year int) []*discordgo.MessageEmbed {
	const animePerEmbed = 20
	totalEmbeds := (len(media) + animePerEmbed - 1) / animePerEmbed
	embeds := make([]*discordgo.MessageEmbed, 0, totalEmbeds)

	for i := 0; i < totalEmbeds; i++ {
		startIndex := i * animePerEmbed
		endIndex := startIndex + animePerEmbed
		if endIndex > len(media) {
			endIndex = len(media)
		}

		animeSlice := media[startIndex:endIndex]
		var description strings.Builder

		for j, anime := range animeSlice {
			title := anime.Title.Romaji
			if anime.Title.English != nil && *anime.Title.English != "" {
				title = *anime.Title.English
			}

			statusEmoji := getStatusEmoji(anime.Status)
			description.WriteString(fmt.Sprintf("%d. **%s** %s (ID: %d)\n", startIndex+j+1, title, statusEmoji, anime.ID))
		}

		embedTitle := fmt.Sprintf("%s %d Anime", strings.Title(season), year)
		if totalEmbeds > 1 {
			embedTitle += fmt.Sprintf(" (Part %d/%d)", i+1, totalEmbeds)
		}

		embed := &discordgo.MessageEmbed{
			Title:       embedTitle,
			Description: description.String(),
			Color:       0x02A9FF,
		}

		if i == 0 {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Showing all %d anime from %s %d", len(media), season, year),
			}
		}

		embeds = append(embeds, embed)
	}

	return embeds
}

// getStatusEmoji returns an emoji for the anime status
func getStatusEmoji(status string) string {
	switch status {
	case "RELEASING":
		return "🟢"
	case "FINISHED":
		return "✅"
	case "NOT_YET_RELEASED":
		return "🔜"
	case "CANCELLED":
		return "❌"
	case "HIATUS":
		return "⏸️"
	default:
		return "❓"
	}
}
