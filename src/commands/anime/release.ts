import { ChatInputCommandInteraction } from 'discord.js'
import { createReleaseAnimeEmbed } from '@/embeds/releaseAnimeEmbed'
import { getReleasingAnime } from '@/services/anime/release'
import { formatCompactDateTime } from '@/utils/formatters'

export async function handleReleaseCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply()

  try {
    // Get pagination options from the interaction
    const page = interaction.options.getInteger('page') || 1
    const perPage = interaction.options.getInteger('perpage') || 15
    const releasingAnime = await getReleasingAnime(page, perPage)

    if (releasingAnime.media.length === 0) {
      await interaction.editReply('No releasing anime found.')
      return
    }

    // Create a list of currently releasing anime
    const animeList = releasingAnime.media
      .map((anime) => {
        const title = anime.title.english || anime.title.romaji
        let nextEpisodeInfo = ''

        if (anime.nextAiringEpisode) {
          const airingDate = new Date(anime.nextAiringEpisode.airingAt * 1000)
          const formattedTime = formatCompactDateTime(airingDate)
          nextEpisodeInfo = ` - Ep ${anime.nextAiringEpisode.episode} on ${formattedTime}`
        } else {
          nextEpisodeInfo = ' - No schedule'
        }

        return `**${title}** (ID: ${anime.id})${nextEpisodeInfo}`
      })
      .join('\n')

    const embed = createReleaseAnimeEmbed(
      animeList,
      releasingAnime.media.length,
      releasingAnime.pageInfo.total
    )

    await interaction.editReply({ embeds: [embed] })
  } catch (error) {
    console.error('Error in release command:', error)
    await interaction.editReply('An error occurred while fetching releasing anime.')
  }
}
