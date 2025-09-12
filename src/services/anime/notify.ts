import { Client, TextChannel } from 'discord.js'
import { getAnimeById } from './next'
import type { NotificationEntry } from '@/types/anilist'
import { redisCache } from '@/services/redis/cache'

let notifications: Map<string, NotificationEntry> = new Map()
let client: Client | null = null

/**
 * Set the Discord client and load notifications on startup
 */
export function setClient(newClient: Client) {
  client = newClient
  // Load existing notifications on startup
  loadNotifications()
}

/**
 * Load notifications from Redis and reschedule them
 */
async function loadNotifications() {
  try {
    // Get all notification keys from Redis
    const notificationKeys = await redisCache.getKeys('notification:*')

    if (notificationKeys.length === 0) {
      console.log('No notifications found in Redis')
      return
    }

    const now = Date.now()
    let loadedCount = 0

    for (const redisKey of notificationKeys) {
      const persistedNotification = await redisCache.get<NotificationEntry>(redisKey, true)

      if (!persistedNotification) continue

      // Skip expired notifications
      if (persistedNotification.airingAt <= now) {
        await redisCache.delete(redisKey)
        continue
      }

      const notificationKey = `${persistedNotification.animeId}-${persistedNotification.channelId}-${persistedNotification.userId}`
      const timeUntilAiring = persistedNotification.airingAt - now

      const entry: NotificationEntry = {
        ...persistedNotification,
        timeoutId: setTimeout(() => {
          sendNotification(entry)
        }, timeUntilAiring)
      }

      notifications.set(notificationKey, entry)
      loadedCount++
    }

    console.log(`Loaded ${loadedCount} active notifications from Redis`)
  } catch (error) {
    console.error('Error loading notifications from Redis:', error)
  }
}

/**
 * Save a single notification to Redis
 */
async function saveNotification(notificationKey: string, entry: NotificationEntry): Promise<void> {
  try {
    const redisKey = `notification:${notificationKey}`
    const persistedEntry = {
      animeId: entry.animeId,
      channelId: entry.channelId,
      userId: entry.userId,
      airingAt: entry.airingAt,
      episode: entry.episode
    }

    // Calculate TTL based on airing time (with some buffer)
    const now = Date.now()
    const ttlSeconds = Math.max(Math.ceil((entry.airingAt - now) / 1000) + 3600, 60) // At least 1 minute TTL

    await redisCache.set(redisKey, persistedEntry, ttlSeconds)
    console.log(`[Save] Successfully saved notification ${notificationKey} to Redis`)
  } catch (error) {
    console.error(`[Save] Error saving notification ${notificationKey} to Redis:`, error)
  }
}

/**
 * Remove a notification from Redis
 */
async function removeNotificationFromRedis(notificationKey: string): Promise<void> {
  try {
    const redisKey = `notification:${notificationKey}`
    await redisCache.delete(redisKey)
    console.log(`[Remove] Successfully removed notification ${notificationKey} from Redis`)
  } catch (error) {
    console.error(`[Remove] Error removing notification ${notificationKey} from Redis:`, error)
  }
}

/**
 * Add a notification for an anime episode
 */
