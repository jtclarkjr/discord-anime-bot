package bot

import (
	"fmt"
	"log"

	"discord-anime-bot/internal/commands"
	"discord-anime-bot/internal/config"
	"discord-anime-bot/internal/services/anilist"
	"discord-anime-bot/internal/services/redis"

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
	// Initialize Redis connection
	if err := redis.InitRedis(cfg.RedisURL); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		log.Println("Bot will continue without Redis (features may not work properly)")
	}

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
		if err := b.session.Close(); err != nil {
			log.Printf("Error closing Discord session: %v", err)
		}
	}
	// Close Redis connection
	if err := redis.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}
}

// ready is called when the bot is ready
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as %s!", s.State.User.Username)

	// Get all commands from the commands package
	commands := commands.GetAllCommands(b.config)

	// Register commands to the first guild (for testing)
	// In production, you might want to register globally
	guilds := s.State.Guilds
	if len(guilds) > 0 {
		for _, cmd := range commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guilds[0].ID, cmd)
			if err != nil {
				log.Printf("Failed to register command %s: %v", cmd.Name, err)
			} else {
				log.Printf("Successfully registered command: %s", cmd.Name)
			}
		}
	} else {
		log.Println("Warning: No guilds found. Commands may not be registered.")
	}
}
