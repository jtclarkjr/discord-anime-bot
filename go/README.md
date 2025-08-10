# Discord Anime Bot - Go Version

A Discord bot written in Go that helps users find and search for anime using AI-powered descriptions and the AniList API.

## Features

- **AI-Powered Anime Search**: Use natural language descriptions to find anime with GPT-5
- **Traditional Search**: Search anime by title using AniList API
- **Rich Discord Embeds**: Beautiful embedded responses with anime details
- **Slash Commands**: Modern Discord slash command interface

## Commands

### `/anime find <prompt>`

Find anime using AI based on a description.

**Example**: `/anime find "anime about a kid who becomes a pirate"`

### `/anime search <query>`

Search for anime by title.

**Example**: `/anime search "One Piece"`

## Setup

### Prerequisites

- Go 1.21 or higher
- Discord Bot Token
- OpenAI API Key
- AniList API endpoint (typically `https://graphql.anilist.co`)

### Environment Variables

Create a `.env` file in the root directory:

```env
DISCORD_BOT_TOKEN=your_discord_bot_token_here
CHANNEL_ID=your_channel_id_here
ANILIST_API=https://graphql.anilist.co
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
│   │   ├── bot.go                  # Bot initialization
│   │   └── handlers.go             # Command handlers
│   ├── config/                     # Configuration management
│   │   └── config.go
│   ├── services/                   # External service integrations
│   │   ├── anilist/                # AniList API integration
│   │   │   ├── search.go
│   │   │   └── find.go
│   │   └── openai/                 # OpenAI API integration
│   │       └── completions.go
│   └── types/                      # Type definitions
│       ├── anilist.go
│       └── openai.go
├── go.mod                          # Go module definition
├── go.sum                          # Dependency checksums
├── Makefile                        # Build automation
└── README.md                       # This file
```

## Dependencies

- **discordgo**: Discord API library for Go
- **go-openai**: OpenAI API client for Go
- **godotenv**: Environment variable loading

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see the original project's LICENSE file for details.
