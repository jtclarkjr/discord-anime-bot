import type { AnimeDetailsResponse, AnimeDetails } from '@/types/anilist'
import { makeAniListRequest } from '@/utils/request'
import { GET_ANIME_DETAILS } from '@/graphql'

/**
 * Get anime details by ID including next airing episode
 */
export async function getAnimeById(animeId: number): Promise<AnimeDetails | null> {
  const variables = {
    id: animeId
  }

  const json = (await makeAniListRequest(GET_ANIME_DETAILS, variables)) as AnimeDetailsResponse
  return json.data.Media
}
