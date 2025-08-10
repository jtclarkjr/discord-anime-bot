import type { AnimeDetailsResponse, AnimeDetails } from '@/types/anilist'
import { ANILIST_API } from '@/config/constants'

/**
 * Get anime details by ID including next airing episode
 */
export async function getAnimeById(animeId: number): Promise<AnimeDetails | null> {
  const query = `
    query ($id: Int!) {
      Media(id: $id, type: ANIME) {
        id
        title { romaji english native }
        status
        format
        episodes
        nextAiringEpisode {
          episode
          airingAt
          timeUntilAiring
        }
        coverImage { large }
        siteUrl
      }
    }
  `
  
  const variables = {
    id: animeId
  }

  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })

  const json = (await res.json()) as AnimeDetailsResponse
  return json.data.Media
}
