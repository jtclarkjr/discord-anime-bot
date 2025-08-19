import { IS_OPENAI_ENABLED } from './constants'

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
      name: 'notify',
      type: 2, // SUB_COMMAND_GROUP type
      description: 'Episode notification commands',
      options: [
        {
          name: 'add',
          type: 1, // SUB_COMMAND type
          description: 'Set notification for next episode',
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
          name: 'list',
          type: 1, // SUB_COMMAND type
          description: 'List your active episode notifications'
        },
        {
          name: 'cancel',
          type: 1, // SUB_COMMAND type
          description: 'Cancel notification for an anime',
          options: [
            {
              name: 'id',
              type: 4, // INTEGER type
              description: 'AniList ID of the anime',
              required: true
            }
          ]
        }
      ]
    },
    {
      name: 'watchlist',
      type: 2, // SUB_COMMAND_GROUP type
      description: 'Personal anime watchlist commands',
      options: [
        {
          name: 'add',
          type: 1, // SUB_COMMAND type
          description: 'Add an anime to your watchlist',
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
          name: 'list',
          type: 1, // SUB_COMMAND type
          description: 'Show your anime watchlist'
        },
        {
          name: 'remove',
          type: 1, // SUB_COMMAND type
          description: 'Remove an anime from your watchlist',
          options: [
            {
              name: 'id',
              type: 4, // INTEGER type
              description: 'AniList ID of the anime',
              required: true
            }
          ]
        }
      ]
    },
    {
      name: 'release',
      type: 1, // SUB_COMMAND type
      description: 'Show all currently releasing anime',
      options: [
        {
          name: 'page',
          type: 4, // INTEGER type
          description: 'Page number',
          required: false
        },
        {
          name: 'perpage',
          type: 4, // INTEGER type
          description: 'Anime per page (max 50)',
          required: false
        }
      ]
    },
    {
      name: 'season',
      type: 1, // SUB_COMMAND type
      description: 'Get all anime from a specific season and year',
      options: [
        {
          name: 'season',
          type: 3, // STRING type
          description: 'Season (winter, spring, summer, fall)',
          required: true,
          choices: [
            { name: 'Winter', value: 'winter' },
            { name: 'Spring', value: 'spring' },
            { name: 'Summer', value: 'summer' },
            { name: 'Fall', value: 'fall' }
          ]
        },
        {
          name: 'year',
          type: 4, // INTEGER type
          description: 'Year (defaults to current year)',
          required: false
        }
      ]
    },
    ...(IS_OPENAI_ENABLED
      ? [
          {
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
          }
        ]
      : []),
    {
      name: 'help',
      type: 1, // SUB_COMMAND type
      description: 'Show help for all /anime commands'
    }
  ]
} as const
