import { redisCache } from '@/services/redis/cache'

// Redis key prefix for watchlists
const WATCHLIST_KEY_PREFIX = 'watchlist:user:'

export async function getUserWatchlist(userId: string): Promise<number[]> {
  try {
    const redisKey = `${WATCHLIST_KEY_PREFIX}${userId}`
    const watchlistItems = await redisCache.getSetMembers(redisKey)
    return watchlistItems.map((item) => parseInt(item, 10)).filter((id) => !isNaN(id))
  } catch (error) {
    console.error(`Error getting watchlist for user ${userId}:`, error)
    return []
  }
}

export async function addToWatchlist(
  userId: string,
  animeId: number
): Promise<{ success: boolean; message: string }> {
  try {
    const redisKey = `${WATCHLIST_KEY_PREFIX}${userId}`

    // Check if anime is already in watchlist
    const isAlreadyInWatchlist = await redisCache.isInSet(redisKey, animeId)
    if (isAlreadyInWatchlist) {
      return { success: false, message: 'Anime already in your watchlist.' }
    }

    // Add anime to watchlist
    await redisCache.addToSet(redisKey, animeId)

    // Set a reasonable TTL for the watchlist (30 days)
    await redisCache.expire(redisKey, 30 * 24 * 60 * 60)

    return { success: true, message: 'Anime added to your watchlist.' }
  } catch (error) {
    console.error(`Error adding anime ${animeId} to watchlist for user ${userId}:`, error)
    return { success: false, message: 'Failed to add anime to watchlist.' }
  }
}

export async function removeFromWatchlist(
  userId: string,
  animeId: number
): Promise<{ success: boolean; message: string }> {
  try {
    const redisKey = `${WATCHLIST_KEY_PREFIX}${userId}`

    // Check if anime is in watchlist
    const isInWatchlist = await redisCache.isInSet(redisKey, animeId)
    if (!isInWatchlist) {
      return { success: false, message: 'Anime not found in your watchlist.' }
    }

    // Remove anime from watchlist
    await redisCache.removeFromSet(redisKey, animeId)

    return { success: true, message: 'Anime removed from your watchlist.' }
  } catch (error) {
    console.error(`Error removing anime ${animeId} from watchlist for user ${userId}:`, error)
    return { success: false, message: 'Failed to remove anime from watchlist.' }
  }
}
