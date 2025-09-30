import { ChatInputCommandInteraction } from 'discord.js'
import { createNextAnimeEmbed } from '@/embeds/nextAnimeEmbed'
import { getAnimeById } from '@/services/anime/next'
import { formatCountdown, formatAirDate } from '@/utils/formatters'

export async function handleNextCommand(interaction: ChatInputCommandInteraction) {
  const animeId = interaction.options.getInteger('id')
  if (!animeId) {
    await interaction.reply('Please provide a valid anime ID.')
    return
  }

  await interaction.deferReply()

  try {
    const anime = await getAnimeById(animeId)

    if (!anime) {
      await interaction.editReply(`No anime found with ID ${animeId}.`)
      return
    }

    const isFinished = anime.status === 'FINISHED' || anime.status === 'CANCELLED'
    let nextAiringEpisode
    if (anime.nextAiringEpisode && !isFinished) {
      const airingDate = new Date(anime.nextAiringEpisode.airingAt * 1000)
      nextAiringEpisode = {
        airingAt: anime.nextAiringEpisode.airingAt,
        timeUntilAiring: anime.nextAiringEpisode.timeUntilAiring,
        episode: anime.nextAiringEpisode.episode,
        formattedAirDate: formatAirDate(airingDate),
        timeString: formatCountdown(anime.nextAiringEpisode.timeUntilAiring || 0)
      }
    }

    const embed = createNextAnimeEmbed(anime, { isFinished, nextAiringEpisode })

    await interaction.editReply({ embeds: [embed] })
  } catch (error) {
    console.error('Error in next command:', error)
    await interaction.editReply('An error occurred while fetching anime data.')
  }
}
