import { EmbedBuilder } from 'discord.js'

export function createNotifySuccessEmbed(message: string, airingDate?: Date): EmbedBuilder {
  const embed = new EmbedBuilder()
    .setTitle('✅ Notification Set!')
    .setDescription(message)
    .setColor(0x00FF00)

  if (airingDate) {
    embed.addFields({
      name: 'Airing Date',
      value: airingDate.toLocaleString(),
      inline: false
    })
  }
  return embed
}

export function createNotifyErrorEmbed(message: string): EmbedBuilder {
  return new EmbedBuilder()
    .setTitle('❌ Notification Error')
    .setDescription(message)
    .setColor(0xFF0000)
}

export function createNotifyListEmbed(list: string): EmbedBuilder {
  return new EmbedBuilder()
    .setTitle('🔔 Your Anime Notifications')
    .setDescription(list)
    .setColor(0x02A9FF)
}

export function createNotifyCancelEmbed(message: string): EmbedBuilder {
  return new EmbedBuilder()
    .setTitle('🚫 Notification Cancelled')
    .setDescription(message)
    .setColor(0xFFA500)
}
