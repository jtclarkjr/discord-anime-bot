package commands

import (
	"discord-anime-bot/internal/commands/anime"
	"discord-anime-bot/internal/config"
	"github.com/bwmarrin/discordgo"
)

// GetAllCommands returns all Discord application commands
func GetAllCommands(cfg *config.Config) []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		anime.GetAnimeCommand(cfg),
	}
}