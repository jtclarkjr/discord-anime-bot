package anime

import "github.com/bwmarrin/discordgo"

// GetFindCommandOption returns the find command option
func GetFindCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "find",
		Description: "Find anime by description using AI",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "Describe the anime you're looking for",
				Required:    true,
			},
		},
	}
}