import { EmbedBuilder } from 'discord.js'
import type { AnimeMedia } from '@/types/anilist'

export function createSearchAnimeEmbed(anime: AnimeMedia): EmbedBuilder {
  return new EmbedBuilder()
    .setTitle(anime.title.english || anime.title.romaji)
    .setURL(anime.siteUrl)
    .setThumbnail(anime.coverImage.large)
    .addFields(
      { name: 'Romaji Title', value: anime.title.romaji, inline: true },
      { name: 'Native Title', value: anime.title.native, inline: true },
      { name: 'Format', value: anime.format, inline: true },
      { name: 'Status', value: anime.status, inline: true },
      { name: 'AniList ID', value: anime.id.toString(), inline: true }
    )
    .setColor(0x02a9ff)
}
