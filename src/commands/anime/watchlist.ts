import { ChatInputCommandInteraction } from 'discord.js'
import { getUserWatchlist, addToWatchlist, removeFromWatchlist } from '@/services/anime/watchlist'
import { createWatchlistEmbed } from '@/embeds/watchlistAnimeEmbed'

export async function handleWatchlistCommand(interaction: ChatInputCommandInteraction) {
  const subcommand = interaction.options.getSubcommand()

  switch (subcommand) {
    case 'add':
      await handleWatchlistAddCommand(interaction)
      break
    case 'list':
      await handleWatchlistListCommand(interaction)
      break
    case 'remove':
      await handleWatchlistRemoveCommand(interaction)
      break
    default:
      await interaction.reply({ content: '❌ Unknown watchlist subcommand.', flags: 1 << 6 })
  }
}

async function handleWatchlistAddCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply({ flags: 1 << 6 })
  const animeId = interaction.options.getInteger('id')
  if (!animeId) {
    await interaction.editReply('❌ Please provide a valid anime ID.')
    return
  }
  try {
    const result = await addToWatchlist(interaction.user.id, animeId)
    if (result.success) {
      // Fetch anime name for confirmation
      const anime = await (await import('@/services/anime/next')).getAnimeById(animeId)
      const title = anime ? anime.title.english || anime.title.romaji : `Anime ID ${animeId}`
      await interaction.editReply(`✅ Added **${title}** (ID: ${animeId}) to your watchlist.`)
    } else {
      await interaction.editReply(`❌ ${result.message}`)
    }
  } catch (error) {
    console.error('Error adding to watchlist:', error)
    await interaction.editReply('❌ An error occurred while adding to your watchlist.')
  }
}

async function handleWatchlistListCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply({ flags: 1 << 6 })
  try {
    const watchlist = await getUserWatchlist(interaction.user.id)
    if (!watchlist || watchlist.length === 0) {
      await interaction.editReply('📋 Your watchlist is empty.')
      return
    }
    const embed = await createWatchlistEmbed(watchlist)
    await interaction.editReply({ embeds: [embed] })
  } catch (error) {
    console.error('Error fetching watchlist:', error)
    await interaction.editReply('❌ An error occurred while fetching your watchlist.')
  }
}

async function handleWatchlistRemoveCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply({ flags: 1 << 6 })
  const animeId = interaction.options.getInteger('id')
  if (!animeId) {
    await interaction.editReply('❌ Please provide a valid anime ID.')
    return
  }
  try {
    const result = await removeFromWatchlist(interaction.user.id, animeId)
    if (result.success) {
      // Fetch anime name for confirmation
      const anime = await (await import('@/services/anime/next')).getAnimeById(animeId)
      const title = anime ? anime.title.english || anime.title.romaji : `Anime ID ${animeId}`
      await interaction.editReply(`✅ Removed **${title}** (ID: ${animeId}) from your watchlist.`)
    } else {
      await interaction.editReply(`❌ ${result.message}`)
    }
  } catch (error) {
    console.error('Error removing from watchlist:', error)
    await interaction.editReply('❌ An error occurred while removing from your watchlist.')
  }
}
