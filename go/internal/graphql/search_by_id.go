package graphql

// SearchAnimeByIDQuery is the GraphQL query for searching anime by ID
const SearchAnimeByIDQuery = `
	query ($id: Int!) {
		Media(id: $id, type: ANIME) {
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
	}`
