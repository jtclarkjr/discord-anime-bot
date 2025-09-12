import { redisConnection, type RedisClientType } from './connection'

export class RedisCache {
  private client: RedisClientType | null = null

  constructor() {
    this.initialize()
  }

  private async initialize(): Promise<void> {
    try {
      this.client = await redisConnection.connect()
    } catch (error) {
      console.error('Failed to initialize Redis cache:', error)
      throw error
    }
  }

  private async ensureConnection(): Promise<RedisClientType> {
    if (!this.client || !redisConnection.isConnected()) {
      this.client = await redisConnection.connect()
    }
    return this.client
  }

  /**
   * Set a key-value pair with optional expiration
   */
  async set(key: string, value: string | object, ttlSeconds?: number): Promise<void> {
    const client = await this.ensureConnection()
    const serializedValue = typeof value === 'string' ? value : JSON.stringify(value)

    if (ttlSeconds) {
      await client.setEx(key, ttlSeconds, serializedValue)
    } else {
      await client.set(key, serializedValue)
    }
  }

  /**
   * Get a value by key and optionally parse as JSON
   */
  async get<T = string>(key: string, parseJson = false): Promise<T | null> {
    const client = await this.ensureConnection()
    const value = await client.get(key)

    if (value === null) return null

    if (parseJson) {
      try {
        return JSON.parse(value) as T
      } catch (error) {
        console.error(`Failed to parse JSON for key ${key}:`, error)
        return null
      }
    }

    return value as T
  }

  /**
   * Delete a key
   */
  async delete(key: string): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.del(key)
    return result > 0
  }

  /**
   * Check if a key exists
   */
  async exists(key: string): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.exists(key)
    return result === 1
  }

  /**
   * Get all keys matching a pattern
   */
  async getKeys(pattern: string): Promise<string[]> {
    const client = await this.ensureConnection()
    return await client.keys(pattern)
  }

  /**
   * Add an item to a set
   */
  async addToSet(key: string, value: string | number): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.sAdd(key, String(value))
    return result > 0
  }

  /**
   * Remove an item from a set
   */
  async removeFromSet(key: string, value: string | number): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.sRem(key, String(value))
    return result > 0
  }

  /**
   * Get all members of a set
   */
  async getSetMembers(key: string): Promise<string[]> {
    const client = await this.ensureConnection()
    return await client.sMembers(key)
  }

  /**
   * Check if a value is in a set
   */
  async isInSet(key: string, value: string | number): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.sIsMember(key, String(value))
    return result === 1
  }

  /**
   * Add an item to a hash
   */
  async setHashField(key: string, field: string, value: string | object): Promise<void> {
    const client = await this.ensureConnection()
    const serializedValue = typeof value === 'string' ? value : JSON.stringify(value)
    await client.hSet(key, field, serializedValue)
  }

  /**
   * Get an item from a hash
   */
  async getHashField<T = string>(key: string, field: string, parseJson = false): Promise<T | null> {
    const client = await this.ensureConnection()
    const value = await client.hGet(key, field)

    if (value === undefined || value === null) return null

    if (parseJson) {
      try {
        return JSON.parse(value) as T
      } catch (error) {
        console.error(`Failed to parse JSON for hash ${key}:${field}:`, error)
        return null
      }
    }

    return value as T
  }

  /**
   * Get all fields and values from a hash
   */
  async getHashAll<T = Record<string, string>>(key: string, parseJson = false): Promise<T | null> {
    const client = await this.ensureConnection()
    const hash = await client.hGetAll(key)

    if (!hash || Object.keys(hash).length === 0) return null

    if (parseJson) {
      const parsed: Record<string, unknown> = {}
      for (const [field, value] of Object.entries(hash)) {
        try {
          parsed[field] = JSON.parse(value)
        } catch (error) {
          console.error(`Failed to parse JSON for hash ${key}:${field}:`, error)
          parsed[field] = value
        }
      }
      return parsed as T
    }

    return hash as T
  }

  /**
   * Delete a field from a hash
   */
  async deleteHashField(key: string, field: string): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.hDel(key, field)
    return result > 0
  }

  /**
   * Get all field names from a hash
   */
  async getHashFields(key: string): Promise<string[]> {
    const client = await this.ensureConnection()
    return await client.hKeys(key)
  }

  /**
   * Set expiration on a key
   */
  async expire(key: string, seconds: number): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.expire(key, seconds)
    return result === 1
  }

  /**
   * Remove expiration from a key
   */
  async persist(key: string): Promise<boolean> {
    const client = await this.ensureConnection()
    const result = await client.persist(key)
    return result === 1
  }

  /**
   * Get time to live for a key
   */
  async getTTL(key: string): Promise<number> {
    const client = await this.ensureConnection()
    return await client.ttl(key)
  }

  /**
   * Increment a numeric value
   */
  async increment(key: string, by = 1): Promise<number> {
    const client = await this.ensureConnection()
    return await client.incrBy(key, by)
  }

  /**
   * Decrement a numeric value
   */
  async decrement(key: string, by = 1): Promise<number> {
    const client = await this.ensureConnection()
    return await client.decrBy(key, by)
  }

  /**
   * Execute multiple operations atomically
   */
  async multi(operations: Array<() => Promise<unknown>>): Promise<unknown[]> {
    const client = await this.ensureConnection()
    const multi = client.multi()

    for (const operation of operations) {
      await operation()
    }

    return await multi.exec()
  }
}

export const redisCache = new RedisCache()
