import { ChatInputCommandInteraction, EmbedBuilder } from 'discord.js'
import { findAnimeWithDetails } from '@/services/anime/find'

export async function handleFindCommand(interaction: ChatInputCommandInteraction) {
  const prompt = interaction.options.getString('prompt')
  if (!prompt) {
    await interaction.reply('❌ Please provide a description to search for anime.')
    return
  }

  await interaction.deferReply()

  try {
    const matches = await findAnimeWithDetails(prompt)
    
    if (matches.length === 0) {
      await interaction.editReply(`❌ No anime found matching the description: "${prompt}"`)
      return
    }

    // Create embed for the best match
    const bestMatch = matches[0]
    const anime = bestMatch.anime
    
    const embed = new EmbedBuilder()
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
      .setColor(0x00FF00)
      .setFooter({ text: 'Powered by ChatGPT-5 + AniList' })

    let responseText = `🤖 **AI Found Anime Based on:** "${prompt}"\n\n`
    
    if (matches.length > 1) {
      responseText += `**Other possible matches:**\n`
      responseText += matches.slice(1, 3).map((match, index) => 
        `${index + 2}. **${match.anime.title.english || match.anime.title.romaji}** (${Math.round(match.confidence * 100)}% match)`
      ).join('\n')
    }

    await interaction.editReply({ content: responseText, embeds: [embed] })

  } catch (error) {
    console.error('Error in find command:', error)
    await interaction.editReply('❌ An error occurred while finding anime. Please try again with a different description.')
  }
}
