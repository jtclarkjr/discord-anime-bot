import OpenAI from 'openai'
import { OPENAI_API_KEY, IS_OPENAI_ENABLED } from '@/config/constants'
import type { AnimeRecommendation } from '@/types/openai'

let openai: OpenAI | null = null

if (IS_OPENAI_ENABLED && OPENAI_API_KEY) {
  openai = new OpenAI({
    apiKey: OPENAI_API_KEY
  })
}

/**
 * Find anime recommendations based on a description using ChatGPT-5
 */
export async function findAnimeByDescription(description: string): Promise<AnimeRecommendation[]> {
  if (!openai || !IS_OPENAI_ENABLED) {
    throw new Error('OpenAI is not configured. Please set OPENAI_API_KEY environment variable.')
  }

  try {
    const prompt = `Find 3 anime that match: "${description}"

Respond with only valid JSON:
[
  {"title": "Anime Name", "reason": "Why it matches", "confidence": 0.9}
]`

    const completion = await openai.chat.completions.create({
      model: 'gpt-5',
      messages: [
        {
          role: 'user',
          content: prompt
        }
      ]
      //   max_completion_tokens: 1000
    })

    // console.log('OpenAI completion:', JSON.stringify(completion, null, 2))

    const response = completion.choices[0]?.message?.content
    if (!response) {
      console.error('No response content. Full completion:', completion)
      throw new Error('No response from OpenAI')
    }

    // console.log('OpenAI response:', response)

    // Parse the JSON response
    const recommendations = JSON.parse(response) as AnimeRecommendation[]

    // Validate the response structure
    if (!Array.isArray(recommendations)) {
      throw new Error('Invalid response format')
    }

    return recommendations.filter(
      (rec) => rec.title && rec.reason && typeof rec.confidence === 'number'
    )
  } catch (error) {
    console.error('Error finding anime by description:', error)
    throw new Error('Failed to find anime recommendations')
  }
}
