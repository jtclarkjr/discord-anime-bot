import { EmbedBuilder } from 'discord.js'
import { getAnimeById } from '@/services/anime/next'

export async function createWatchlistEmbed(watchlist: number[]) {
  let description = ''
  for (const animeId of watchlist) {
    const anime = await getAnimeById(animeId)
    if (anime) {
      const title = anime.title.english || anime.title.romaji
      description += `• **${title}** (ID: ${animeId})\n`
    }
  }
  return new EmbedBuilder()
    .setTitle('Your Anime Watchlist')
    .setDescription(description || 'No anime in your watchlist.')
    .setColor(0x0099ff)
}
