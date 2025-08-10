import { searchAnime } from './search'
import { findAnimeByDescription } from '@/services/openai/completions'
import type { AnimeMatch } from '@/types/anilist'

/**
 * Find anime using AI description and return AniList details
 */
export async function findAnimeWithDetails(description: string): Promise<AnimeMatch[]> {
  try {
    // Get AI recommendations
    const recommendations = await findAnimeByDescription(description)
    
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
