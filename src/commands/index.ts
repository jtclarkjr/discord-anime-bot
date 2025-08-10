import { ChatInputCommandInteraction } from 'discord.js'
import { handleSearchCommand } from './anime/search'
import { handleNextCommand } from './anime/next'
import { handleReleaseCommand } from './anime/release'
import { handleFindCommand } from './anime/find'
import { IS_OPENAI_ENABLED } from '@/config/constants'

export async function handleAnimeCommand(interaction: ChatInputCommandInteraction) {
  const subcommand = interaction.options.getSubcommand()

  switch (subcommand) {
    case 'search':
      await handleSearchCommand(interaction)
      break
    case 'next':
      await handleNextCommand(interaction)
      break
    case 'release':
      await handleReleaseCommand(interaction)
      break
    case 'find':
      if (!IS_OPENAI_ENABLED) {
        await interaction.reply('❌ The find command is disabled because OpenAI API key is not configured.')
        return
      }
      await handleFindCommand(interaction)
      break
    default:
      await interaction.reply('❌ Unknown subcommand.')
  }
}

export const animeCommandDefinition = {
  name: 'anime',
  description: 'Anime-related commands',
  options: [
    {
      name: 'search',
      type: 1, // SUB_COMMAND type
      description: 'Search for anime information',
      options: [
        {
          name: 'query',
          type: 3, // STRING type
          description: 'Name of the anime to search for',
          required: true
        }
      ]
    },
    {
      name: 'next',
      type: 1, // SUB_COMMAND type
      description: 'Get next airing episode information',
      options: [
        {
          name: 'id',
          type: 4, // INTEGER type
          description: 'AniList ID of the anime',
          required: true
        }
      ]
    },
    {
      name: 'release',
      type: 1, // SUB_COMMAND type
      description: 'Show all currently releasing anime'
    },
    ...(IS_OPENAI_ENABLED ? [{
      name: 'find',
      type: 1, // SUB_COMMAND type
      description: 'Find anime using AI based on description',
      options: [
        {
          name: 'prompt',
          type: 3, // STRING type
          description: 'Describe the anime you are looking for',
          required: true
        }
      ]
    }] : [])
  ]
}
