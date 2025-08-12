import { Client, TextChannel } from 'discord.js'
import { getAnimeById } from './next'
import type { NotificationEntry } from '@/types/anilist'



class NotificationService {
  private notifications: Map<string, NotificationEntry> = new Map()
  private client: Client | null = null
  private readonly storageFile = './data/notifications.json'

  setClient(client: Client) {
    this.client = client
    // Load existing notifications on startup
    this.loadNotifications()
  }

  /**
   * Load notifications from storage and reschedule them
   */
  private async loadNotifications() {
    try {
      const file = Bun.file(this.storageFile)
      
      if (!(await file.exists())) {
        return
      }

      const data = await file.text()
      const persistedNotifications = JSON.parse(data)
      const now = Date.now()

      for (const notification of persistedNotifications) {
        // Skip expired notifications
        if (notification.airingAt <= now) {
          continue
        }

        const notificationKey = `${notification.animeId}-${notification.channelId}-${notification.userId}`
        const timeUntilAiring = notification.airingAt - now

        const entry: NotificationEntry = {
          ...notification,
          timeoutId: setTimeout(() => {
            this.sendNotification(entry)
          }, timeUntilAiring)
        }

        this.notifications.set(notificationKey, entry)
      }

      console.log(`✅ Loaded ${this.notifications.size} active notifications from storage`)
    } catch (error) {
      console.error('Error loading notifications:', error)
    }
  }

    /**
   * Save notifications to persistent storage
   */
  private async saveNotifications() {
    try {
      // Convert to persistable format (without timeoutId)
      const persistedNotifications = Array.from(this.notifications.values()).map(entry => ({
        animeId: entry.animeId,
        channelId: entry.channelId,
        userId: entry.userId,
        airingAt: entry.airingAt,
        episode: entry.episode
      }))

      // Ensure we always save an array, never null
      const safeNotifications = persistedNotifications || []
      
      // Use Bun.write for efficient file writing
      await Bun.write(this.storageFile, JSON.stringify(safeNotifications, null, 2))
      console.log(`💾 [Save] Successfully saved ${safeNotifications.length} notifications to storage`)
    } catch (error) {
      console.error('💾 [Save] Error saving notifications:', error)
    }
  }

  /**
   * Add a notification for an anime episode
   */
  async addNotification(animeId: number, channelId: string, userId: string): Promise<{ success: boolean; message: string; airingDate?: Date }> {
    if (!this.client) {
      return { success: false, message: 'Bot client not initialized' }
    }

    try {
      const anime = await getAnimeById(animeId)
      
      if (!anime) {
        return { success: false, message: `No anime found with ID ${animeId}` }
      }

      if (!anime.nextAiringEpisode) {
        return { success: false, message: `No upcoming episodes scheduled for ${anime.title.english || anime.title.romaji}` }
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
      if (this.notifications.has(notificationKey)) {
        return { success: false, message: `You already have a notification set for ${anime.title.english || anime.title.romaji}` }
      }
      
      // Remove any other notifications for this anime by this user (in case they exist in other channels)
      const keysToRemove: string[] = []
      for (const [key, entry] of this.notifications.entries()) {
        if (entry.animeId === animeId && entry.userId === userId) {
          keysToRemove.push(key)
        }
      }
      
      // Remove the found notifications (without saving each time to avoid race condition)
      for (const key of keysToRemove) {
        this.removeNotificationInternal(key)
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
        this.sendNotification(entry)
      }, timeUntilAiring)

      this.notifications.set(notificationKey, entry)

      // Save to persistent storage
      await this.saveNotifications()

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
  private removeNotificationInternal(notificationKey: string): boolean {
    const entry = this.notifications.get(notificationKey)
    if (entry) {
      // Clear timeout if it exists
      if (entry.timeoutId) {
        clearTimeout(entry.timeoutId)
      }
      this.notifications.delete(notificationKey)
      return true
    }
    return false
  }

  /**
   * Remove a notification
   */
  async removeNotification(notificationKey: string): Promise<boolean> {
    const removed = this.removeNotificationInternal(notificationKey)
    if (removed) {
      await this.saveNotifications()
    }
    return removed
  }

  /**
   * Remove notifications for a specific user and anime
   */
  async removeUserNotification(animeId: number, channelId: string, userId: string): Promise<boolean> {
    const notificationKey = `${animeId}-${channelId}-${userId}`
    return await this.removeNotification(notificationKey)
  }

  /**
   * Send the notification when episode airs
   */
  private async sendNotification(entry: NotificationEntry) {
    if (!this.client) return

    try {
      const channel = await this.client.channels.fetch(entry.channelId) as TextChannel
      if (!channel) return

      const anime = await getAnimeById(entry.animeId)
      if (!anime) return

      const title = anime.title.english || anime.title.romaji
      const message = `🎉 <@${entry.userId}> Episode ${entry.episode} of **${title}** has just aired!\n\n🔗 [AniList Details](${anime.siteUrl})`

      await channel.send(message)
      
      // Remove the notification after sending
      const notificationKey = `${entry.animeId}-${entry.channelId}-${entry.userId}`
      this.notifications.delete(notificationKey)
      await this.saveNotifications()
    } catch (error) {
      console.error('Error sending notification:', error)
    }
  }

  /**
   * Get all active notifications for a user in a specific channel
   */
  getUserNotifications(userId: string, channelId?: string): NotificationEntry[] {
    return Array.from(this.notifications.values()).filter(entry => 
      entry.userId === userId && (!channelId || entry.channelId === channelId)
    )
  }

  /**
   * Clean up expired notifications (for safety)
   */
  async cleanup() {
    const now = Date.now()
    let hasChanges = false
    
    for (const [key, entry] of this.notifications.entries()) {
      if (entry.airingAt <= now) {
        if (entry.timeoutId) {
          clearTimeout(entry.timeoutId)
        }
        this.notifications.delete(key)
        hasChanges = true
      }
    }
    
    if (hasChanges) {
      await this.saveNotifications()
    }
  }
}

export const notificationService = new NotificationService()
