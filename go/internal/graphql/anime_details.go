package graphql

// GetAnimeDetailsQuery is the GraphQL query for getting anime details by ID including next airing episode
const GetAnimeDetailsQuery = `
	query ($id: Int!) {
		Media(id: $id, type: ANIME) {
			id
			title { 
				romaji 
				english 
				native 
			}
			status
			format
			episodes
			nextAiringEpisode {
				episode
				airingAt
				timeUntilAiring
			}
			coverImage { 
				large 
			}
			siteUrl
		}
	}`
