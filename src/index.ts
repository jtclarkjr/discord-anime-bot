import { Client, GatewayIntentBits } from 'discord.js'
import { handleAnimeCommand, animeCommandDefinition } from '@/commands/index'
import { DISCORD_TOKEN, IS_CLAUDE_ENABLED, IS_OPENAI_ENABLED } from '@/config/constants'
import { setClient, cleanup } from '@/services/anime/notify'
import { redisConnection } from '@/services/redis'

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent
  ]
})

client.on('interactionCreate', async (interaction) => {
  if (!interaction.isChatInputCommand()) return

  if (interaction.commandName === 'anime') {
    try {
      await handleAnimeCommand(interaction)
    } catch (err) {
      console.error('Error handling anime command:', err)
      const reply = {
        content: 'An error occurred while processing your command.',
        ephemeral: true
      }

      if (interaction.replied || interaction.deferred) {
        await interaction.editReply(reply)
      } else {
        await interaction.reply(reply)
      }
    }
  }
})

client.once('clientReady', async () => {
  console.log(`Logged in as ${client.user?.tag}!`)

  // Initialize Redis connection
  try {
    await redisConnection.connect()
    console.log('Redis connection established')
  } catch (error) {
    console.error('Failed to connect to Redis:', error)
    console.warn('Bot will continue without Redis (features may not work properly)')
  }

  // Initialize notification service
  setClient(client)

  // Clean up expired notifications every hour
  setInterval(
    async () => {
      await cleanup()
    },
    60 * 60 * 1000
  ) // 1 hour

  if (IS_OPENAI_ENABLED && !IS_CLAUDE_ENABLED) {
    console.log('OpenAI features enabled (find command available)')
  } else if (!IS_OPENAI_ENABLED && IS_CLAUDE_ENABLED) {
    console.log('Claude features enabled (find command available)')
  } else {
    console.log('No AI features enabled (find command not available)')
  }

  // Register the /anime command with subcommands
  const guild = client.guilds.cache.first()
  if (guild) {
    try {
      await guild.commands.create(animeCommandDefinition)
      console.log('Successfully registered anime command!')
    } catch (error) {
      console.error('Failed to register anime command:', error)
    }
  } else {
    console.warn('No guild found. Commands may not be registered.')
  }
})

// Graceful shutdown handling
process.on('SIGINT', async () => {
  console.log('\nShutting down gracefully...')

  try {
    await redisConnection.disconnect()
  } catch (error) {
    console.error('Error disconnecting from Redis:', error)
  }

  client.destroy()
  process.exit(0)
})

process.on('SIGTERM', async () => {
  console.log('\nReceived SIGTERM, shutting down gracefully...')

  try {
    await redisConnection.disconnect()
  } catch (error) {
    console.error('Error disconnecting from Redis:', error)
  }

  client.destroy()
  process.exit(0)
})

client.login(DISCORD_TOKEN)
