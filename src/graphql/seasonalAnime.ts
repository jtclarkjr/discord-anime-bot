/**
 * GraphQL query to get anime for a specific season and year
 */
export const GET_SEASONAL_ANIME = `
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
