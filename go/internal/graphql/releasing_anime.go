package graphql

// GetReleasingAnimeQuery is the GraphQL query for getting currently releasing anime
const GetReleasingAnimeQuery = `
	query ($page: Int, $perPage: Int) {
		Page(page: $page, perPage: $perPage) {
			media(type: ANIME, status: RELEASING, sort: [POPULARITY_DESC]) {
				id
				title { 
					romaji 
					english 
				}
				nextAiringEpisode {
					episode
					airingAt
				}
			}
			pageInfo { 
				total 
				currentPage 
				lastPage 
				hasNextPage 
			}
		}
	}`
