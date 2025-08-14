import { searchAnime } from './search'
import { findAnimeByDescription as findAnimeByDescriptionOpenAI } from '@/services/openai/completions'
import { findAnimeByDescriptionClaude } from '@/services/claude/completions'
import { IS_OPENAI_ENABLED, IS_CLAUDE_ENABLED, OPENAI_API_KEY, CLAUDE_API_KEY } from '@/config/constants'
import type { AnimeRecommendation } from '@/types/openai'
import type { AnimeMatch } from '@/types/anilist'

/**
 * Find anime using AI description and return AniList details
 */
export async function findAnimeWithDetails(description: string): Promise<AnimeMatch[]> {
  if (!IS_OPENAI_ENABLED && !IS_CLAUDE_ENABLED) {
    throw new Error('AI is not configured. Please set OPENAI_API_KEY or CLAUDE_API_KEY environment variable to use AI-powered anime search.')
  }

  try {
    // Prefer OpenAI if both keys are present
  let recommendations: AnimeRecommendation[] = []
    if (IS_OPENAI_ENABLED && OPENAI_API_KEY) {
      recommendations = await findAnimeByDescriptionOpenAI(description)
    } else if (IS_CLAUDE_ENABLED && CLAUDE_API_KEY) {
      recommendations = await findAnimeByDescriptionClaude(description)
    }

    const matches: AnimeMatch[] = []

    // Search for each recommendation on AniList
    for (const rec of recommendations) {
      try {
        const searchResults = await searchAnime(rec.title, 1, 5)
        if (searchResults.media.length > 0) {
          // Use the first (most relevant) result
          const anime = searchResults.media[0]
          matches.push({
            anime,
            reason: rec.reason,
            confidence: rec.confidence
          })
        }
      } catch (error) {
        console.warn(`Could not find anime "${rec.title}" on AniList:`, error)
        // Continue with other recommendations
      }
    }

    // Sort by confidence score
    return matches.sort((a, b) => b.confidence - a.confidence)

  } catch (error) {
    console.error('Error finding anime with details:', error)
    throw new Error('Failed to find anime based on description')
  }
}
