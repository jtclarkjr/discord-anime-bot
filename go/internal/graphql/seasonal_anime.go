package graphql

// GetSeasonalAnimeQuery is the GraphQL query for getting anime from a specific season and year
const GetSeasonalAnimeQuery = `
	query SeasonAnime($season: MediaSeason, $seasonYear: Int, $type: MediaType, $page: Int, $perPage: Int) {
		Page(page: $page, perPage: $perPage) {
			media(season: $season, seasonYear: $seasonYear, type: $type, sort: [POPULARITY_DESC]) {
				id
				title { 
					romaji 
					english 
				}
				coverImage {
					medium
					large
				}
				status
			}
			pageInfo { 
				total 
				currentPage 
				lastPage 
				hasNextPage 
			}
		}
	}`
