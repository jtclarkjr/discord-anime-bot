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

export type AnimeDetails = {
  id: number
  title: AnimeTitle
  status: string
  format: string
  episodes: number | null
  nextAiringEpisode: NextAiringEpisode | null
  coverImage: CoverImage
  siteUrl: string
}

export type ReleasingAnime = {
  id: number
  title: Pick<AnimeTitle, 'romaji' | 'english'>
  nextAiringEpisode: Omit<NextAiringEpisode, 'timeUntilAiring'> | null
}

// Response types
export type AnimeSearchResponse = {
  data: {
    Page: {
      media: AnimeMedia[]
      pageInfo: PageInfo
    }
  }
}

export type AnimeDetailsResponse = {
  data: {
    Media: AnimeDetails
  }
}

export type ReleasingAnimeResponse = {
  data: {
    Page: {
      media: ReleasingAnime[]
      pageInfo: PageInfo
    }
  }
}

// Utility types
// OpenAI types
export interface AnimeMatch {
  anime: AnimeMedia
  reason: string
  confidence: number
}
