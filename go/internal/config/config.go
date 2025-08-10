package config

import (
	"log"
	"os"
)

// Config holds all configuration values for the bot
type Config struct {
	DiscordToken string
	ChannelID    string
	AniListAPI   string
	OpenAIAPIKey string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		DiscordToken: getEnv("DISCORD_BOT_TOKEN"),
		ChannelID:    getEnv("CHANNEL_ID"),
		AniListAPI:   getEnv("ANILIST_API"),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY"),
	}

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
	if cfg.OpenAIAPIKey == "" {
		log.Fatal("❌ OPENAI_API_KEY is not set in environment variables.")
	}
}
