import { createClient, type RedisClientType } from 'redis'
import { REDIS_URL } from '@config/constants'

class RedisConnection {
  private static instance: RedisConnection
  private client: RedisClientType | null = null
  private isConnecting = false
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5

  private constructor() {}

  public static getInstance(): RedisConnection {
    if (!RedisConnection.instance) {
      RedisConnection.instance = new RedisConnection()
    }
    return RedisConnection.instance
  }

  public async connect(): Promise<RedisClientType> {
    if (this.client?.isOpen) {
      return this.client
    }

    if (this.isConnecting) {
      // Wait for the connection attempt to complete
      return new Promise((resolve, reject) => {
        const checkConnection = () => {
          if (this.client?.isOpen) {
            resolve(this.client)
          } else if (!this.isConnecting) {
            reject(new Error('Failed to connect to Redis'))
          } else {
            setTimeout(checkConnection, 100)
          }
        }
        checkConnection()
      })
    }

    this.isConnecting = true

    try {
      this.client = createClient({
        url: REDIS_URL,
        socket: {
          reconnectStrategy: (retries) => {
            if (retries > this.maxReconnectAttempts) {
              console.error('Redis: Max reconnection attempts reached')
              return new Error('Max reconnection attempts reached')
            }
            const delay = Math.min(retries * 100, 3000)
            console.log(`Redis: Reconnecting in ${delay}ms (attempt ${retries})`)
            return delay
          }
        }
      })

      this.client.on('error', (err) => {
        console.error('Redis Client Error:', err)
      })

      this.client.on('connect', () => {
        console.log('Redis: Connected')
        this.reconnectAttempts = 0
      })

      this.client.on('ready', () => {
        console.log('Redis: Ready for commands')
      })

      this.client.on('end', () => {
        console.log('Redis: Connection closed')
      })

      this.client.on('reconnecting', () => {
        this.reconnectAttempts++
        console.log(`Redis: Reconnecting... (attempt ${this.reconnectAttempts})`)
      })

      await this.client.connect()
      this.isConnecting = false

      console.log('Redis connection established')
      return this.client
    } catch (error) {
      this.isConnecting = false
      console.error('Failed to connect to Redis:', error)
      throw error
    }
  }

  public async disconnect(): Promise<void> {
    if (this.client?.isOpen) {
      await this.client.quit()
      console.log('Redis: Disconnected')
    }
  }

  public getClient(): RedisClientType | null {
    return this.client
  }

  public isConnected(): boolean {
    return this.client?.isOpen ?? false
  }
}

export const redisConnection = RedisConnection.getInstance()
export type { RedisClientType }
