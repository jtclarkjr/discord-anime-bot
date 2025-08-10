# Discord Anime Bot

A Discord bot for searching anime information and tracking release schedules using the AniList API.

## Commands

- `/anime search <query>` - Search for anime by name
- `/anime next <id>` - Get next episode information for an anime by AniList ID
- `/anime release` - Show all currently releasing anime
- `/anime find <prompt>` - Find anime using AI based on description (powered by ChatGPT-5)

## Project Structure

```
src/
├── config/
│   └── constants.ts        # Environment variables and configuration constants
├── types/
│   ├── anilist.ts          # TypeScript type definitions for AniList API responses
│   └── openai.ts           # TypeScript type definitions for OpenAI API responses
├── services/
│   ├── anime/
│   │   ├── search.ts       # Search anime API service
│   │   ├── next.ts         # Next episode API service
│   │   ├── release.ts      # Releasing anime API service
│   │   ├── find.ts         # AI-powered anime finder service
│   │   └── index.ts        # Anime services exports
│   ├── openai/
│   │   ├── completions.ts  # OpenAI ChatGPT completion service
│   │   └── index.ts        # OpenAI services exports
│   └── index.ts            # All services exports
├── commands/
│   ├── anime/
│   │   ├── search.ts       # Search anime command handler
│   │   ├── next.ts         # Next episode command handler
│   │   ├── release.ts      # Releasing anime command handler
│   │   ├── find.ts         # AI find anime command handler
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

### Release Command

- Lists currently releasing anime
- Shows next episode numbers and air dates
- Sorted by popularity
- Displays up to 15 anime with pagination info

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

## Setup

1. Install dependencies:

   ```bash
   bun i
   ```

2. Create a `.env` file with:

   ```
   DISCORD_BOT_TOKEN=your_discord_bot_token
   CHANNEL_ID=your_channel_id
   ANILIST_API=https://graphql.anilist.co
   OPENAI_API_KEY=your_openai_api_key
   ```

3. Run the bot:
   ```bash
   bun dev
   docker compose up --build
   ```

## Environment Variables

- `DISCORD_BOT_TOKEN` - Your Discord bot token
- `CHANNEL_ID` - The Discord channel ID where the bot operates
- `ANILIST_API` - AniList GraphQL API endpoint (https://graphql.anilist.co)
- `OPENAI_API_KEY` - Your OpenAI API key for ChatGPT-5 access

## Technologies Used

- **Bun** - JavaScript runtime and package manager
- **Discord.js** - Discord API library
- **TypeScript** - Type-safe JavaScript
- **AniList API** - Anime and manga database API
- **OpenAI ChatGPT-5** - AI-powered anime recommendations
