# Docker Setup Guide

This guide covers how to run the Discord Anime Bot using Docker with Redis for caching.

## Docker Components

### Services

- **app**: The Discord bot application (Bun/TypeScript)
- **redis**: Redis cache server for notifications and watchlists

### Files

- `Dockerfile`: Production optimized build
- `docker-compose.yml`: Deployment with Redis

## Quick Start

1. **Prepare environment**:

   ```bash
   # Copy and configure environment
   cp .env.example .env
   # Edit .env with your Discord bot token and other config
   ```

2. **Start services**:

   ```bash
   # Build and start in detached mode
   bun run docker:up

   # View logs
   bun run docker:logs
   ```

3. **Stop services**:
   ```bash
   bun run docker:down
   ```

## Configuration

### Environment Variables

The bot requires these environment variables:

```bash
# Required
DISCORD_BOT_TOKEN=your_bot_token
CHANNEL_ID=your_channel_id
ANILIST_API=https://graphql.anilist.co

# Redis (auto-configured in Docker)
REDIS_URL=redis://redis:6379

# Optional AI features
OPENAI_API_KEY=your_openai_key
CLAUDE_API_KEY=your_claude_key

# Environment
BUN_ENV=production
PORT=8081
```

### Redis Configuration

The Redis service is configured with:

- **Memory limit**: 512MB
- **Persistence**: AOF (Append Only File) enabled
- **Eviction**: LRU (Least Recently Used) when memory is full
- **Health checks**: Ping every 10s with retries

## Monitoring & Management

### View Logs

```bash
# App logs
bun run docker:logs

# All services
docker-compose logs -f

# Redis logs only
docker-compose logs -f redis
```

### Redis Management

```bash
# Connect to Redis CLI
bun run docker:redis-cli

# Example Redis commands:
# KEYS notification:*     # List all notifications
# KEYS watchlist:*        # List all watchlists
# INFO memory             # Memory usage
# FLUSHDB                 # Clear database (be careful!)
```

## Data Persistence

### Redis Data

- **Storage**: `redis-data` volume persists Redis data
- **Location**: Docker managed volumes

### Backup Redis Data

```bash
# Create backup
docker-compose exec redis redis-cli BGSAVE

# Copy backup file out
docker cp $(docker-compose ps -q redis):/data/dump.rdb ./redis-backup.rdb
```

### Restore Redis Data

```bash
# Stop services
docker-compose down

# Copy backup file in
docker run --rm -v discord-anime-bot_redis-data:/data -v $(pwd):/backup alpine cp /backup/redis-backup.rdb /data/dump.rdb

# Start services
docker-compose up -d
```

## Debugging

```bash
# Shell into app container
docker-compose exec app /bin/bash

# Shell into Redis container
docker-compose exec redis /bin/sh

# View app container details
docker-compose exec app bun --version
```

## Health Checks

### Redis Health Check

- **Command**: `redis-cli ping`
- **Interval**: 10s
- **Timeout**: 5s
- **Retries**: 5

### App Dependencies

The app service waits for Redis to be healthy before starting via `depends_on` with `condition: service_healthy`.

## Scaling

### Horizontal Scaling

```bash
# Scale app service to 3 replicas
docker-compose up -d --scale app=3
```

**Note**: Multiple app instances will share the same Redis cache, which is perfect for notifications and watchlists.

### Resource Limits

Add to `docker-compose.yml` under each service:

```yaml
deploy:
  resources:
    limits:
      memory: 512M
      cpus: '0.5'
    reservations:
      memory: 256M
      cpus: '0.25'
```

## Security

### Production Security

- Non-root user in container (`anime` user)
- No unnecessary packages in final image
- Environment variables for secrets
- Network isolation via custom bridge network
- Redis is not exposed externally by default

### Network Security

```yaml
# Only expose app port externally
ports:
  - '8081:8081' # App only

# Redis only accessible internally
# Remove Redis ports section to make it internal-only
```

## Troubleshooting

### Common Issues

1. **Redis connection failed**:

   ```bash
   # Check Redis health
   docker-compose ps

   # Check Redis logs
   docker-compose logs redis

   # Test connection manually
   bun run docker:redis-cli
   ```

2. **App won't start**:

   ```bash
   # Check app logs
   bun run docker:logs

   # Verify environment variables
   docker-compose exec app env | grep -E "(DISCORD|REDIS)"
   ```

3. **Permission issues**:

   ```bash
   # Rebuild with no cache
   docker-compose build --no-cache
   ```

4. **Port conflicts**:

   ```bash
   # Check what's using ports
   lsof -i :6379
   lsof -i :8081

   # Change ports in docker-compose.yml
   ```

### Clean Slate Reset

```bash
# Stop and remove everything
docker-compose down -v --remove-orphans

# Remove images
docker-compose down --rmi all

# Rebuild from scratch
docker-compose up --build
```

## Available Scripts

| Script                     | Description                  |
| -------------------------- | ---------------------------- |
| `bun run docker:build`     | Build images                 |
| `bun run docker:up`        | Start services in background |
| `bun run docker:down`      | Stop services                |
| `bun run docker:logs`      | View app logs                |
| `bun run docker:redis-cli` | Connect to Redis CLI         |
