import type { SeasonAnimeResponse } from '@/types/anilist'
import { makeAniListRequest } from '@/utils/request'
import { GET_SEASONAL_ANIME } from '@/graphql'

/**
 * Get anime for a specific season and year
 */
export async function getSeasonAnime(
  season: string,
  seasonYear: number,
  page: number = 1,
  perPage: number = 50
) {
  const variables = {
    season: season.toUpperCase(),
    seasonYear,
    type: 'ANIME',
    page,
    perPage
  }

  const json = (await makeAniListRequest(GET_SEASONAL_ANIME, variables)) as SeasonAnimeResponse

  if (json.errors) {
    const [error] = json.errors
    throw new Error(`AniList GraphQL error: ${error?.message}`)
  }

  return json.data.Page
}
