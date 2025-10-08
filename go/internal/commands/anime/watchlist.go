package anime

import "github.com/bwmarrin/discordgo"

// GetWatchlistCommandOption returns the watchlist command option
func GetWatchlistCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "watchlist",
		Description: "Manage your anime watchlist (shows list by default)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "action",
				Description: "Action to perform (add or remove)",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Add to watchlist",
						Value: "add",
					},
					{
						Name:  "Remove from watchlist",
						Value: "remove",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "AniList ID of the anime (required for add/remove)",
				Required:    false,
			},
		},
	}
}