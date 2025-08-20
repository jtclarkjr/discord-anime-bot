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
      type: 1, // SUB_COMMAND type
      description: 'Manage episode notifications (shows list by default)',
      options: [
        {
          name: 'action',
          type: 3, // STRING type
          description: 'Action to perform (add or cancel)',
          required: false,
          choices: [
            { name: 'Add notification', value: 'add' },
            { name: 'Cancel notification', value: 'cancel' }
          ]
        },
        {
          name: 'id',
          type: 4, // INTEGER type
          description: 'AniList ID of the anime (required for add/cancel)',
          required: false
        }
      ]
    },
    {
      name: 'watchlist',
      type: 1, // SUB_COMMAND type
      description: 'Manage your anime watchlist (shows list by default)',
      options: [
        {
          name: 'action',
          type: 3, // STRING type
          description: 'Action to perform (add or remove)',
          required: false,
          choices: [
            { name: 'Add to watchlist', value: 'add' },
            { name: 'Remove from watchlist', value: 'remove' }
          ]
        },
        {
          name: 'id',
          type: 4, // INTEGER type
          description: 'AniList ID of the anime (required for add/remove)',
          required: false
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
