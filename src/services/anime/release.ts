import type { ReleasingAnimeResponse } from '@/types/anilist'
import { makeAniListRequest } from '@/utils/request'
import { GET_RELEASING_ANIME } from '@/graphql'

/**
 * Get all currently releasing anime
 */
export async function getReleasingAnime(page: number, perPage: number) {
  const variables = {
    page,
    perPage
  }

  const json = (await makeAniListRequest(GET_RELEASING_ANIME, variables)) as ReleasingAnimeResponse
  return json.data.Page
}
