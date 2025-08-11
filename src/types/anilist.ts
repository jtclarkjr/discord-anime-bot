// Base types
export type AnimeStatus = 'RELEASING' | 'FINISHED' | 'NOT_YET_RELEASED' | 'CANCELLED' | 'HIATUS'

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
  medium: string
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

export type SeasonAnime = {
  id: number
  title: Pick<AnimeTitle, 'romaji' | 'english'>
  coverImage: Pick<CoverImage, 'medium' | 'large'>
  status: AnimeStatus
}

// Response types
export type AniListPageResponse<T> = {
  data: {
    Page: {
      media: T[]
      pageInfo: PageInfo
    }
  }
  errors?: Array<{ message: string }>
}

export type AniListSingleResponse<T> = {
  data: {
    Media: T
  }
}

export type AnimeSearchResponse = AniListPageResponse<AnimeMedia>

export type AnimeDetailsResponse = AniListSingleResponse<AnimeDetails>

export type ReleasingAnimeResponse = AniListPageResponse<ReleasingAnime>

export type SeasonAnimeResponse = AniListPageResponse<SeasonAnime>

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
