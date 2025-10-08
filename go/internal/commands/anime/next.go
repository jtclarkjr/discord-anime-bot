package anime

import "github.com/bwmarrin/discordgo"

// GetNextCommandOption returns the next command option
func GetNextCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "next",
		Description: "Get next episode information for an anime",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "The AniList ID of the anime",
				Required:    true,
			},
		},
	}
}