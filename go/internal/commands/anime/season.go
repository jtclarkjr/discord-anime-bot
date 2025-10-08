package anime

import "github.com/bwmarrin/discordgo"

// GetSeasonCommandOption returns the season command option
func GetSeasonCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "season",
		Description: "Get all anime from a specific season and year",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "season",
				Description: "Season (winter, spring, summer, fall)",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Winter",
						Value: "winter",
					},
					{
						Name:  "Spring",
						Value: "spring",
					},
					{
						Name:  "Summer",
						Value: "summer",
					},
					{
						Name:  "Fall",
						Value: "fall",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "year",
				Description: "Year (defaults to current year)",
				Required:    false,
			},
		},
	}
}