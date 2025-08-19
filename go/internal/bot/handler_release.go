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
	// Default values
	page := 1
	perPage := 15

	// Parse options for page and perpage
	if len(i.ApplicationCommandData().Options) > 0 {
		for _, opt := range i.ApplicationCommandData().Options[0].Options {
			switch opt.Name {
			case "page":
				if v, ok := opt.Value.(float64); ok {
					page = int(v)
				}
			case "perpage":
				if v, ok := opt.Value.(float64); ok {
					perPage = int(v)
				}
			}
		}
	}

	releasingAnime, err := anilist.GetReleasingAnime(page, perPage)
	if err != nil {
		log.Printf("Error getting releasing anime: %v", err)
		b.respondWithError(s, i, "❌ An error occurred while fetching releasing anime.")
		return
	}

	if len(releasingAnime.Data.Page.Media) == 0 {
		b.respondWithError(s, i, "❌ No releasing anime found.")
		return
	}

	var animeList []string
	for _, anime := range releasingAnime.Data.Page.Media {
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

	pageInfo := releasingAnime.Data.Page.PageInfo
	embed := &discordgo.MessageEmbed{
		Title:       "Currently Releasing Anime",
		Description: strings.Join(animeList, "\n"),
		Color:       0x02A9FF,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d • Showing %d of %d releasing anime", pageInfo.CurrentPage, pageInfo.LastPage, len(releasingAnime.Data.Page.Media), pageInfo.Total),
		},
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}
