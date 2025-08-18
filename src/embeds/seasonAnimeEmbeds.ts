import { EmbedBuilder } from 'discord.js'

export function createSeasonAnimeEmbeds({
  season,
  year,
  animeMedia,
  getStatusEmoji
}: {
  season: string
  year: number
  animeMedia: import('@/types/anilist').SeasonAnime[]
  getStatusEmoji: (status: import('@/types/anilist').AnimeStatus) => string
}): EmbedBuilder[] {
  const animePerEmbed = 20
  const totalEmbeds = Math.ceil(animeMedia.length / animePerEmbed)

  return Array.from({ length: totalEmbeds }, (_, i) => {
    const startIndex = i * animePerEmbed
    const endIndex = Math.min(startIndex + animePerEmbed, animeMedia.length)
    const animeSlice = animeMedia.slice(startIndex, endIndex)

    const animeList = animeSlice
      .map((anime, index) => {
        const title = anime.title.english || anime.title.romaji
        const statusEmoji = getStatusEmoji(anime.status)
        return `${startIndex + index + 1}. **${title}** ${statusEmoji} (ID: ${anime.id})`
      })
      .join('\n')

    const embed = new EmbedBuilder()
      .setTitle(
        `${season.charAt(0).toUpperCase() + season.slice(1)} ${year} Anime${totalEmbeds > 1 ? ` (Part ${i + 1}/${totalEmbeds})` : ''}`
      )
      .setDescription(animeList)
      .setColor(0x02a9ff)

    if (i === 0) {
      embed.setFooter({ text: `Showing all ${animeMedia.length} anime from ${season} ${year}` })
    }

    return embed
  })
}
