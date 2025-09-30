import Anthropic from '@anthropic-ai/sdk'
import type { AnimeRecommendation } from '@/types/openai'
import { CLAUDE_API_KEY } from '@/config/constants'

let anthropic: Anthropic | null = null

if (CLAUDE_API_KEY) {
  anthropic = new Anthropic({
    apiKey: CLAUDE_API_KEY
  })
}

export async function findAnimeByDescriptionClaude(prompt: string): Promise<AnimeRecommendation[]> {
  if (!anthropic) {
    throw new Error('Claude is not configured. Please set CLAUDE_API_KEY environment variable.')
  }

  const systemPrompt = `You are an anime expert. Given a description, recommend 3 anime titles that match. Respond ONLY with a valid JSON array in this format:
   [
     {"title":"Anime Name","reason":"Why it matches","confidence":0.9}
   ]`

  const msg = await anthropic.messages.create({
    model: 'claude-sonnet-4-5',
    max_tokens: 2048, // Increased token limit for longer responses if needed
    messages: [{ role: 'user', content: `${systemPrompt}\n\n${prompt}` }]
  })

  try {
    const contentText = Array.isArray(msg.content)
      ? msg.content
          .map((block: unknown) => {
            if (
              typeof block === 'object' &&
              block !== null &&
              'text' in block &&
              typeof (block as { text?: string }).text === 'string'
            ) {
              return (block as { text: string }).text
            }
            return ''
          })
          .join('')
      : msg.content
    return JSON.parse(contentText) as AnimeRecommendation[]
  } catch {
    return []
  }
}
