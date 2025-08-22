/**
 * GraphQL query to get anime details by ID including next airing episode
 */
export const GET_ANIME_DETAILS = `
  query ($id: Int!) {
    Media(id: $id, type: ANIME) {
      id
      title { romaji english native }
      status
      format
      episodes
      nextAiringEpisode {
        episode
        airingAt
        timeUntilAiring
      }
      coverImage { large }
      siteUrl
    }
  }
`
