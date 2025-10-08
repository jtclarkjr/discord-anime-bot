package anime

import (
	"discord-anime-bot/internal/config"
	"github.com/bwmarrin/discordgo"
)

// GetAnimeCommandOptions returns all anime subcommand options
func GetAnimeCommandOptions(cfg *config.Config) []*discordgo.ApplicationCommandOption {
	// Base command options
	commandOptions := []*discordgo.ApplicationCommandOption{
		GetSearchCommandOption(),
		GetNextCommandOption(),
		GetNotifyCommandOption(),
		GetWatchlistCommandOption(),
		GetReleaseCommandOption(),
		GetSeasonCommandOption(),
	}

	// Conditionally add the find command if OpenAI is enabled
	if cfg.IsOpenAIEnabled {
		// Prepend find command to the beginning for better UX
		commandOptions = append([]*discordgo.ApplicationCommandOption{GetFindCommandOption()}, commandOptions...)
	}

	// Add help command at the end
	commandOptions = append(commandOptions, GetHelpCommandOption())

	return commandOptions
}

// GetAnimeCommand returns the complete anime command definition
func GetAnimeCommand(cfg *config.Config) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "anime",
		Description: "Anime-related commands",
		Options:     GetAnimeCommandOptions(cfg),
	}
}