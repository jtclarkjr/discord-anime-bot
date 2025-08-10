# Discord Anime Bot

A Discord bot for searching anime information and tracking release schedules using the AniList API.

## Commands

- `/anime search <query>` - Search for anime by name
- `/anime next <id>` - Get next episode information for an anime by AniList ID
- `/anime notify add <id>` - Set up notification for when the next episode airs
- `/anime notify list` - List your active episode notifications
- `/anime notify cancel <id>` - Cancel notification for an anime
- `/anime release` - Show all currently releasing anime
- `/anime find <prompt>` - Find anime using AI based on description (powered by GPT-5) _(requires OpenAI API key)_

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
│   │   ├── notify.ts       # Episode notification service
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

### Notification Commands

- **Add**: Set up automatic notifications for when episodes air with `/anime notify add <id>`
- **List**: View all your active episode notifications with `/anime notify list`
- **Cancel**: Remove notifications for specific anime with `/anime notify cancel <id>`
- Notifications are sent automatically when episodes air
- One notification per anime per user is maintained

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

## Episode Notifications

The bot includes an automatic notification system for new episode releases:

- **Set Notifications**: Use `/anime notify add <id>` to get notified when an episode airs
- **Automatic Delivery**: Notifications are sent in the same channel where you set them up
- **Smart Scheduling**: The bot calculates exact air times and schedules notifications accordingly
- **One Per Anime**: Only one notification per anime per user is maintained
- **Easy Management**: List active notifications with `/anime notify list` and cancel with `/anime notify cancel <id>`
- **Error Handling**: Gracefully handles finished anime, invalid IDs, and scheduling conflicts

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

**Required:**

- `DISCORD_BOT_TOKEN` - Your Discord bot token
- `CHANNEL_ID` - The Discord channel ID where the bot operates
- `ANILIST_API` - AniList GraphQL API endpoint (https://graphql.anilist.co)

**Optional:**

- `OPENAI_API_KEY` - Your OpenAI API key for ChatGPT-5 access (enables `/anime find` command)

## Technologies Used

- **Bun** - JavaScript runtime and package manager
- **Discord.js** - Discord API library
- **TypeScript** - Type-safe JavaScript
- **AniList API** - Anime and manga database API
- **OpenAI ChatGPT-5** - AI-powered anime recommendations
