/**
 * GraphQL query to search for anime by text
 */
export const SEARCH_ANIME_BY_TEXT = `
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
