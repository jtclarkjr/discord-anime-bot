import { ANILIST_API } from "@config/constants"

/**
 * Helper function to make GraphQL requests to AniList API
 */
export const makeAniListRequest = async (query: string, variables: any) => {
  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })
  
  if (!res.ok) {
    throw new Error(`AniList API request failed with status: ${res.status}`)
  }
  
  return res.json()
}
