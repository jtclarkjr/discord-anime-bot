/**
 * GraphQL query to get currently releasing anime
 */
export const GET_RELEASING_ANIME = `
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
