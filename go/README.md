# Discord Anime Bot - Go Version

A Discord bot written in Go that helps users find, search, and track anime using AI-powered descriptions and the AniList API. Features episode notifications with persistent storage.

## Features

- **AI-Powered Anime Search**: Use natural language descriptions to find anime with GPT-5 _(requires OpenAI API key)_
- **Traditional Search**: Search anime by title using AniList API
- **Episode Notifications**: Get notified when new anime episodes air
- **Currently Releasing**: View currently airing anime with schedules
- **Next Episode Info**: Check when the next episode of any anime airs
- **Persistent Storage**: Notifications persist across bot restarts using JSON storage
- **Rich Discord Embeds**: Beautiful embedded responses with anime details
- **Slash Commands**: Modern Discord slash command interface

## Commands

### `/anime find <prompt>`

Find anime using AI based on a description. _(Requires OpenAI API key)_

**Example**: `/anime find "anime about a kid who becomes a pirate"`

### `/anime search <query>`

Search for anime by title.

**Example**: `/anime search "One Piece"`

### `/anime release`

Display currently releasing anime with their next episode schedules.

### `/anime next <id>`

Get information about the next episode of a specific anime.

**Example**: `/anime next 21` (for One Piece)

### `/anime notify` commands

Manage episode notifications:

- `/anime notify add <id>` - Set notification for next episode
- `/anime notify list` - View your active notifications
- `/anime notify cancel <id>` - Cancel notification for an anime

**Examples**:

- `/anime notify add 21` - Get notified for One Piece episodes
- `/anime notify list` - See all your notifications
- `/anime notify cancel 21` - Stop One Piece notifications

## Setup

### Prerequisites

- Go 1.21 or higher
- Discord Bot Token
- AniList API endpoint (typically `https://graphql.anilist.co`)
- OpenAI API Key (optional, for AI-powered search features)

### Environment Variables

Create a `.env` file in the root directory:

**Required:**

```env
DISCORD_BOT_TOKEN=your_discord_bot_token_here
CHANNEL_ID=your_channel_id_here
ANILIST_API=https://graphql.anilist.co
```

**Optional (for AI features):**

```env
OPENAI_API_KEY=your_openai_api_key_here
```

### Installation

1. Clone the repository
2. Navigate to the `go` directory
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build and run:
   ```bash
   make build
   make run
   ```

### Development

- `make dev` - Run in development mode with auto-reload
- `make build` - Build the binary
- `make clean` - Clean build artifacts
- `make test` - Run tests

## Project Structure

```
go/
├── main.go                          # Entry point
├── internal/
│   ├── bot/                        # Discord bot logic
│   │   ├── bot.go                  # Bot initialization and setup
│   │   ├── handler_main.go         # Main interaction router
│   │   ├── handler_find.go         # AI-powered anime search
│   │   ├── handler_search.go       # Traditional anime search
│   │   ├── handler_release.go      # Currently releasing anime
│   │   ├── handler_next.go         # Next episode information
│   │   └── handler_notify.go       # Episode notification system
│   ├── config/                     # Configuration management
│   │   └── config.go
│   ├── services/                   # External service integrations
│   │   ├── anilist/                # AniList API integration
│   │   │   ├── search.go           # Anime search functionality
│   │   │   ├── find.go             # AI-powered search
│   │   │   ├── release.go          # Currently releasing anime
│   │   │   ├── next.go             # Next episode data
│   │   │   └── notify.go           # Notification service
│   │   └── openai/                 # OpenAI API integration
│   │       └── completions.go
│   ├── types/                      # Type definitions
│   │   ├── anilist.go              # AniList API types
│   │   └── openai.go               # OpenAI API types
│   └── utils/                      # Utility functions
│       └── formatters.go           # Time and date formatting
├── data/                           # Data storage
│   └── notifications.json         # Persistent notification storage
├── go.mod                          # Go module definition
├── go.sum                          # Dependency checksums
├── Makefile                        # Build automation
└── README.md                       # This file
```

## Dependencies

- **discordgo**: Discord API library for Go
- **go-openai**: OpenAI API client for Go
- **godotenv**: Environment variable loading

## Architecture

The bot uses a modular architecture with clear separation of concerns:

- **Handlers**: Each command type has its own handler file for maintainability
- **Services**: External API integrations (AniList, OpenAI) are encapsulated
- **Persistent Storage**: JSON-based storage for notifications with automatic cleanup
- **Concurrent Safety**: Thread-safe notification management with proper synchronization
- **Error Handling**: Comprehensive error handling with user-friendly Discord responses

## Notification System

The notification system provides:

- **Persistent Storage**: Notifications survive bot restarts
- **Automatic Scheduling**: Uses Go's `time.AfterFunc` for precise timing
- **Cleanup**: Automatic removal of expired notifications
- **User Management**: Per-user notification tracking
- **Memory Efficient**: Minimal memory footprint with cleanup on notification fire

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see the original project's LICENSE file for details.
