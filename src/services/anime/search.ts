import type { AnimeSearchResponse } from '@/types/anilist'
import { ANILIST_API } from '@/config/constants'

/**
 * Search for anime using AniList API
 */
export async function searchAnime(searchQuery: string, page: number = 1, perPage: number = 10) {
  const query = `
    query ($q: String!, $page: Int = 1, $perPage: Int = 10) {
      Page(page: $page, perPage: $perPage) {
        media(search: $q, type: ANIME, sort: [SEARCH_MATCH, POPULARITY_DESC]) {
          id
          title { romaji english native }
          format
          status
          coverImage { large }
          siteUrl
        }
        pageInfo { total currentPage lastPage hasNextPage }
      }
    }
  `
  
  const variables = {
    q: searchQuery,
    page,
    perPage
  }

  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })

  const json = (await res.json()) as AnimeSearchResponse
  return json.data.Page
}
