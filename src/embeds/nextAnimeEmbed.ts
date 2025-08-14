import { EmbedBuilder } from 'discord.js'
import type { AnimeDetails } from '@/types/anilist'

export function createNextAnimeEmbed(anime: AnimeDetails, options: {
  isFinished: boolean
  nextAiringEpisode?: {
    airingAt: number
    timeUntilAiring?: number
    episode: number
    formattedAirDate: string
    timeString: string
  }
}): EmbedBuilder {
  const embed = new EmbedBuilder()
    .setTitle(anime.title.english || anime.title.romaji)
    .setURL(anime.siteUrl)
    .setThumbnail(anime.coverImage.large)
    .setColor(0x02A9FF)

  if (options.isFinished) {
    embed.setDescription('This anime has finished airing.')
    embed.addFields(
      { name: 'Status', value: anime.status, inline: true },
      { name: 'Format', value: anime.format, inline: true },
      { name: 'Total Episodes', value: anime.episodes?.toString() || 'Unknown', inline: true }
    )
  } else if (options.nextAiringEpisode) {
    embed.setDescription(`Next episode airs in ${options.nextAiringEpisode.timeString} (${options.nextAiringEpisode.formattedAirDate})`)
    embed.addFields(
      { name: 'Status', value: anime.status, inline: true },
      { name: 'Format', value: anime.format, inline: true },
      { name: 'Next Episode', value: options.nextAiringEpisode.episode.toString(), inline: true }
    )
  }

  return embed
}
