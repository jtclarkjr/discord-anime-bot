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
	ClaudeAPIKey    string
	IsOpenAIEnabled bool
	IsClaudeEnabled bool
	IsAIEnabled     bool
	UseOpenAI       bool // true if OpenAI should be used, false if Claude should be used
	StorageFile     string
	WatchlistFile   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		DiscordToken:  getEnv("DISCORD_BOT_TOKEN"),
		ChannelID:     getEnv("CHANNEL_ID"),
		AniListAPI:    getEnv("ANILIST_API"),
		OpenAIAPIKey:  getEnvOptional("OPENAI_API_KEY"),
		ClaudeAPIKey:  getEnvOptional("CLAUDE_API_KEY"),
		StorageFile:   getEnv("STORAGE_FILE"),
		WatchlistFile: getEnv("WATCHLIST_FILE"),
	}

	cfg.IsOpenAIEnabled = cfg.OpenAIAPIKey != ""
	cfg.IsClaudeEnabled = cfg.ClaudeAPIKey != ""
	cfg.IsAIEnabled = cfg.IsOpenAIEnabled || cfg.IsClaudeEnabled
	cfg.UseOpenAI = cfg.IsOpenAIEnabled || (!cfg.IsOpenAIEnabled && !cfg.IsClaudeEnabled)

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

	// AI logic
	if cfg.IsOpenAIEnabled {
		log.Println("✅ OpenAI features enabled (find command available)")
	} else if cfg.IsClaudeEnabled {
		log.Println("✅ Claude features enabled (find command available)")
	} else {
		log.Println("⚠️ No AI features enabled (find command not available)")
	}
}
