import { ChatInputCommandInteraction, EmbedBuilder } from 'discord.js'
import { getReleasingAnime } from '@/services/anime/release'
import { formatCompactDateTime } from '@/utils/formatters'

export async function handleReleaseCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply()

  try {
    const releasingAnime = await getReleasingAnime()
    
    if (releasingAnime.media.length === 0) {
      await interaction.editReply('❌ No releasing anime found.')
      return
    }

    // Create a list of currently releasing anime
    const animeList = releasingAnime.media.slice(0, 15).map(anime => {
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
    }).join('\n')

    const embed = new EmbedBuilder()
      .setTitle('Currently Releasing Anime')
      .setDescription(animeList)
      .setColor(0x02A9FF)
      .setFooter({ text: `Showing ${Math.min(15, releasingAnime.media.length)} of ${releasingAnime.pageInfo.total} releasing anime` })

    await interaction.editReply({ embeds: [embed] })
  } catch (error) {
    console.error('Error in release command:', error)
    await interaction.editReply('❌ An error occurred while fetching releasing anime.')
  }
}
