package graphql

// SearchAnimeByTextQuery is the GraphQL query for searching anime by text
const SearchAnimeByTextQuery = `
	query ($search: String, $page: Int, $perPage: Int) {
		Page(page: $page, perPage: $perPage) {
			pageInfo {
				total
				currentPage
				lastPage
				hasNextPage
			}
			media(search: $search, type: ANIME) {
				id
				title {
					romaji
					english
					native
				}
				format
				status
				coverImage {
					large
				}
				siteUrl
			}
		}
	}`
