import { EmbedBuilder } from 'discord.js'
import type { AnimeMedia, AnimeMatch } from '@/types/anilist'

export function createFindAnimeEmbed(anime: AnimeMedia, bestMatch: AnimeMatch): EmbedBuilder {
  return new EmbedBuilder()
    .setTitle(`🎯 ${anime.title.english || anime.title.romaji}`)
    .setURL(anime.siteUrl)
    .setThumbnail(anime.coverImage.large)
    .setDescription(bestMatch.reason)
    .addFields(
      { name: 'Romaji Title', value: anime.title.romaji, inline: true },
      { name: 'Native Title', value: anime.title.native, inline: true },
      { name: 'Format', value: anime.format, inline: true },
      { name: 'Status', value: anime.status, inline: true },
      { name: 'AniList ID', value: anime.id.toString(), inline: true },
      { name: 'AI Confidence', value: `${Math.round(bestMatch.confidence * 100)}%`, inline: true }
    )
    .setColor(0x00ff00)
    .setFooter({ text: 'Powered by GPT-5 + AniList' })
}
