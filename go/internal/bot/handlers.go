package bot

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"discord-anime-bot/internal/services/anilist"
	"discord-anime-bot/internal/utils"

	"github.com/bwmarrin/discordgo"
)

// interactionCreate handles slash command interactions
func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "anime" {
		return
	}

	// Defer the response to give us more time to process
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Failed to defer interaction response: %v", err)
		return
	}

	// Get the subcommand
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		b.respondWithError(s, i, "No subcommand provided")
		return
	}

	subcommand := options[0]
	switch subcommand.Name {
	case "find":
		if !b.config.IsOpenAIEnabled {
			b.respondWithError(s, i, "❌ The find command is disabled because OpenAI API key is not configured.")
			return
		}
		b.handleFindCommand(s, i, subcommand.Options)
	case "search":
		b.handleSearchCommand(s, i, subcommand.Options)
	case "release":
		b.handleReleaseCommand(s, i, subcommand.Options)
	case "next":
		b.handleNextCommand(s, i, subcommand.Options)
	default:
		b.respondWithError(s, i, "Unknown subcommand")
	}
}

// handleFindCommand handles the anime find subcommand
func (b *Bot) handleFindCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if !b.config.IsOpenAIEnabled {
		b.respondWithError(s, i, "❌ The find command is disabled because OpenAI API key is not configured. Please set the OPENAI_API_KEY environment variable to use AI-powered anime search.")
		return
	}

	if len(options) == 0 {
		b.respondWithError(s, i, "Please provide a description to search for anime.")
		return
	}

	prompt := options[0].StringValue()
	if prompt == "" {
		b.respondWithError(s, i, "Please provide a description to search for anime.")
		return
	}

	// Find anime using AI
	matches, err := anilist.FindAnimeWithDetails(prompt, b.config.OpenAIAPIKey)
	if err != nil {
		log.Printf("Error finding anime: %v", err)
		b.respondWithError(s, i, "❌ An error occurred while searching for anime.")
		return
	}

	if len(matches) == 0 {
		b.respondWithError(s, i, fmt.Sprintf("❌ No anime found matching the description: \"%s\"", prompt))
		return
	}

	// Create embed for the best match
	bestMatch := matches[0]
	anime := bestMatch.Anime

	// Build the title, preferring English over Romaji
	title := anime.Title.Romaji
	if anime.Title.English != nil && *anime.Title.English != "" {
		title = *anime.Title.English
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🎯 %s", title),
		URL:         anime.SiteURL,
		Description: bestMatch.Reason,
		Color:       0x00FF00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: anime.CoverImage.Large,
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Romaji Title", Value: anime.Title.Romaji, Inline: true},
			{Name: "Native Title", Value: anime.Title.Native, Inline: true},
			{Name: "Format", Value: anime.Format, Inline: true},
			{Name: "Status", Value: anime.Status, Inline: true},
			{Name: "AniList ID", Value: fmt.Sprintf("%d", anime.ID), Inline: true},
			{Name: "AI Confidence", Value: fmt.Sprintf("%d%%", int(math.Round(bestMatch.Confidence*100))), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Powered by GPT-5 + AniList",
		},
	}

	responseText := fmt.Sprintf("🤖 **AI Found Anime Based on:** \"%s\"\n\n", prompt)

	if len(matches) > 1 {
		responseText += "**Other possible matches:**\n"
		var otherMatches []string
		for i, match := range matches[1:] {
			if i >= 2 { // Only show top 3 total (best + 2 others)
				break
			}
			matchTitle := match.Anime.Title.Romaji
			if match.Anime.Title.English != nil && *match.Anime.Title.English != "" {
				matchTitle = *match.Anime.Title.English
			}
			otherMatches = append(otherMatches, fmt.Sprintf("%d. **%s** (%d%% match)",
				i+2, matchTitle, int(math.Round(match.Confidence*100))))
		}
		responseText += strings.Join(otherMatches, "\n")
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &responseText,
		Embeds:  &[]*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}

// handleSearchCommand handles the anime search subcommand
func (b *Bot) handleSearchCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if len(options) == 0 {
		b.respondWithError(s, i, "Please provide a search query.")
		return
	}

	query := options[0].StringValue()
	if query == "" {
		b.respondWithError(s, i, "Please provide a search query.")
		return
	}

	// Search anime using AniList
	searchResults, err := anilist.SearchAnime(query, 1, 5)
	if err != nil {
		log.Printf("Error searching anime: %v", err)
		b.respondWithError(s, i, "❌ An error occurred while searching for anime.")
		return
	}

	if len(searchResults.Data.Page.Media) == 0 {
		b.respondWithError(s, i, fmt.Sprintf("❌ No anime found for query: \"%s\"", query))
		return
	}

	// Create embeds for search results
	var embeds []*discordgo.MessageEmbed
	for i, anime := range searchResults.Data.Page.Media {
		if i >= 3 { // Limit to 3 results
			break
		}

		title := anime.Title.Romaji
		if anime.Title.English != nil && *anime.Title.English != "" {
			title = *anime.Title.English
		}

		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("📺 %s", title),
			URL:   anime.SiteURL,
			Color: 0x0099FF,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: anime.CoverImage.Large,
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Romaji Title", Value: anime.Title.Romaji, Inline: true},
				{Name: "Native Title", Value: anime.Title.Native, Inline: true},
				{Name: "Format", Value: anime.Format, Inline: true},
				{Name: "Status", Value: anime.Status, Inline: true},
				{Name: "AniList ID", Value: fmt.Sprintf("%d", anime.ID), Inline: true},
			},
		}
		embeds = append(embeds, embed)
	}

	responseText := fmt.Sprintf("🔍 **Search Results for:** \"%s\"\n\nFound %d results (showing top %d):",
		query, searchResults.Data.Page.PageInfo.Total, len(embeds))

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &responseText,
		Embeds:  &embeds,
	})
	if err != nil {
		log.Printf("Failed to edit interaction response: %v", err)
	}
}

// respondWithError responds to an interaction with an error message
func (b *Bot) respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &message,
	})
	if err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}

// handleReleaseCommand handles the anime release subcommand
func (b *Bot) handleReleaseCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
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
		b.respondWithError(s, i, fmt.Sprintf("❌ No anime found with ID %d.", animeID))
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
