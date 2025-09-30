package bot

import (
	"fmt"
	"log"
	"math"
	"strings"

	"discord-anime-bot/internal/services/anilist"

	"github.com/bwmarrin/discordgo"
)

// handleFindCommand handles the anime find subcommand
func (b *Bot) handleFindCommand(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	if !b.config.IsOpenAIEnabled {
		b.respondWithError(s, i, "The find command is disabled because OpenAI API key is not configured. Please set the OPENAI_API_KEY environment variable to use AI-powered anime search.")
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

	// Find anime using AI (OpenAI or Claude, based on config)
	matches, err := anilist.FindAnimeWithDetails(prompt, b.config)
	if err != nil {
		log.Printf("Error finding anime: %v", err)
		b.respondWithError(s, i, "An error occurred while searching for anime.")
		return
	}

	if len(matches) == 0 {
		b.respondWithError(s, i, fmt.Sprintf("No anime found matching the description: \"%s\"", prompt))
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
