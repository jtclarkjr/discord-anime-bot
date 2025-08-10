package bot

import (
	"log"

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
		b.handleReleaseCommand(s, i)
	case "next":
		b.handleNextCommand(s, i, subcommand.Options)
	case "notify":
		b.handleNotifyCommand(s, i, subcommand.Options)
	default:
		b.respondWithError(s, i, "Unknown subcommand")
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
