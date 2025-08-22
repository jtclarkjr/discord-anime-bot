/**
 * GraphQL query to search for anime by ID
 */
export const SEARCH_ANIME_BY_ID = `
  query ($id: Int!) {
    Media(id: $id, type: ANIME) {
      id
      title { romaji english native }
      format
      status
      coverImage { large }
      siteUrl
    }
  }
`
