#!/usr/bin/env bun

import { redisConnection, redisCache } from '../src/services/redis'

async function testRedis() {
  console.log('Testing Redis connection and operations...\n')

  try {
    // Test connection
    console.log('1. Testing connection...')
    await redisConnection.connect()
    console.log('Redis connection successful\n')

    // Test basic operations
    console.log('2. Testing basic key-value operations...')
    await redisCache.set('test:key', 'test-value')
    const value = await redisCache.get('test:key')
    console.log(`Set and get: ${value}\n`)

    // Test JSON operations
    console.log('3. Testing JSON operations...')
    const testObj = { name: 'Test Anime', id: 12345 }
    await redisCache.set('test:json', testObj)
    const retrievedObj = await redisCache.get('test:json', true)
    console.log(`JSON set and get:`, retrievedObj)
    console.log()

    // Test set operations (for watchlists)
    console.log('4. Testing set operations (watchlist simulation)...')
    await redisCache.addToSet('test:watchlist:user123', 123)
    await redisCache.addToSet('test:watchlist:user123', 456)
    await redisCache.addToSet('test:watchlist:user123', 789)
    
    const watchlistItems = await redisCache.getSetMembers('test:watchlist:user123')
    console.log(`Watchlist items: [${watchlistItems.join(', ')}]`)
    
    const isInWatchlist = await redisCache.isInSet('test:watchlist:user123', 456)
    console.log(`Is 456 in watchlist: ${isInWatchlist}`)
    
    await redisCache.removeFromSet('test:watchlist:user123', 456)
    const updatedWatchlist = await redisCache.getSetMembers('test:watchlist:user123')
    console.log(`After removal: [${updatedWatchlist.join(', ')}]\n`)

    // Test TTL
    console.log('5. Testing TTL operations...')
    await redisCache.set('test:ttl', 'expires-soon', 5)
    const ttl = await redisCache.getTTL('test:ttl')
    console.log(`TTL for test key: ${ttl} seconds\n`)

    // Cleanup
    console.log('6. Cleaning up test data...')
    await redisCache.delete('test:key')
    await redisCache.delete('test:json')
    await redisCache.delete('test:watchlist:user123')
    await redisCache.delete('test:ttl')
    console.log('Cleanup complete\n')

    console.log('All Redis tests passed!')
  } catch (error) {
    console.error('Redis test failed:', error)
  } finally {
    // Disconnect
    await redisConnection.disconnect()
    console.log('Disconnected from Redis')
    process.exit(0)
  }
}

testRedis()