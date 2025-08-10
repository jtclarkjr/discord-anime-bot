export const DISCORD_TOKEN = process.env.DISCORD_BOT_TOKEN!
export const ANILIST_API = process.env.ANILIST_API!
export const OPENAI_API_KEY = process.env.OPENAI_API_KEY!

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
if (!OPENAI_API_KEY) {
  console.error('❌ OPENAI_API_KEY is not set in environment variables.')
  process.exit(1)
}
