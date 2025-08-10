import { ChatInputCommandInteraction, EmbedBuilder } from 'discord.js'
import { getAnimeById } from '@/services/anime/next'
import { formatAirDate } from '@/utils/formatters'
import { notificationService } from '@/services/anime/notify'

export async function handleNotifyCommand(interaction: ChatInputCommandInteraction) {
  const subcommand = interaction.options.getSubcommand()
  
  switch (subcommand) {
    case 'add':
      await handleNotifyAddCommand(interaction)
      break
    case 'list':
      await handleNotifyListCommand(interaction)
      break
    case 'cancel':
      await handleNotifyCancelCommand(interaction)
      break
    default:
      await interaction.reply('❌ Unknown notify subcommand.')
  }
}

async function handleNotifyAddCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply()

  const animeId = interaction.options.getInteger('id')
  if (!animeId) {
    await interaction.editReply('❌ Please provide a valid anime ID.')
    return
  }

  // console.log(`🔔 [NotifyAdd] Starting add notification for anime ${animeId}`)

  try {
    // console.log(`🔔 [NotifyAdd] Calling notificationService.addNotification`)
    const result = await notificationService.addNotification(
      animeId,
      interaction.channelId,
      interaction.user.id
    )

    // console.log(`🔔 [NotifyAdd] Result:`, result)

    if (result.success) {
      const embed = new EmbedBuilder()
        .setTitle('✅ Notification Set!')
        .setDescription(result.message)
        .setColor(0x00FF00)
        
      if (result.airingDate) {
        embed.addFields({
          name: 'Airing Date',
          value: formatAirDate(result.airingDate),
          inline: false
        })
      }

      await interaction.editReply({ embeds: [embed] })
    } else {
      await interaction.editReply(`❌ ${result.message}`)
    }
  } catch (error) {
    console.error('🔔 [NotifyAdd] Error in notify add command:', error)
    try {
      await interaction.editReply('❌ An error occurred while setting up the notification.')
    } catch (replyError) {
      console.error('🔔 [NotifyAdd] Failed to send error reply:', replyError)
    }
  }
}

async function handleNotifyListCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply()

  try {
    const notifications = notificationService.getUserNotifications(interaction.user.id)
    
    if (notifications.length === 0) {
      await interaction.editReply('📋 You have no active episode notifications.')
      return
    }

    const embed = new EmbedBuilder()
      .setTitle('📋 Your Active Episode Notifications')
      .setColor(0x02A9FF)

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

    embed.setDescription(description)
    await interaction.editReply({ embeds: [embed] })
  } catch (error) {
    console.error('Error in notify list command:', error)
    await interaction.editReply('❌ An error occurred while fetching your notifications.')
  }
}

async function handleNotifyCancelCommand(interaction: ChatInputCommandInteraction) {
  await interaction.deferReply()

  const animeId = interaction.options.getInteger('id')
  if (!animeId) {
    await interaction.editReply('❌ Please provide a valid anime ID.')
    return
  }

  try {
    const removed = await notificationService.removeUserNotification(
      animeId,
      interaction.channelId,
      interaction.user.id
    )

    if (removed) {
      const anime = await getAnimeById(animeId)
      const title = anime ? (anime.title.english || anime.title.romaji) : `Anime ID ${animeId}`
      
      const embed = new EmbedBuilder()
        .setTitle('🗑️ Notification Canceled')
        .setDescription(`Notification canceled for **${title}**`)
        .setColor(0xFF9900)

      await interaction.editReply({ embeds: [embed] })
    } else {
      await interaction.editReply('❌ No active notification found for this anime.')
    }
  } catch (error) {
    console.error('Error in notify cancel command:', error)
    try {
      await interaction.editReply('❌ An error occurred while canceling the notification.')
    } catch (replyError) {
      console.error('Failed to send error reply:', replyError)
    }
  }
}
