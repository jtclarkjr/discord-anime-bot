import { ChatInputCommandInteraction } from 'discord.js'
import { handleSearchCommand } from './anime/search'
import { handleNextCommand } from './anime/next'
import { handleReleaseCommand } from './anime/release'
import { handleSeasonCommand } from './anime/season'
import { handleFindCommand } from './anime/find'
import { handleNotifyCommand } from './anime/notify'
import { IS_OPENAI_ENABLED } from '@/config/constants'
import { handleHelpCommand } from './anime/help'
import { handleWatchlistCommand } from './anime/watchlist'
export { animeCommandDefinition } from '@/config/commands'

export async function handleAnimeCommand(interaction: ChatInputCommandInteraction) {
  const subcommand = interaction.options.getSubcommand()

  switch (subcommand) {
    case 'search':
      await handleSearchCommand(interaction)
      break
    case 'next':
      await handleNextCommand(interaction)
      break
    case 'notify':
      await handleNotifyCommand(interaction)
      break
    case 'watchlist':
      await handleWatchlistCommand(interaction)
      break
    case 'release':
      await handleReleaseCommand(interaction)
      break
    case 'season':
      await handleSeasonCommand(interaction)
      break
    case 'find':
      if (!IS_OPENAI_ENABLED) {
        await interaction.reply(
          'The find command is disabled because OpenAI API key is not configured.'
        )
        return
      }
      await handleFindCommand(interaction)
      break
    case 'help':
      await handleHelpCommand(interaction)
      break
    default:
      await interaction.reply('Unknown subcommand.')
  }
}
