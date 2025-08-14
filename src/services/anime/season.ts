import type { SeasonAnimeResponse } from '@/types/anilist'
import { makeAniListRequest } from '@/utils/request'

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

  const json = (await makeAniListRequest(query, variables)) as SeasonAnimeResponse
  
  if (json.errors) {
    const [error] = json.errors
    throw new Error(`AniList GraphQL error: ${error?.message}`)
  }

  return json.data.Page
}
