package types

// PageInfo represents pagination information from AniList API
type PageInfo struct {
	Total       int  `json:"total"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
}

// AnimeTitle represents different title formats for anime
type AnimeTitle struct {
	Romaji  string  `json:"romaji"`
	English *string `json:"english"`
	Native  string  `json:"native"`
}

// CoverImage represents anime cover image URLs
type CoverImage struct {
	Large string `json:"large"`
}

// EpisodeInfo represents basic episode and airing information
type EpisodeInfo struct {
	Episode  int `json:"episode"`
	AiringAt int `json:"airingAt"`
}

// NextAiringEpisode represents information about the next airing episode
type NextAiringEpisode struct {
	EpisodeInfo
	TimeUntilAiring int `json:"timeUntilAiring,omitempty"`
}

// AnimeMedia represents basic anime information
type AnimeMedia struct {
	ID         int        `json:"id"`
	Title      AnimeTitle `json:"title"`
	Format     string     `json:"format"`
	Status     string     `json:"status"`
	CoverImage CoverImage `json:"coverImage"`
	SiteURL    string     `json:"siteUrl"`
}

// AnimeDetails represents detailed anime information
type AnimeDetails struct {
	ID                int                `json:"id"`
	Title             AnimeTitle         `json:"title"`
	Status            string             `json:"status"`
	Format            string             `json:"format"`
	Episodes          *int               `json:"episodes"`
	NextAiringEpisode *NextAiringEpisode `json:"nextAiringEpisode"`
	CoverImage        CoverImage         `json:"coverImage"`
	SiteURL           string             `json:"siteUrl"`
}

// ReleasingAnime represents anime that is currently releasing
type ReleasingAnime struct {
	ID                int                `json:"id"`
	Title             AnimeTitle         `json:"title"`
	NextAiringEpisode *NextAiringEpisode `json:"nextAiringEpisode"`
}

// Generic response types
type AniListPageResponse[T any] struct {
	Data struct {
		Page struct {
			PageInfo PageInfo `json:"pageInfo"`
			Media    []T      `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

type AniListSingleResponse[T any] struct {
	Data struct {
		Media T `json:"Media"`
	} `json:"data"`
}

// SearchResponse represents the response from AniList search API
type SearchResponse = AniListPageResponse[AnimeMedia]

// AnimeDetailsResponse represents the response from AniList anime details API
type AnimeDetailsResponse = AniListSingleResponse[AnimeDetails]

// ReleasingAnimeResponse represents the response from AniList releasing anime API
type ReleasingAnimeResponse = AniListPageResponse[ReleasingAnime]

// AnimeMatch represents a match found by AI with confidence and reasoning
type AnimeMatch struct {
	Anime      AnimeMedia `json:"anime"`
	Reason     string     `json:"reason"`
	Confidence float64    `json:"confidence"`
}

// OpenAIRecommendation represents a single anime recommendation from OpenAI
type OpenAIRecommendation struct {
	Title      string  `json:"title"`
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
}

// GraphQLRequest represents a generic GraphQL request structure
type GraphQLRequest[T any] struct {
	Query     string `json:"query"`
	Variables T      `json:"variables"`
}

// GraphQLSearchVariables represents variables for GraphQL search query
type GraphQLSearchVariables struct {
	Search  string `json:"search"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

// GraphQLNextVariables represents variables for GraphQL next episode query
type GraphQLNextVariables struct {
	ID int `json:"id"`
}

// NotificationEntry represents a notification entry with timer
type NotificationEntry struct {
	AnimeID   int    `json:"animeId"`
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
	Episode   int    `json:"episode"`
	AiringAt  int64  `json:"airingAt"`
}

// PersistedNotification represents a notification entry for storage (same as NotificationEntry)
type PersistedNotification = NotificationEntry