export async function addNotification(
  animeId: number,
  channelId: string,
  userId: string
): Promise<{ success: boolean; message: string; airingDate?: Date }> {
  if (!client) {
    return { success: false, message: 'Bot client not initialized' }
  }
  try {
    const anime = await getAnimeById(animeId)

    if (!anime) {
      return { success: false, message: `No anime found with ID ${animeId}` }
    }

    if (!anime.nextAiringEpisode) {
      return {
        success: false,
        message: `No upcoming episodes scheduled for ${anime.title.english || anime.title.romaji}`
      }
    }

    if (anime.status === 'FINISHED') {
      return { success: false, message: `This anime has finished airing` }
    }

    if (anime.status === 'CANCELLED') {
      return { success: false, message: `This anime has been cancelled` }
    }

    const airingAt = anime.nextAiringEpisode.airingAt * 1000 // Convert to milliseconds
    const now = Date.now()

    if (airingAt <= now) {
      return { success: false, message: 'This episode has already aired' }
    }

    const notificationKey = `${animeId}-${channelId}-${userId}`

    // Check if notification already exists
    if (notifications.has(notificationKey)) {
      return {
        success: false,
        message: `You already have a notification set for ${anime.title.english || anime.title.romaji}`
      }
    }

    // Remove any other notifications for this anime by this user (in case they exist in other channels)
    const keysToRemove: string[] = []
    for (const [key, entry] of notifications.entries()) {
      if (entry.animeId === animeId && entry.userId === userId) {
        keysToRemove.push(key)
      }
    }

    // Remove the found notifications (without saving each time to avoid race condition)
    for (const key of keysToRemove) {
      removeNotificationInternal(key)
    }

    const airingDate = new Date(airingAt)
    const timeUntilAiring = airingAt - now

    const entry: NotificationEntry = {
      animeId,
      channelId,
      userId,
      airingAt,
      episode: anime.nextAiringEpisode.episode
    }

    // Schedule the notification
    entry.timeoutId = setTimeout(() => {
      sendNotification(entry)
    }, timeUntilAiring)

    notifications.set(notificationKey, entry)

    // Save to Redis
    await saveNotification(notificationKey, entry)

    return {
      success: true,
      message: `Notification set for ${anime.title.english || anime.title.romaji} Episode ${anime.nextAiringEpisode.episode}`,
      airingDate
    }
  } catch (error) {
    console.error('Error adding notification:', error)
    return { success: false, message: 'An error occurred while setting up the notification' }
  }
}

/**
 * Remove a notification (internal method without saving)
 */
function removeNotificationInternal(notificationKey: string): boolean {
  const entry = notifications.get(notificationKey)
  if (entry?.timeoutId) {
    clearTimeout(entry.timeoutId)
  }
  return notifications.delete(notificationKey)
}

/**
 * Remove a notification
 */
export async function removeNotification(notificationKey: string): Promise<boolean> {
  const removed = removeNotificationInternal(notificationKey)
  if (removed) {
    await removeNotificationFromRedis(notificationKey)
  }
  return removed
}

/**
 * Remove notifications for a specific user and anime
 */
export async function removeUserNotification(
  animeId: number,
  channelId: string,
  userId: string
): Promise<boolean> {
  const notificationKey = `${animeId}-${channelId}-${userId}`
  return await removeNotification(notificationKey)
}

/**
 * Send the notification when episode airs
 */
async function sendNotification(entry: NotificationEntry) {
  if (!client) return
  try {
    const channel = (await client.channels.fetch(entry.channelId)) as TextChannel
    if (!channel) return
    const anime = await getAnimeById(entry.animeId)
    if (!anime) return
    const title = anime.title.english || anime.title.romaji
    const message = `<@${entry.userId}> Episode ${entry.episode} of **${title}** has just aired!\n\n[AniList Details](${anime.siteUrl})`
    await channel.send(message)
    // Remove the notification after sending
    const notificationKey = `${entry.animeId}-${entry.channelId}-${entry.userId}`
    notifications.delete(notificationKey)
    await removeNotificationFromRedis(notificationKey)
  } catch (error) {
    console.error('Error sending notification:', error)
  }
}

/**
 * Get all active notifications for a user in a specific channel
 */
export function getUserNotifications(userId: string, channelId?: string): NotificationEntry[] {
  return Array.from(notifications.values()).filter(
    (entry) => entry.userId === userId && (!channelId || entry.channelId === channelId)
  )
}

/**
 * Clean up expired notifications (for safety)
 */
export async function cleanup() {
  const now = Date.now()
  const keysToRemove: string[] = []

  for (const [key, entry] of notifications.entries()) {
    if (entry.airingAt <= now) {
      if (entry.timeoutId) {
        clearTimeout(entry.timeoutId)
      }
      notifications.delete(key)
      keysToRemove.push(key)
    }
  }

  // Remove expired notifications from Redis
  for (const key of keysToRemove) {
    await removeNotificationFromRedis(key)
  }

  if (keysToRemove.length > 0) {
    console.log(`[Cleanup] Removed ${keysToRemove.length} expired notifications`)
  }
}
