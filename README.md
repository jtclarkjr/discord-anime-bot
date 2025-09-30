# Discord Anime Bot

A Discord bot for searching anime information and tracking release schedules using the AniList API.

## Commands

- `/anime help` - Show all available anime commands, subcommands, and arguments
- `/anime search <query>` - Search for anime by name or AniList ID
- `/anime next <id>` - Get next episode information for an anime by AniList ID
- `/anime notify` - List your active episode notifications (default)
- `/anime notify action:add id:<id>` - Set up notification for when the next episode airs
- `/anime notify action:cancel id:<id>` - Cancel notification for an anime
- `/anime watchlist` - Show your personal anime watchlist (default, only visible to you)
- `/anime watchlist action:add id:<id>` - Add an anime to your personal watchlist
- `/anime watchlist action:remove id:<id>` - Remove an anime from your personal watchlist
- `/anime release [page] [perpage]` - Show currently releasing anime with pagination
- `/anime season <season> [year]` - Get all anime from a specific season and year
- `/anime find <prompt>` - Find anime using AI based on description (powered by GPT-5) _(requires OpenAI/ Claude API key)_

## Help Command

Use `/anime help` to see a full list of all available commands, subcommands, and their arguments. The help output is always up-to-date with the bot's features and shows usage for every command, including subcommand groups and required/optional arguments.

## Project Structure

```
src/
├── config/
│   └── constants.ts        # Environment variables and configuration constants
├── graphql/
│   ├── searchById.ts       # GraphQL query for anime search by ID
│   ├── searchByText.ts     # GraphQL query for anime text search
│   ├── animeDetails.ts     # GraphQL query for anime details with next episode
│   ├── releasingAnime.ts   # GraphQL query for currently releasing anime
│   ├── seasonalAnime.ts    # GraphQL query for seasonal anime
│   └── index.ts            # GraphQL queries exports
├── types/
│   ├── anilist.ts          # TypeScript type definitions for AniList API responses
│   └── openai.ts           # TypeScript type definitions for OpenAI API responses
├── services/
│   ├── anime/
│   │   ├── search.ts       # Search anime API service
│   │   ├── next.ts         # Next episode API service
│   │   ├── release.ts      # Releasing anime API service
│   │   ├── season.ts       # Seasonal anime API service
│   │   ├── find.ts         # AI-powered anime finder service
│   │   ├── notify.ts       # Episode notification service
│   │   ├── watchlist.ts    # Watchlist anime API service
│   │   └── index.ts        # Anime services exports
│   ├── claude/
│   │   ├── claude.ts       # Claude completions service
│   │   └── index.ts        # Claude services exports
│   ├── openai/
│   │   ├── completions.ts  # OpenAI ChatGPT completion service
│   │   └── index.ts        # OpenAI services exports
│   └── index.ts            # All services exports
├── commands/
│   ├── anime/
│   │   ├── search.ts       # Search anime command handler
│   │   ├── next.ts         # Next episode command handler
│   │   ├── release.ts      # Releasing anime command handler
│   │   ├── season.ts       # Seasonal anime command handler
│   │   ├── find.ts         # AI find anime command handler
│   │   ├── notify.ts       # Episode notification command handler
│   │   ├── watchlist.ts    # Watchlist command handler
│   │   ├── help.ts         # Help command handler
│   │   └── index.ts        # Command routing and definitions
├── utils/
│   └── formatters.ts       # Utility functions for formatting dates and times
└── index.ts                # Main bot entry point
```

## Features

### Search Command

- Search for anime by name
- Returns detailed information including title, format, status, and AniList ID
- Shows multiple matches when available

### Next Episode Command

- Get next airing episode information by AniList ID
- Shows countdown timer (days, hours, minutes)
- Displays air date in user-friendly format
- Handles finished/cancelled anime appropriately

### Notification Commands

- **List (Default)**: View all your active episode notifications with `/anime notify`
- **Add**: Set up automatic notifications for when episodes air with `/anime notify action:add id:<id>`
- **Cancel**: Remove notifications for specific anime with `/anime notify action:cancel id:<id>`
- Notifications are sent automatically when episodes air
- One notification per anime per user is maintained
- **Examples**:
  - `/anime notify` - Shows your current notifications
  - `/anime notify action:add id:21` - Get notified for One Piece episodes
  - `/anime notify action:cancel id:21` - Stop One Piece notifications

