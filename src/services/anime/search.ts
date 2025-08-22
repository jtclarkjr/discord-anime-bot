import type { AnimeSearchResponse } from '@/types/anilist'
import { makeAniListRequest } from '@/utils/request'
import { SEARCH_ANIME_BY_ID, SEARCH_ANIME_BY_TEXT } from '@/graphql'

/**
 * Search for anime using AniList API
 * Supports both text search and ID lookup
 * @param searchQuery - Either a text query to search for or a numeric anime ID
 * @param page - Page number for text search results (ignored for ID search)
 * @param perPage - Number of results per page for text search (ignored for ID search)
 * @returns Page containing matching anime with pagination info
 */
export async function searchAnime(searchQuery: string, page: number = 1, perPage: number = 10) {
  // Check if the search query is a numeric ID
  const numericId = parseInt(searchQuery.trim())
  const isNumericSearch = !isNaN(numericId) && numericId.toString() === searchQuery.trim()

  if (isNumericSearch) {
    // Search by ID - return single result in page format
    return searchAnimeById(numericId)
  } else {
    // Text search
    return searchAnimeByText(searchQuery, page, perPage)
  }
}

/**
 * Search for anime by ID and return it in page format
 */
async function searchAnimeById(animeId: number) {
  const variables = { id: animeId }
  const json = await makeAniListRequest(SEARCH_ANIME_BY_ID, variables)

  // Return in the same format as search results
  if (json.data?.Media) {
    return {
      media: [json.data.Media],
      pageInfo: {
        total: 1,
        currentPage: 1,
        lastPage: 1,
        hasNextPage: false
      }
    }
  } else {
    // No anime found with that ID
    return {
      media: [],
      pageInfo: {
        total: 0,
        currentPage: 1,
        lastPage: 1,
        hasNextPage: false
      }
    }
  }
}

/**
 * Search for anime by text query
 */
async function searchAnimeByText(searchQuery: string, page: number, perPage: number) {
  const variables = {
    q: searchQuery,
    page,
    perPage
  }

  const json = (await makeAniListRequest(SEARCH_ANIME_BY_TEXT, variables)) as AnimeSearchResponse
  return json.data.Page
}
