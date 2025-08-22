# Discord Anime Bot - Go Version

A Discord bot written in Go that helps users find, search, and track anime using AI-powered descriptions and the AniList API. Features episode notifications with persistent storage.

## Features

- **AI-Powered Anime Search**: Use natural language descriptions to find anime with GPT-5 _(requires OpenAI/ Claude API key)_
- **Traditional Search**: Search anime by title using AniList API
- **Episode Notifications**: Get notified when new anime episodes air
- **Watchlist Management**: Track your personal anime watchlist
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

Search for anime by title or AniList ID.

**Example**: `/anime search "One Piece"`

### `/anime release`

Display currently releasing anime with their next episode schedules.

Supports pagination:

- `/anime release page:<number> perpage:<number>`
  - `page` (optional): Page number to view (default: 1)
  - `perpage` (optional): Number of anime per page (default: 15, max: 50)
  - Example: `/anime release page:2 perpage:25`
    Shows page info and total count in the embed footer.

### `/anime season <season> [year]`

Get all anime from a specific season and year.

**Examples**:

- `/anime season summer` - Shows all Summer 2025 anime
- `/anime season winter 2023` - Shows all Winter 2023 anime
- `/anime season fall 2024` - Shows all Fall 2024 anime

### `/anime next <id>`

Get information about the next episode of a specific anime.

**Example**: `/anime next 21` (for One Piece)

### `/anime notify` commands

Manage episode notifications:

- `/anime notify` - View your active notifications (default)
- `/anime notify action:add id:<id>` - Set notification for next episode
- `/anime notify action:cancel id:<id>` - Cancel notification for an anime

**Examples**:

- `/anime notify` - See all your current notifications
- `/anime notify action:add id:21` - Get notified for One Piece episodes
- `/anime notify action:cancel id:21` - Stop One Piece notifications

### `/anime watchlist` commands

Manage your personal anime watchlist:

- `/anime watchlist` - View your anime watchlist (default)
- `/anime watchlist action:add id:<id>` - Add an anime to your watchlist
- `/anime watchlist action:remove id:<id>` - Remove an anime from your watchlist

**Examples**:

- `/anime watchlist` - See your current watchlist
- `/anime watchlist action:add id:21` - Add One Piece to your watchlist
- `/anime watchlist action:remove id:21` - Remove One Piece from your watchlist

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
STORAGE_FILE=./data/notifications.json
WATCHLIST_FILE=./data/watchlists.json
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
│   │   ├── handler_season.go       # Seasonal anime listings
│   │   ├── handler_next.go         # Next episode information
│   │   ├── handler_notify.go       # Episode notification system
│   │   ├── handler_watchlist.go    # Watchlist management
│   │   ├── handler_help.go         # Help command handler
│   ├── config/                     # Configuration management
│   │   └── config.go
│   ├── graphql/                    # GraphQL query definitions
│   │   ├── search_by_id.go         # Anime search by ID query
│   │   ├── search_by_text.go       # Anime text search query
│   │   ├── anime_details.go        # Anime details with next episode query
│   │   ├── releasing_anime.go      # Currently releasing anime query
│   │   └── seasonal_anime.go       # Seasonal anime query
│   ├── services/                   # External service integrations
│   │   ├── anilist/                # AniList API integration
│   │   │   ├── search.go           # Anime search functionality
│   │   │   ├── find.go             # AI-powered search
│   │   │   ├── release.go          # Currently releasing anime
│   │   │   ├── season.go           # Seasonal anime data
│   │   │   ├── next.go             # Next episode data
│   │   │   ├── notify.go           # Notification service
│   │   │   ├── watchlist.go        # Watchlist service
│   │   └── claude/                 # Claude API integration
│   │   │   └── claude.go           # Claude completions
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

## License

MIT License - see the original project's LICENSE file for details.
