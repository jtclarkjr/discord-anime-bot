import type { SeasonAnimeResponse } from '@/types/anilist'
import { ANILIST_API } from '@/config/constants'

/**
 * Get anime for a specific season and year
 */
export async function getSeasonAnime(season: string, seasonYear: number, page: number = 1, perPage: number = 50) {
  const query = `
    query SeasonAnime($season: MediaSeason, $seasonYear: Int, $type: MediaType, $page: Int, $perPage: Int) {
      Page(page: $page, perPage: $perPage) {
        media(season: $season, seasonYear: $seasonYear, type: $type, sort: [POPULARITY_DESC]) {
          id
          title { romaji english }
          coverImage { medium large }
          status
        }
        pageInfo { total currentPage lastPage hasNextPage }
      }
    }
  `
  
  const variables = {
    season: season.toUpperCase(),
    seasonYear,
    type: 'ANIME',
    page,
    perPage
  }

  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })

  if (!res.ok) {
    throw new Error(`AniList API error: ${res.status}`)
  }

  const json = (await res.json()) as SeasonAnimeResponse
  
  if (json.errors) {
    throw new Error(`AniList GraphQL error: ${json.errors[0]?.message}`)
  }

  return json.data.Page
}
