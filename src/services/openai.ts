import { OpenAI } from 'openai'
import { OPENAI_API_KEY } from '../config/constants'

const openai = new OpenAI({
  apiKey: OPENAI_API_KEY
})

export interface AnimeRecommendation {
  title: string
  reason: string
  confidence: number
}

/**
 * Find anime recommendations based on a description using ChatGPT-5
 */
export async function findAnimeByDescription(description: string): Promise<AnimeRecommendation[]> {
  try {
    const prompt = `Based on this description, recommend 3 anime titles that match:

Description: "${description}"

Please respond with ONLY a valid JSON array in this exact format:
[
  {
    "title": "Exact anime title (English or most common name)",
    "reason": "Brief explanation why this matches the description",
    "confidence": 0.95
  }
]

Focus on popular, well-known anime. Use confidence scores between 0.1-1.0 based on how well the anime matches the description. Return only the JSON array, no other text.`

    const completion = await openai.chat.completions.create({
      model: "gpt-4o",
      messages: [
        {
          role: "system",
          content: "You are an anime expert. You help people find anime based on descriptions. Always respond with valid JSON only."
        },
        {
          role: "user",
          content: prompt
        }
      ],
      temperature: 0.3,
      max_tokens: 500
    })

    const response = completion.choices[0]?.message?.content
    if (!response) {
      throw new Error('No response from OpenAI')
    }

    // Parse the JSON response
    const recommendations = JSON.parse(response) as AnimeRecommendation[]
    
    // Validate the response structure
    if (!Array.isArray(recommendations)) {
      throw new Error('Invalid response format')
    }

    return recommendations.filter(rec => 
      rec.title && rec.reason && typeof rec.confidence === 'number'
    )

  } catch (error) {
    console.error('Error finding anime by description:', error)
    throw new Error('Failed to find anime recommendations')
  }
}
