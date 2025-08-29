import { Client, GatewayIntentBits, EmbedBuilder } from 'discord.js'

const DISCORD_TOKEN = process.env.DISCORD_BOT_TOKEN!
const CHANNEL_ID = process.env.CHANNEL_ID!
const ANILIST_API = process.env.ANILIST_API!

type AnimeSearchResponse = {
  data: {
    Page: {
      media: {
        id: number
        title: { romaji: string; english: string | null; native: string }
        format: string
        status: string
        coverImage: { large: string }
        siteUrl: string
      }[]
      pageInfo: { total: number; currentPage: number; lastPage: number; hasNextPage: boolean }
    }
  }
}

type AnimeDetailsResponse = {
  data: {
    Media: {
      id: number
      title: { romaji: string; english: string | null; native: string }
      status: string
      format: string
      episodes: number | null
      nextAiringEpisode: {
        episode: number
        airingAt: number
        timeUntilAiring: number
      } | null
      coverImage: { large: string }
      siteUrl: string
    }
  }
}

type ReleasingAnimeResponse = {
  data: {
    Page: {
      media: {
        id: number
        title: { romaji: string; english: string | null }
        nextAiringEpisode: {
          episode: number
          airingAt: number
        } | null
      }[]
      pageInfo: { total: number; currentPage: number; lastPage: number; hasNextPage: boolean }
    }
  }
}

// Check for required env vars
if (!DISCORD_TOKEN) {
  console.error('❌ DISCORD_BOT_TOKEN is not set in environment variables.')
  process.exit(1)
}
if (!CHANNEL_ID) {
  console.error('❌ CHANNEL_ID is not set in environment variables.')
  process.exit(1)
}

const client = new Client({
  intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMessages, GatewayIntentBits.MessageContent]
})

/**
 * Search for anime using AniList API
 */
async function searchAnime(searchQuery: string, page: number = 1, perPage: number = 10) {
  const query = `
    query ($q: String!, $page: Int = 1, $perPage: Int = 10) {
      Page(page: $page, perPage: $perPage) {
        media(search: $q, type: ANIME, sort: [SEARCH_MATCH, POPULARITY_DESC]) {
          id
          title { romaji english native }
          format
          status
          coverImage { large }
          siteUrl
        }
        pageInfo { total currentPage lastPage hasNextPage }
      }
    }
  `
  
  const variables = {
    q: searchQuery,
    page,
    perPage
  }

  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })

  const json = (await res.json()) as AnimeSearchResponse

  return json.data.Page
}

/**
 * Get anime details by ID including next airing episode
 */
async function getAnimeById(animeId: number) {
  const query = `
    query ($id: Int!) {
      Media(id: $id, type: ANIME) {
        id
        title { romaji english native }
        status
        format
        episodes
        nextAiringEpisode {
          episode
          airingAt
          timeUntilAiring
        }
        coverImage { large }
        siteUrl
      }
    }
  `
  
  const variables = {
    id: animeId
  }

  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })

  const json = (await res.json()) as AnimeDetailsResponse

  return json.data.Media
}

/**
 * Get all currently releasing anime
 */
async function getReleasingAnime(page: number = 1, perPage: number = 25) {
  const query = `
    query ($page: Int, $perPage: Int) {
      Page(page: $page, perPage: $perPage) {
        media(type: ANIME, status: RELEASING, sort: [POPULARITY_DESC]) {
          id
          title { romaji english }
          nextAiringEpisode {
            episode
            airingAt
          }
        }
        pageInfo { total currentPage lastPage hasNextPage }
      }
    }
  `
  
  const variables = {
    page,
    perPage
  }

  const res = await fetch(ANILIST_API, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query, variables })
  })

  const json = (await res.json()) as ReleasingAnimeResponse

  return json.data.Page
}



