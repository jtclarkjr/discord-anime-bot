import type { ReleasingAnimeResponse } from '@/types/anilist'
import { ANILIST_API } from '@/config/constants'

/**
 * Get all currently releasing anime
 */
export async function getReleasingAnime(page: number = 1, perPage: number = 25) {
  const query = `
    query ($page: Int, $perPage: Int) {
      Page(page: $page, perPage: $perPage) {
        media(type: ANIME, status: RELEASING, sort: [POPULARITY_DESC]) {
          id
          title { romaji english }
          nextAiringEpisode {
            episode
            airingAt
          }
        }
        pageInfo { total currentPage lastPage hasNextPage }
      }
    }
  `
  
  const variables = {
    page,
    perPage
  }

  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })

  const json = (await res.json()) as ReleasingAnimeResponse
  return json.data.Page
}
