import { ChatInputCommandInteraction, EmbedBuilder } from 'discord.js'
import { getAnimeById } from '@/services/anime/next'
import { formatAirDate } from '@/utils/formatters'
import {
  addNotification,
  getUserNotifications,
  removeUserNotification
} from '@/services/anime/notify'
import { createNotifyListEmbed, createNotifyCancelEmbed } from '@/embeds/notifyAnimeEmbed'

export async function handleNotifyCommand(interaction: ChatInputCommandInteraction) {
  const action = interaction.options.getString('action')
  const animeId = interaction.options.getInteger('id')

  // If no action is specified, show the list
  if (!action) {
    await handleNotifyListCommand(interaction)
    return
  }

  // Validate that ID is provided for actions that require it
  if ((action === 'add' || action === 'cancel') && !animeId) {
    await interaction.reply({
      content: 'Please provide an anime ID for this action.',
      flags: 1 << 6
    })
    return
  }

  switch (action) {
    case 'add':
      await handleNotifyAddCommand(interaction, animeId!)
      break
    case 'cancel':
      await handleNotifyCancelCommand(interaction, animeId!)
      break
    default:
      await interaction.reply('Unknown notify action.')
  }
}

async function handleNotifyAddCommand(interaction: ChatInputCommandInteraction, animeId: number) {
  await interaction.deferReply({ flags: 1 << 6 })

  // console.log(`[NotifyAdd] Starting add notification for anime ${animeId}`)

  try {
    // console.log(`[NotifyAdd] Calling notificationService.addNotification`)
    const result = await addNotification(animeId, interaction.channelId, interaction.user.id)

    // console.log(`[NotifyAdd] Result:`, result)

    if (result.success) {
      const embed = new EmbedBuilder()
        .setTitle('Notification Set!')
        .setDescription(result.message)
        .setColor(0x00ff00)

      if (result.airingDate) {
        embed.addFields({
          name: 'Airing Date',
          value: formatAirDate(result.airingDate),
          inline: false
        })
      }

      await interaction.editReply({ embeds: [embed] })
    } else {
      await interaction.editReply(`${result.message}`)
    }
  } catch (error) {
    console.error('[NotifyAdd] Error in notify add command:', error)
    try {
      await interaction.editReply('An error occurred while setting up the notification.')
    } catch (replyError) {
      console.error('[NotifyAdd] Failed to send error reply:', replyError)
    }
  }
}

async function handleNotifyListCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply({ flags: 1 << 6 })

  try {
    const notifications = getUserNotifications(interaction.user.id)

    if (notifications.length === 0) {
      await interaction.editReply('You have no active episode notifications.')
      return
    }

    let description = ''
    for (const notification of notifications) {
      const anime = await getAnimeById(notification.animeId)
      if (anime) {
        const title = anime.title.english || anime.title.romaji
        const airingDate = new Date(notification.airingAt)
        description += `• **${title}** (ID: ${notification.animeId})\n`
        description += `  Episode ${notification.episode} - ${formatAirDate(airingDate)}\n\n`
      }
    }
    const embed = createNotifyListEmbed(description)
    await interaction.editReply({ embeds: [embed] })
  } catch (error) {
    console.error('Error in notify list command:', error)
    await interaction.editReply('An error occurred while fetching your notifications.')
  }
}

async function handleNotifyCancelCommand(
  interaction: ChatInputCommandInteraction,
  animeId: number
) {
  await interaction.deferReply({ flags: 1 << 6 })

  try {
    const removed = await removeUserNotification(
      animeId,
      interaction.channelId,
      interaction.user.id
    )

    if (removed) {
      const anime = await getAnimeById(animeId)
      const title = anime ? anime.title.english || anime.title.romaji : `Anime ID ${animeId}`
      const embed = createNotifyCancelEmbed(`Notification canceled for **${title}**`)
      await interaction.editReply({ embeds: [embed] })
    } else {
      await interaction.editReply('No active notification found for this anime.')
    }
  } catch (error) {
    console.error('Error in notify cancel command:', error)
    try {
      await interaction.editReply('An error occurred while canceling the notification.')
    } catch (replyError) {
      console.error('Failed to send error reply:', replyError)
    }
  }
}