client.on('interactionCreate', async (interaction) => {
  if (!interaction.isChatInputCommand()) return

  if (interaction.commandName === 'anime') {
    try {
      const subcommand = interaction.options.getSubcommand()

      if (subcommand === 'search') {
        const searchQuery = interaction.options.getString('query')
        if (!searchQuery) {
          await interaction.reply('❌ Please provide an anime name to search for.')
          return
        }

        await interaction.deferReply()

        const searchResults = await searchAnime(searchQuery)
        
        if (searchResults.media.length === 0) {
          await interaction.editReply(`❌ No anime found for "${searchQuery}".`)
          return
        }

        // Create embed for the first result
        const anime = searchResults.media[0]
        const embed = new EmbedBuilder()
          .setTitle(anime.title.english || anime.title.romaji)
          .setURL(anime.siteUrl)
          .setThumbnail(anime.coverImage.large)
          .addFields(
            { name: 'Romaji Title', value: anime.title.romaji, inline: true },
            { name: 'Native Title', value: anime.title.native, inline: true },
            { name: 'Format', value: anime.format, inline: true },
            { name: 'Status', value: anime.status, inline: true },
            { name: 'AniList ID', value: anime.id.toString(), inline: true }
          )
          .setColor(0x02A9FF)

        let responseText = `Found ${searchResults.pageInfo.total} result(s) for "${searchQuery}"`
        if (searchResults.media.length > 1) {
          responseText += `\n\nShowing top result. Other matches:\n${searchResults.media.slice(1, 5).map(a => `• ${a.title.english || a.title.romaji}`).join('\n')}`
        }

        await interaction.editReply({ content: responseText, embeds: [embed] })

      } else if (subcommand === 'next') {
        const animeId = interaction.options.getInteger('id')
        if (!animeId) {
          await interaction.reply('❌ Please provide a valid anime ID.')
          return
        }

        await interaction.deferReply()

        const anime = await getAnimeById(animeId)
        
        if (!anime) {
          await interaction.editReply(`❌ No anime found with ID ${animeId}.`)
          return
        }

        const embed = new EmbedBuilder()
          .setTitle(anime.title.english || anime.title.romaji)
          .setURL(anime.siteUrl)
          .setThumbnail(anime.coverImage.large)
          .setColor(0x02A9FF)

        // Check if anime is finished or has next episode
        if (anime.status === 'FINISHED' || anime.status === 'CANCELLED') {
          embed.setDescription(`This anime has finished airing.`)
          embed.addFields(
            { name: 'Status', value: anime.status, inline: true },
            { name: 'Format', value: anime.format, inline: true },
            { name: 'Total Episodes', value: anime.episodes?.toString() || 'Unknown', inline: true }
          )
        } else if (anime.nextAiringEpisode) {
          const airingDate = new Date(anime.nextAiringEpisode.airingAt * 1000)
          const timeUntil = anime.nextAiringEpisode.timeUntilAiring
          
          // Convert seconds to days, hours, minutes
          const days = Math.floor(timeUntil / (24 * 60 * 60))
          const hours = Math.floor((timeUntil % (24 * 60 * 60)) / (60 * 60))
          const minutes = Math.floor((timeUntil % (60 * 60)) / 60)
          
          let timeString = ''
          if (days > 0) timeString += `${days} day${days > 1 ? 's' : ''} `
          if (hours > 0) timeString += `${hours} hour${hours > 1 ? 's' : ''} `
          if (minutes > 0) timeString += `${minutes} minute${minutes > 1 ? 's' : ''}`
          
          // Format date and time in user-friendly format
          const dateOptions: Intl.DateTimeFormatOptions = {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric'
          }
          const timeOptions: Intl.DateTimeFormatOptions = {
            hour: 'numeric',
            minute: '2-digit',
            hour12: true
          }
          
          const formattedDate = airingDate.toLocaleDateString('en-US', dateOptions)
          const formattedTime = airingDate.toLocaleTimeString('en-US', timeOptions)
          
          embed.setDescription(`Episode ${anime.nextAiringEpisode.episode} airs in ${timeString.trim()}`)
          embed.addFields(
            { name: 'Next Episode', value: anime.nextAiringEpisode.episode.toString(), inline: true },
            { name: 'Air Date', value: `${formattedDate} at ${formattedTime}`, inline: false },
            { name: 'Status', value: anime.status, inline: true }
          )
        } else {
          embed.setDescription(`No upcoming episodes scheduled.`)
          embed.addFields(
            { name: 'Status', value: anime.status, inline: true },
            { name: 'Format', value: anime.format, inline: true }
          )
        }

        await interaction.editReply({ embeds: [embed] })
      
      } else if (subcommand === 'release') {
        await interaction.deferReply()

        const releasingAnime = await getReleasingAnime()
        
        if (releasingAnime.media.length === 0) {
          await interaction.editReply('❌ No releasing anime found.')
          return
        }

        // Create a list of currently releasing anime
        const animeList = releasingAnime.media.slice(0, 15).map(anime => {
          const title = anime.title.english || anime.title.romaji
          let nextEpisodeInfo = ''
          
          if (anime.nextAiringEpisode) {
            const airingDate = new Date(anime.nextAiringEpisode.airingAt * 1000)
            const timeOptions: Intl.DateTimeFormatOptions = {
              month: 'short',
              day: 'numeric',
              hour: 'numeric',
              minute: '2-digit',
              hour12: true
            }
            const formattedTime = airingDate.toLocaleDateString('en-US', timeOptions)
            nextEpisodeInfo = ` - Ep ${anime.nextAiringEpisode.episode} on ${formattedTime}`
          } else {
            nextEpisodeInfo = ' - No schedule'
          }
          
          return `**${title}** (ID: ${anime.id})${nextEpisodeInfo}`
        }).join('\n')

        const embed = new EmbedBuilder()
          .setTitle('Currently Releasing Anime')
          .setDescription(animeList)
          .setColor(0x02A9FF)
          .setFooter({ text: `Showing ${Math.min(15, releasingAnime.media.length)} of ${releasingAnime.pageInfo.total} releasing anime` })

        await interaction.editReply({ embeds: [embed] })
      }

    } catch (err) {
      console.error('Error fetching anime:', err)
      await interaction.editReply('❌ An error occurred while fetching anime data.')
    }
  }
})

client.once('ready', async () => {
  console.log(`Logged in as ${client.user?.tag}!`)

  // Register the /anime command with subcommands
  const guild = client.guilds.cache.first()
  if (guild) {
    await guild.commands.create({
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
        }
      ]
    })
  }
})

client.login(DISCORD_TOKEN)
