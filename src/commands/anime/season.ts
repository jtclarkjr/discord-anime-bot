import { ChatInputCommandInteraction, EmbedBuilder } from 'discord.js'
import { getSeasonAnime } from '@/services/anime/season'
import type { AnimeStatus } from '@/types/anilist'

export async function handleSeasonCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply()

  try {
    const season = interaction.options.getString('season', true)
    const year = interaction.options.getInteger('year') || new Date().getFullYear()

    // Validate season
    const validSeasons = ['winter', 'spring', 'summer', 'fall']
    if (!validSeasons.includes(season.toLowerCase())) {
      await interaction.editReply('❌ Invalid season. Please use: winter, spring, summer, or fall.')
      return
    }

    const seasonAnime = await getSeasonAnime(season, year)
    
    if (seasonAnime.media.length === 0) {
      await interaction.editReply(`❌ No anime found for ${season} ${year}.`)
      return
    }

    // Create multiple embeds if necessary (Discord has a 4096 character limit for description)
    const animePerEmbed = 20
    const totalEmbeds = Math.ceil(seasonAnime.media.length / animePerEmbed)
    const embeds: EmbedBuilder[] = []

    for (let i = 0; i < totalEmbeds; i++) {
      const startIndex = i * animePerEmbed
      const endIndex = Math.min(startIndex + animePerEmbed, seasonAnime.media.length)
      const animeSlice = seasonAnime.media.slice(startIndex, endIndex)

      const animeList = animeSlice.map((anime, index) => {
        const title = anime.title.english || anime.title.romaji
        const statusEmoji = getStatusEmoji(anime.status)
        return `${startIndex + index + 1}. **${title}** ${statusEmoji} (ID: ${anime.id})`
      }).join('\n')

      const embed = new EmbedBuilder()
        .setTitle(`${season.charAt(0).toUpperCase() + season.slice(1)} ${year} Anime${totalEmbeds > 1 ? ` (Part ${i + 1}/${totalEmbeds})` : ''}`)
        .setDescription(animeList)
        .setColor(0x02A9FF)
        
      if (i === 0) {
        embed.setFooter({ text: `Showing all ${seasonAnime.media.length} anime from ${season} ${year}` })
      }

      embeds.push(embed)
    }

    // Discord allows up to 10 embeds per message
    if (embeds.length <= 10) {
      await interaction.editReply({ embeds })
    } else {
      // Send first 10 embeds initially
      await interaction.editReply({ embeds: embeds.slice(0, 10) })
      
      // Send remaining embeds as follow-up messages
      for (let i = 10; i < embeds.length; i += 10) {
        const embedSlice = embeds.slice(i, i + 10)
        await interaction.followUp({ embeds: embedSlice })
      }
    }
  } catch (error) {
    console.error('Error in season command:', error)
    await interaction.editReply('❌ An error occurred while fetching seasonal anime.')
  }
}

function getStatusEmoji(status: AnimeStatus): string {
  switch (status) {
    case 'RELEASING':
      return '🟢'
    case 'FINISHED':
      return '✅'
    case 'NOT_YET_RELEASED':
      return '🔜'
    case 'CANCELLED':
      return '❌'
    case 'HIATUS':
      return '⏸️'
    default:
      return '❓'
  }
}
