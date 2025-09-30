package bot

import (
	"fmt"
	"log"
	"time"

	"discord-anime-bot/internal/services/anilist"
	"discord-anime-bot/internal/utils"

	"github.com/bwmarrin/discordgo"
)

// handleNextCommand handles the anime next subcommand
func (b *Bot) handleNextCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondWithError(s, i, "Please provide a valid anime ID.")
		return
	}

	animeID := int(options[0].Value.(float64))
	if animeID <= 0 {
		b.respondWithError(s, i, "Please provide a valid anime ID.")
		return
	}

	// Get anime details
	anime, err := anilist.GetAnimeByID(animeID)
	if err != nil {
		log.Printf("Error getting anime by ID %d: %v", animeID, err)
		b.respondWithError(s, i, fmt.Sprintf("No anime found with ID %d.", animeID))
		return
	}

	title := anime.Title.Romaji
	if anime.Title.English != nil && *anime.Title.English != "" {
		title = *anime.Title.English
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		URL:   anime.SiteURL,
		Color: 0x02A9FF,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: anime.CoverImage.Large,
		},
	}

	// Check if anime is finished or has next episode
	if anime.Status == "FINISHED" || anime.Status == "CANCELLED" {
		embed.Description = "This anime has finished airing."
		episodesStr := "Unknown"
		if anime.Episodes != nil {
			episodesStr = fmt.Sprintf("%d", *anime.Episodes)
		}
		embed.Fields = []*discordgo.MessageEmbedField{
			{Name: "Status", Value: anime.Status, Inline: true},
			{Name: "Format", Value: anime.Format, Inline: true},
			{Name: "Total Episodes", Value: episodesStr, Inline: true},
		}
	} else if anime.NextAiringEpisode != nil {
		airingDate := time.Unix(int64(anime.NextAiringEpisode.AiringAt), 0)
		timeString := utils.FormatCountdown(anime.NextAiringEpisode.TimeUntilAiring)
		formattedAirDate := utils.FormatAirDate(airingDate)

		embed.Description = fmt.Sprintf("Episode %d airs in %s", anime.NextAiringEpisode.Episode, timeString)
		embed.Fields = []*discordgo.MessageEmbedField{
			{Name: "Next Episode", Value: fmt.Sprintf("%d", anime.NextAiringEpisode.Episode), Inline: true},
			{Name: "Air Date", Value: formattedAirDate, Inline: false},
			{Name: "Status", Value: anime.Status, Inline: true},
		}
	} else {
		embed.Description = "No upcoming episodes scheduled."
		embed.Fields = []*discordgo.MessageEmbedField{
			{Name: "Status", Value: anime.Status, Inline: true},
			{Name: "Format", Value: anime.Format, Inline: true},
		}
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}
