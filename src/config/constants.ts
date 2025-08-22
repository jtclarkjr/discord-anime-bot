export const DISCORD_TOKEN = process.env.DISCORD_BOT_TOKEN!
export const ANILIST_API = process.env.ANILIST_API!
export const OPENAI_API_KEY = process.env.OPENAI_API_KEY
export const CLAUDE_API_KEY = process.env.CLAUDE_API_KEY

const CHANNEL_ID = process.env.CHANNEL_ID!

// Validate that required environment variables are set
if (!DISCORD_TOKEN) {
  console.error('❌ DISCORD_BOT_TOKEN is not set in environment variables.')
  process.exit(1)
}
if (!CHANNEL_ID) {
  console.error('❌ CHANNEL_ID is not set in environment variables.')
  process.exit(1)
}
if (!ANILIST_API) {
  console.error('❌ ANILIST_API is not set in environment variables.')
  process.exit(1)
}

// AI is enabled if either OpenAI or Claude key is present
if (!OPENAI_API_KEY && !CLAUDE_API_KEY) {
  console.warn('⚠️ No AI API key is set. AI-powered features (like /anime find) will be disabled.')
}

export const IS_OPENAI_ENABLED = !!OPENAI_API_KEY
export const IS_CLAUDE_ENABLED = !!CLAUDE_API_KEY
export const IS_AI_ENABLED = !!OPENAI_API_KEY || !!CLAUDE_API_KEY

export const storageFile = process.env.STORAGE_FILE!
export const watchlistFile = process.env.WATCHLIST_FILE!
