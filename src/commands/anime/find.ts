import { ChatInputCommandInteraction } from 'discord.js'
import { createFindAnimeEmbed } from '@/embeds/findAnimeEmbed'
import { findAnimeWithDetails } from '@/services/anime/find'
import { IS_AI_ENABLED } from '@/config/constants'

export async function handleFindCommand(interaction: ChatInputCommandInteraction) {
  if (!IS_AI_ENABLED) {
    await interaction.reply(
      '❌ The find command is disabled because no AI API key is configured. Please set the OPENAI_API_KEY or CLAUDE_API_KEY environment variable to use AI-powered anime search.'
    )
    return
  }

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
    const embed = createFindAnimeEmbed(anime, bestMatch)

    let responseText = `🤖 **AI Found Anime Based on:** "${prompt}"\n\n`

    if (matches.length > 1) {
      responseText += `**Other possible matches:**\n`
      responseText += matches
        .slice(1, 3)
        .map(
          (match, index) =>
            `${index + 2}. **${match.anime.title.english || match.anime.title.romaji}** (${Math.round(match.confidence * 100)}% match)`
        )
        .join('\n')
    }

    await interaction.editReply({ content: responseText, embeds: [embed] })
  } catch (error) {
    console.error('Error in find command:', error)
    await interaction.editReply(
      '❌ An error occurred while finding anime. Please try again with a different description.'
    )
  }
}
