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

// handleReleaseCommand handles the anime release subcommand
func (b *Bot) handleReleaseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get currently releasing anime
	releasingAnime, err := anilist.GetReleasingAnime(1, 25)
	if err != nil {
		log.Printf("Error getting releasing anime: %v", err)
		b.respondWithError(s, i, "❌ An error occurred while fetching releasing anime.")
		return
	}

	if len(releasingAnime.Data.Page.Media) == 0 {
		b.respondWithError(s, i, "❌ No releasing anime found.")
		return
	}

	// Create a list of currently releasing anime (limit to 15)
	var animeList []string
	maxItems := 15
	if len(releasingAnime.Data.Page.Media) < maxItems {
		maxItems = len(releasingAnime.Data.Page.Media)
	}

	for _, anime := range releasingAnime.Data.Page.Media[:maxItems] {
		title := anime.Title.Romaji
		if anime.Title.English != nil && *anime.Title.English != "" {
			title = *anime.Title.English
		}

		var nextEpisodeInfo string
		if anime.NextAiringEpisode != nil {
			airingDate := time.Unix(int64(anime.NextAiringEpisode.AiringAt), 0)
			formattedTime := utils.FormatCompactDateTime(airingDate)
			nextEpisodeInfo = fmt.Sprintf(" - Ep %d on %s", anime.NextAiringEpisode.Episode, formattedTime)
		} else {
			nextEpisodeInfo = " - No schedule"
		}

		animeList = append(animeList, fmt.Sprintf("**%s** (ID: %d)%s", title, anime.ID, nextEpisodeInfo))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Currently Releasing Anime",
		Description: strings.Join(animeList, "\n"),
		Color:       0x02A9FF,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Showing %d of %d releasing anime", maxItems, releasingAnime.Data.Page.PageInfo.Total),
		},
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}
