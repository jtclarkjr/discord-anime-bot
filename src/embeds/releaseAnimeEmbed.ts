import { EmbedBuilder } from 'discord.js'

export function createReleaseAnimeEmbed(
  animeList: string,
  count: number,
  total: number
): EmbedBuilder {
  return new EmbedBuilder()
    .setTitle('Currently Releasing Anime')
    .setDescription(animeList)
    .setColor(0x02a9ff)
    .setFooter({ text: `Showing ${count} of ${total} releasing anime` })
}
