package anime

import "github.com/bwmarrin/discordgo"

// GetNotifyCommandOption returns the notify command option
func GetNotifyCommandOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "notify",
		Description: "Manage episode notifications (shows list by default)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "action",
				Description: "Action to perform (add or cancel)",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Add notification",
						Value: "add",
					},
					{
						Name:  "Cancel notification",
						Value: "cancel",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "id",
				Description: "AniList ID of the anime (required for add/cancel)",
				Required:    false,
			},
		},
	}
}