### Release Command

- Lists currently releasing anime
- Shows next episode numbers and air dates
- Sorted by popularity
- Supports pagination: use `/anime release page:<number> perpage:<number>`
- `page` (optional): Page number to view (default: 1)
- `perpage` (optional): Number of anime per page (default: 15, max: 50)
- Displays up to `perpage` anime per page, with navigation info

### Season Command

- Get all anime from a specific season and year
- Supports all four seasons: Winter, Spring, Summer, Fall
- Year parameter is optional (defaults to current year)
- Shows **complete** seasonal listings (not truncated like other commands)
- Displays status indicators:
  - Currently Releasing
  - Finished
  - Not Yet Released
  - Cancelled
  - On Hiatus
- Automatically handles multiple embeds for large seasonal catalogs
- Sorted by popularity from AniList
- Examples:
  - `/anime season summer` - Shows all Summer 2025 anime
  - `/anime season winter 2023` - Shows all Winter 2023 anime
  - `/anime season fall 2024` - Shows all Fall 2024 anime

### Find Command (AI-Powered)

- Uses ChatGPT-5 to understand natural language descriptions
- Finds anime based on plot, genre, themes, or characteristics
- Returns AniList details for the best matches
- Shows confidence scores and reasoning
- Examples:
  - "A show about a boy who can turn into a titan"
  - "Romance anime set in high school with supernatural elements"
  - "Sci-fi mecha anime with philosophical themes"
- Displays up to 15 anime with pagination info

## Episode Notifications

The bot includes an automatic notification system for new episode releases:

- **Easy Access**: Simply use `/anime notify` to see all your active notifications
- **Set Notifications**: Use `/anime notify action:add id:<id>` to get notified when an episode airs
- **Automatic Delivery**: Notifications are sent in the same channel where you set them up
- **Smart Scheduling**: The bot calculates exact air times and schedules notifications accordingly
- **One Per Anime**: Only one notification per anime per user is maintained
- **Easy Management**: Cancel notifications with `/anime notify action:cancel id:<id>`
- **Error Handling**: Gracefully handles finished anime, invalid IDs, and scheduling conflicts

## Setup

### Local Development

1. Install dependencies:

   ```bash
   bun i
   ```

2. Set up Redis:

   ```bash
   # Using Docker (recommended)
   docker run -d --name redis -p 6379:6379 redis:7-alpine

   # Or install locally (macOS)
   brew install redis
   brew services start redis
   ```

3. Create a `.env` file with:

   ```bash
   DISCORD_BOT_TOKEN=your_discord_bot_token
   CHANNEL_ID=your_channel_id
   ANILIST_API=https://graphql.anilist.co
   REDIS_URL=redis://localhost:6379
   OPENAI_API_KEY=your_openai_api_key  # Optional
   CLAUDE_API_KEY=your_claude_api_key  # Optional
   ```

4. Test Redis connection:

   ```bash
   bun run test:redis
   ```

5. Run the bot:
   ```bash
   bun run dev
   ```

### Docker Deployment

For production or simplified setup, use Docker:

```bash
# Copy environment template
cp .env.example .env
# Edit .env with your configuration

# Build and start services (includes Redis)
docker-compose up --build

# Stop services
docker-compose down
```

See [DOCKER.md](./DOCKER.md) for complete Docker documentation and [DEPLOYMENT.md](./DEPLOYMENT.md) for production deployment options including Google Cloud Run.

## Environment Variables

**Required:**

- `DISCORD_BOT_TOKEN` - Your Discord bot token
- `CHANNEL_ID` - The Discord channel ID where the bot operates
- `ANILIST_API` - AniList GraphQL API endpoint (https://graphql.anilist.co)
- `STORAGE_FILE` - Path to notifications storage file (default: `./data/notifications.json`)
- `WATCHLIST_FILE` - Path to watchlist storage file (default: `./data/watchlists.json`)

**Optional:**

- `OPENAI_API_KEY` - Your OpenAI API key for ChatGPT-5 access (enables `/anime find` command)

## Technologies Used

- **Bun** - JavaScript runtime and package manager
- **Discord.js** - Discord API library
- **TypeScript** - Type-safe JavaScript
- **AniList API** - Anime and manga database API
- **OpenAI ChatGPT-5** - AI-powered anime recommendations
