import { Client, TextChannel } from 'discord.js'
import { getAnimeById } from './next'
import type { NotificationEntry } from '@/types/anilist'
import { storageFile } from '@config/constants'

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
 * Load notifications from storage and reschedule them
 */
async function loadNotifications() {
  try {
    const file = Bun.file(storageFile)
    if (!(await file.exists())) {
      return
    }
    const data = await file.text()
    const parseJson = <T>(json: string): T => JSON.parse(json) as T
    const persistedNotifications = parseJson<NotificationEntry[]>(data)
    const now = Date.now()
    for (const notification of persistedNotifications) {
      if (notification.airingAt <= now) {
        continue
      }
      const notificationKey = `${notification.animeId}-${notification.channelId}-${notification.userId}`
      const timeUntilAiring = notification.airingAt - now
      const entry: NotificationEntry = {
        ...notification,
        timeoutId: setTimeout(() => {
          sendNotification(entry)
        }, timeUntilAiring)
      }
      notifications.set(notificationKey, entry)
    }
    console.log(`✅ Loaded ${notifications.size} active notifications from storage`)
  } catch (error) {
    console.error('Error loading notifications:', error)
  }
}

/**
 * Save notifications to persistent storage
 */
async function saveNotifications() {
  try {
    const persistedNotifications = Array.from(notifications.values()).map((entry) => ({
      animeId: entry.animeId,
      channelId: entry.channelId,
      userId: entry.userId,
      airingAt: entry.airingAt,
      episode: entry.episode
    }))
    const safeNotifications = persistedNotifications || []
    await Bun.write(storageFile, JSON.stringify(safeNotifications, null, 2))
    console.log(`💾 [Save] Successfully saved ${safeNotifications.length} notifications to storage`)
  } catch (error) {
    console.error('💾 [Save] Error saving notifications:', error)
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

    // Save to persistent storage
    await saveNotifications()

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
    await saveNotifications()
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
    const message = `🎉 <@${entry.userId}> Episode ${entry.episode} of **${title}** has just aired!\n\n🔗 [AniList Details](${anime.siteUrl})`
    await channel.send(message)
    // Remove the notification after sending
    const notificationKey = `${entry.animeId}-${entry.channelId}-${entry.userId}`
    notifications.delete(notificationKey)
    await saveNotifications()
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
  let hasChanges = false
  for (const [key, entry] of notifications.entries()) {
    if (entry.airingAt <= now) {
      if (entry.timeoutId) {
        clearTimeout(entry.timeoutId)
      }
      notifications.delete(key)
      hasChanges = true
    }
  }
  if (hasChanges) {
    await saveNotifications()
  }
}
