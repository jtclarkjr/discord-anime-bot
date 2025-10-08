package anime

import "github.com/bwmarrin/discordgo"

// GetSearchCommandOption returns the search command option
func GetSearchCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "search",
		Description: "Search for anime by title",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "The anime title to search for",
				Required:    true,
			},
		},
	}
}