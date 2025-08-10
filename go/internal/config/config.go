package config

import (
	"log"
	"os"
)

// Config holds all configuration values for the bot
type Config struct {
	DiscordToken    string
	ChannelID       string
	AniListAPI      string
	OpenAIAPIKey    string
	IsOpenAIEnabled bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		DiscordToken: getEnv("DISCORD_BOT_TOKEN"),
		ChannelID:    getEnv("CHANNEL_ID"),
		AniListAPI:   getEnv("ANILIST_API"),
		OpenAIAPIKey: getEnvOptional("OPENAI_API_KEY"),
	}

	cfg.IsOpenAIEnabled = cfg.OpenAIAPIKey != ""

	// Validate required environment variables
	validateConfig(cfg)

	return cfg
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("⚠️  Warning: %s environment variable is not set", key)
	}
	return value
}

func getEnvOptional(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("⚠️ %s environment variable is not set. Related features will be disabled.", key)
	}
	return value
}

func validateConfig(cfg *Config) {
	if cfg.DiscordToken == "" {
		log.Fatal("❌ DISCORD_BOT_TOKEN is not set in environment variables.")
	}
	if cfg.ChannelID == "" {
		log.Fatal("❌ CHANNEL_ID is not set in environment variables.")
	}
	if cfg.AniListAPI == "" {
		log.Fatal("❌ ANILIST_API is not set in environment variables.")
	}

	// OpenAI is optional - if not provided, AI features will be disabled
	if cfg.IsOpenAIEnabled {
		log.Println("✅ OpenAI features enabled (find command available)")
	} else {
		log.Println("⚠️ OpenAI features disabled (find command not available)")
	}
}
