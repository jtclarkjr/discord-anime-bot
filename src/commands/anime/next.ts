import { ChatInputCommandInteraction, EmbedBuilder } from 'discord.js'
import { getAnimeById } from '@/services/anime/next'
import { formatCountdown, formatAirDate } from '@/utils/formatters'

export async function handleNextCommand(interaction: ChatInputCommandInteraction) {
  const animeId = interaction.options.getInteger('id')
  if (!animeId) {
    await interaction.reply('❌ Please provide a valid anime ID.')
    return
  }

  await interaction.deferReply()

  try {
    const anime = await getAnimeById(animeId)
    
    if (!anime) {
      await interaction.editReply(`❌ No anime found with ID ${animeId}.`)
      return
    }

    const embed = new EmbedBuilder()
      .setTitle(anime.title.english || anime.title.romaji)
      .setURL(anime.siteUrl)
      .setThumbnail(anime.coverImage.large)
      .setColor(0x02A9FF)

    // Check if anime is finished or has next episode
    if (anime.status === 'FINISHED' || anime.status === 'CANCELLED') {
      embed.setDescription(`This anime has finished airing.`)
      embed.addFields(
        { name: 'Status', value: anime.status, inline: true },
        { name: 'Format', value: anime.format, inline: true },
        { name: 'Total Episodes', value: anime.episodes?.toString() || 'Unknown', inline: true }
      )
    } else if (anime.nextAiringEpisode) {
      const airingDate = new Date(anime.nextAiringEpisode.airingAt * 1000)
      const timeString = formatCountdown(anime.nextAiringEpisode.timeUntilAiring || 0)
      const formattedAirDate = formatAirDate(airingDate)
      
      embed.setDescription(`Episode ${anime.nextAiringEpisode.episode} airs in ${timeString}`)
      embed.addFields(
        { name: 'Next Episode', value: anime.nextAiringEpisode.episode.toString(), inline: true },
        { name: 'Air Date', value: formattedAirDate, inline: false },
        { name: 'Status', value: anime.status, inline: true }
      )
    } else {
      embed.setDescription(`No upcoming episodes scheduled.`)
      embed.addFields(
        { name: 'Status', value: anime.status, inline: true },
        { name: 'Format', value: anime.format, inline: true }
      )
    }

    await interaction.editReply({ embeds: [embed] })
  } catch (error) {
    console.error('Error in next command:', error)
    await interaction.editReply('❌ An error occurred while fetching anime data.')
  }
}
