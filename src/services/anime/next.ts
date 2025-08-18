import type { AnimeDetailsResponse, AnimeDetails } from '@/types/anilist'
import { makeAniListRequest } from '@/utils/request'

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

  const json = (await makeAniListRequest(query, variables)) as AnimeDetailsResponse
  return json.data.Media
}
