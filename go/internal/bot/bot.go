package bot

import (
	"fmt"
	"log"

	"discord-anime-bot/internal/config"
	"discord-anime-bot/internal/services/anilist"

	"github.com/bwmarrin/discordgo"
)

// Bot represents the Discord bot instance
type Bot struct {
	session             *discordgo.Session
	config              *config.Config
	notificationService *anilist.NotificationService
}

// NewBot creates a new bot instance
func NewBot(cfg *config.Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// Initialize notification service
	notificationService := anilist.NewNotificationService(session)

	bot := &Bot{
		session:             session,
		config:              cfg,
		notificationService: notificationService,
	}

	// Add event handlers
	session.AddHandler(bot.ready)
	session.AddHandler(bot.interactionCreate)

	// Set intents
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	return bot, nil
}

// Start starts the bot
func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %w", err)
	}
	return nil
}

// Stop stops the bot
func (b *Bot) Stop() {
	if b.notificationService != nil {
		b.notificationService.Cleanup()
	}
	if b.session != nil {
		b.session.Close()
	}
}

// ready is called when the bot is ready
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s!", s.State.User.Username)

	// Create base command options
	commandOptions := []*discordgo.ApplicationCommandOption{
		{
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
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "release",
			Description: "Get currently releasing anime",
		},
		{
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
		},
		{
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
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        "notify",
			Description: "Episode notification commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "Set notification for next episode",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "id",
							Description: "The AniList ID of the anime",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "List your active episode notifications",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "cancel",
					Description: "Cancel notification for an anime",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "id",
							Description: "The AniList ID of the anime",
							Required:    true,
						},
					},
				},
			},
		},
	}

	// Conditionally add the find command if OpenAI is enabled
	if b.config.IsOpenAIEnabled {
		findOption := &discordgo.ApplicationCommandOption{
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
		// Prepend find command to the beginning for better UX
		commandOptions = append([]*discordgo.ApplicationCommandOption{findOption}, commandOptions...)
	}

	// Register the /anime command with subcommands
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "anime",
			Description: "Anime-related commands",
			Options:     commandOptions,
		},
	}

	// Register commands to the first guild (for testing)
	// In production, you might want to register globally
	guilds := s.State.Guilds
	if len(guilds) > 0 {
		for _, cmd := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guilds[0].ID, cmd)
			if err != nil {
				log.Printf("❌ Failed to register command %s: %v", cmd.Name, err)
			} else {
				log.Printf("✅ Successfully registered command: %s", cmd.Name)
			}
		}
	} else {
		log.Println("⚠️ No guilds found. Commands may not be registered.")
	}
}
