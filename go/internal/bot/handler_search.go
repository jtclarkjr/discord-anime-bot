package bot

import (
	"fmt"
	"log"

	"discord-anime-bot/internal/services/anilist"

	"github.com/bwmarrin/discordgo"
)

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
		b.respondWithError(s, i, "An error occurred while searching for anime.")
		return
	}

	if len(searchResults.Data.Page.Media) == 0 {
		b.respondWithError(s, i, fmt.Sprintf("No anime found for query: \"%s\"", query))
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
