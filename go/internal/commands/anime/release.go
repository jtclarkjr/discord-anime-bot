package anime

import "github.com/bwmarrin/discordgo"

// GetReleaseCommandOption returns the release command option
func GetReleaseCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "release",
		Description: "Show all currently releasing anime",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "page",
				Description: "Page number",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "perpage",
				Description: "Anime per page (max 50)",
				Required:    false,
			},
		},
	}
}
