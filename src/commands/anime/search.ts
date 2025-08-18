import { ChatInputCommandInteraction } from 'discord.js'
import { createSearchAnimeEmbed } from '@/embeds/searchAnimeEmbed'
import { searchAnime } from '@/services/anime/search'

export async function handleSearchCommand(interaction: ChatInputCommandInteraction) {
  const searchQuery = interaction.options.getString('query')
  if (!searchQuery) {
    await interaction.reply('❌ Please provide an anime name to search for.')
    return
  }

  await interaction.deferReply()

  try {
    const searchResults = await searchAnime(searchQuery)

    if (searchResults.media.length === 0) {
      await interaction.editReply(`❌ No anime found for "${searchQuery}".`)
      return
    }

    // Create embed for the first result
    const anime = searchResults.media[0]
    const embed = createSearchAnimeEmbed(anime)

    let responseText = `Found ${searchResults.pageInfo.total} result(s) for "${searchQuery}"`
    if (searchResults.media.length > 1) {
      responseText += `\n\nShowing top result. Other matches:\n${searchResults.media
        .slice(1, 5)
        .map((a) => `• ${a.title.english || a.title.romaji}`)
        .join('\n')}`
    }

    await interaction.editReply({ content: responseText, embeds: [embed] })
  } catch (error) {
    console.error('Error in search command:', error)
    await interaction.editReply('❌ An error occurred while searching for anime.')
  }
}
