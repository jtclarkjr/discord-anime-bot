package anime

import "github.com/bwmarrin/discordgo"

// GetHelpCommandOption returns the help command option
func GetHelpCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "help",
		Description: "Show help for all /anime commands",
	}
}