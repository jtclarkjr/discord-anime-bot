import 'dotenv/config'
import { Client, GatewayIntentBits } from 'discord.js'
import { handleAnimeCommand, animeCommandDefinition } from '@/commands/index'
import { DISCORD_TOKEN, IS_OPENAI_ENABLED } from '@/config/constants'
import { notificationService } from '@/services/anime/notify'

const client = new Client({
  intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMessages, GatewayIntentBits.MessageContent]
})

client.on('interactionCreate', async (interaction) => {
  if (!interaction.isChatInputCommand()) return

  if (interaction.commandName === 'anime') {
    try {
      await handleAnimeCommand(interaction)
    } catch (err) {
      console.error('Error handling anime command:', err)
      const reply = { content: '❌ An error occurred while processing your command.', ephemeral: true }
      
      if (interaction.replied || interaction.deferred) {
        await interaction.editReply(reply)
      } else {
        await interaction.reply(reply)
      }
    }
  }
})

client.once('ready', async () => {
  console.log(`Logged in as ${client.user?.tag}!`)
  
  // Initialize notification service
  notificationService.setClient(client)
  
  // Clean up expired notifications every hour
  setInterval(async () => {
    await notificationService.cleanup()
  }, 60 * 60 * 1000) // 1 hour
  
  if (IS_OPENAI_ENABLED) {
    console.log('✅ OpenAI features enabled (find command available)')
  } else {
    console.log('⚠️ OpenAI features disabled (find command not available)')
  }

  // Register the /anime command with subcommands
  const guild = client.guilds.cache.first()
  if (guild) {
    try {
      await guild.commands.create(animeCommandDefinition)
      console.log('✅ Successfully registered anime command!')
    } catch (error) {
      console.error('❌ Failed to register anime command:', error)
    }
  } else {
    console.warn('⚠️ No guild found. Commands may not be registered.')
  }
})

client.login(DISCORD_TOKEN)
