// Base types
export type PageInfo = {
  total: number
  currentPage: number
  lastPage: number
  hasNextPage: boolean
}

export type AnimeTitle = {
  romaji: string
  english: string | null
  native: string
}

export type CoverImage = {
  large: string
}

export type NextAiringEpisode = {
  episode: number
  airingAt: number
  timeUntilAiring?: number
}

// Media types
export type AnimeMedia = {
  id: number
  title: AnimeTitle
  format: string
  status: string
  coverImage: CoverImage
  siteUrl: string
}

export type AnimeDetails = AnimeMedia & {
  episodes: number | null
  nextAiringEpisode: NextAiringEpisode | null
}

export type ReleasingAnime = {
  id: number
  title: Pick<AnimeTitle, 'romaji' | 'english'>
  nextAiringEpisode: Omit<NextAiringEpisode, 'timeUntilAiring'> | null
}

// Response types
export type AniListPageResponse<T> = {
  data: {
    Page: {
      media: T[]
      pageInfo: PageInfo
    }
  }
}

export type AniListSingleResponse<T> = {
  data: {
    Media: T
  }
}

export type AnimeSearchResponse = AniListPageResponse<AnimeMedia>

export type AnimeDetailsResponse = AniListSingleResponse<AnimeDetails>

export type ReleasingAnimeResponse = AniListPageResponse<ReleasingAnime>

// Utility types
export type NotificationEntry = Pick<NextAiringEpisode, 'episode' | 'airingAt'> & {
  animeId: number
  channelId: string
  userId: string
  timeoutId?: NodeJS.Timeout
}

// OpenAI types
export type AnimeMatch = {
  anime: AnimeMedia
  reason: string
  confidence: number
}
