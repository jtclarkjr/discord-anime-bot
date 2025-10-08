package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// handleHelpCommand responds with a list of all /anime commands and their arguments
func (b *Bot) handleHelpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	helpLines := []string{
		"Here are the available /anime commands:",
		"",
		"**/anime help**: Show help for all /anime commands",
		"**/anime search <query>**: Search for anime by title",
		"**/anime release**: Get currently releasing anime",
		"**/anime season <season> [year]**: Get all anime from a specific season and year",
		"**/anime next <id>**: Get next episode information for an anime",
		"**/anime notify add <id>**: Set notification for next episode",
		"**/anime notify list**: List your active episode notifications",
		"**/anime notify cancel <id>**: Cancel notification for an anime",
		"**/anime watchlist add <id>**: Add an anime to your personal watchlist",
		"**/anime watchlist list**: Show your personal anime watchlist (only visible to you)",
		"**/anime watchlist remove <id>**: Remove an anime from your personal watchlist",
	}
	if b.config.IsOpenAIEnabled {
		helpLines = append(helpLines, "**/anime find <prompt>**: Find anime by description using AI")
	}
	helpText := ""
	for _, line := range helpLines {
		helpText += line + "\n"
	}
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &helpText,
	})
	if err != nil {
		log.Printf("Failed to send help response: %v", err)
	}
}
