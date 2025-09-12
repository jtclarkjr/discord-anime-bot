# Deployment Guide

This guide covers various deployment options for the Discord Anime Bot, with a focus on cloud platforms.

## Prerequisites

Before deploying, ensure you have:

- Discord Bot Token
- AniList API access (https://graphql.anilist.co)
- OpenAI API Key (optional, for AI features)
- Claude API Key (optional, alternative to OpenAI)
- Redis instance (cloud or self-hosted)

## Environment Variables

All deployments require these environment variables:

### Required
- `DISCORD_BOT_TOKEN` - Your Discord bot token
- `CHANNEL_ID` - Discord channel ID where the bot operates
- `ANILIST_API` - AniList GraphQL endpoint (https://graphql.anilist.co)
- `REDIS_URL` - Redis connection string

### Optional
- `OPENAI_API_KEY` - OpenAI API key for AI-powered search
- `CLAUDE_API_KEY` - Claude API key (alternative to OpenAI)

## Google Cloud Run Deployment

Google Cloud Run is ideal for this Discord bot because it:
- Scales to zero when not in use (cost-effective)
- Handles HTTP and persistent connections
- Supports WebSocket connections for Discord
- Integrates well with other GCP services

### Prerequisites for GCP

1. **Google Cloud Account** with billing enabled
2. **Google Cloud CLI** installed and configured
3. **Docker** installed locally
4. **Redis instance** (Google Memory Store recommended)

### Step 1: Set up Redis

#### Option A: Google Memory Store (Recommended)

```bash
# Create Redis instance
gcloud redis instances create discord-bot-redis \
    --size=1 \
    --region=us-central1 \
    --redis-version=redis_7_0 \
    --network=default \
    --enable-auth

# Get connection details
gcloud redis instances describe discord-bot-redis --region=us-central1
```

#### Option B: External Redis Provider

Use services like:
- Redis Cloud
- Upstash
- ElastiCache (if using AWS resources)

### Step 2: Prepare the Application

1. **Create production Dockerfile** (if not using existing):
   ```dockerfile
   FROM oven/bun:1.2.20-slim AS base
   WORKDIR /app

   FROM base AS build
   RUN apt-get update -qq && \
       apt-get install --no-install-recommends -y build-essential pkg-config python-is-python3
   COPY bun.lock package.json ./
   RUN bun install --ci
   COPY . .

   FROM base
   COPY --from=build /app /app
   RUN groupadd -r anime && useradd -r -g anime anime
   RUN chown -R anime:anime /app
   USER anime

   EXPOSE 8080
   CMD ["bun", "run", "start"]
   ```

2. **Create `.gcloudignore`**:
   ```
   .git
   .gitignore
   README.md
   Dockerfile*
   docker-compose*.yml
   .env*
   node_modules
   data/
   scripts/
   go/
   ```

### Step 3: Configure Environment

Create a production environment file:

```bash
# Create .env.production
cat > .env.production << EOF
DISCORD_BOT_TOKEN=your_bot_token_here
CHANNEL_ID=your_channel_id_here
ANILIST_API=https://graphql.anilist.co
REDIS_URL=redis://[USERNAME]:[PASSWORD]@[HOST]:[PORT]
OPENAI_API_KEY=your_openai_key_here
CLAUDE_API_KEY=your_claude_key_here
PORT=8080
EOF
```

### Step 4: Deploy to Cloud Run

1. **Set up Google Cloud project**:
   ```bash
   # Set project
   gcloud config set project YOUR_PROJECT_ID
   
   # Enable required APIs
   gcloud services enable run.googleapis.com
   gcloud services enable cloudbuild.googleapis.com
   gcloud services enable redis.googleapis.com
   ```

2. **Build and deploy**:
   ```bash
   # Build and deploy in one command
   gcloud run deploy discord-anime-bot \
     --source . \
     --region=us-central1 \
     --platform=managed \
     --allow-unauthenticated \
     --port=8080 \
     --memory=512Mi \
     --cpu=1 \
     --min-instances=0 \
     --max-instances=10 \
     --env-vars-file=.env.production \
     --execution-environment=gen2
   ```

3. **Alternative: Deploy with pre-built image**:
   ```bash
   # Build image locally
   docker build -t gcr.io/YOUR_PROJECT_ID/discord-anime-bot .
   
   # Push to Container Registry
   docker push gcr.io/YOUR_PROJECT_ID/discord-anime-bot
   
   # Deploy from registry
   gcloud run deploy discord-anime-bot \
     --image=gcr.io/YOUR_PROJECT_ID/discord-anime-bot \
     --region=us-central1 \
     --platform=managed \
     --allow-unauthenticated \
     --port=8080 \
     --memory=512Mi \
     --cpu=1 \
     --min-instances=1 \
     --max-instances=5 \
     --set-env-vars="$(cat .env.production | tr '\n' ',' | sed 's/,$//')"
   ```

### Step 5: Configure Networking (for Memory Store)

If using Google Memory Store, configure VPC access:

```bash
# Create VPC connector
gcloud compute networks vpc-access connectors create discord-bot-connector \
  --region=us-central1 \
  --subnet=default \
  --subnet-project=YOUR_PROJECT_ID

# Update Cloud Run service to use VPC connector
gcloud run services update discord-anime-bot \
  --region=us-central1 \
  --vpc-connector=discord-bot-connector \
  --vpc-egress=private-ranges-only
```

### Step 6: Monitor and Maintain

1. **View logs**:
   ```bash
   gcloud run services logs tail discord-anime-bot --region=us-central1
   ```

2. **Update deployment**:
   ```bash
   # Redeploy with changes
   gcloud run deploy discord-anime-bot \
     --source . \
     --region=us-central1
   ```

3. **Set up monitoring**:
   ```bash
   # Enable Cloud Monitoring
   gcloud services enable monitoring.googleapis.com
   ```

## Fly.io Deployment

Fly.io is an excellent choice for Discord bots because it:
- Deploys globally in 35+ regions
- Scales instantly with hardware-virtualized containers
- Offers competitive pricing with generous free tier
- Provides built-in Redis (Upstash integration)
- Supports persistent volumes and networking

### Prerequisites for Fly.io

1. **Fly.io account** (free tier available)
2. **Fly CLI** installed
3. **Docker** (optional, Fly can build from source)

### Step 1: Install Fly CLI

```bash
# macOS
brew install flyctl

# Linux
curl -L https://fly.io/install.sh | sh

# Windows
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"

# Login
flyctl auth login
```

### Step 2: Set up Redis

#### Option A: Fly.io Redis (Upstash Integration)

```bash
# Create Redis instance
flyctl redis create
# Follow prompts to choose region and plan

# Get connection string
flyctl redis status <redis-name>
```

#### Option B: External Redis Provider

Use Redis Cloud, Upstash directly, or other providers.

### Step 3: Initialize Fly App

```bash
# Initialize Fly app
flyctl launch

# This creates fly.toml configuration file
# Choose:
# - App name: discord-anime-bot-yourname
# - Region: Choose closest to your users
# - PostgreSQL: No
# - Deploy now: No (we'll configure first)
```

### Step 4: Configure fly.toml

Edit the generated `fly.toml` file:

```toml
app = "discord-anime-bot-yourname"
primary_region = "sea" # or your preferred region

[build]
  dockerfile = "Dockerfile"

[env]
  PORT = "8080"
  ANILIST_API = "https://graphql.anilist.co"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
  processes = ["app"]

[[vm]]
  cpu_kind = "shared"
  cpus = 1
  memory_mb = 512

[deploy]
  release_command = "echo 'Starting Discord Anime Bot...'"
```

### Step 5: Set Secrets

```bash
# Set environment variables as secrets
flyctl secrets set DISCORD_BOT_TOKEN="your_bot_token_here"
flyctl secrets set CHANNEL_ID="your_channel_id_here"
flyctl secrets set REDIS_URL="redis://your_redis_connection_string"
flyctl secrets set OPENAI_API_KEY="your_openai_key_here"
flyctl secrets set CLAUDE_API_KEY="your_claude_key_here"

# Verify secrets
flyctl secrets list
```

### Step 6: Create Production Dockerfile

Create a Fly.io optimized Dockerfile:

```dockerfile
# Dockerfile.fly
FROM oven/bun:1.2.20-slim AS base

# Bun app lives here
WORKDIR /app

# Set production environment
ENV BUN_ENV="production"

# Install packages needed to build node modules
RUN apt-get update -qq && \
    apt-get install --no-install-recommends -y build-essential pkg-config python-is-python3 && \
    rm -rf /var/lib/apt/lists /var/cache/apt/archives

# Install node modules
COPY bun.lock package.json ./
RUN bun install --ci

# Copy application code
COPY . .

# Create non-root user for security
RUN groupadd -r anime && useradd -r -g anime anime
RUN chown -R anime:anime /app
USER anime

# Start the server by default
EXPOSE 8080
CMD ["bun", "run", "start"]
```

Update `fly.toml` to use this Dockerfile:

```toml
[build]
  dockerfile = "Dockerfile.fly"
```

### Step 7: Deploy

```bash
# Deploy the application
flyctl deploy

# Monitor deployment
flyctl logs

# Check app status
flyctl status
```

### Step 8: Configure Scaling (Optional)

```bash
# Scale to multiple regions for global deployment
flyctl scale count 2 --region sea,lax

# Scale resources if needed
flyctl scale memory 1024  # 1GB RAM
flyctl scale cpu dedicated # Dedicated CPU

# Set minimum machines running (avoid cold starts)
flyctl scale count 1 --min-machines-running=1
```

### Step 9: Monitor and Maintain

1. **View logs**:
   ```bash
   flyctl logs
   flyctl logs --app discord-anime-bot-yourname
   ```

2. **Access app console**:
   ```bash
   flyctl ssh console
   ```

3. **Update deployment**:
   ```bash
   flyctl deploy
   ```

4. **Monitor resources**:
   ```bash
   flyctl status
   flyctl vm status
   ```

### Fly.io Specific Benefits

- **Global Edge Network**: Deploy close to Discord's servers
- **Instant Scaling**: Machines start in milliseconds
- **Cost Effective**: Pay only for resources used
- **Built-in Monitoring**: Integrated metrics and logs
- **Zero Downtime Deploys**: Rolling updates

### Advanced Fly.io Configuration

#### Multi-Region Deployment

```bash
# Deploy to multiple regions
flyctl regions add sea lax fra
flyctl scale count 3
```

#### Persistent Volumes (if needed)

```toml
# Add to fly.toml
[[mounts]]
  source = "discord_bot_data"
  destination = "/app/data"
```

```bash
# Create volume
flyctl volumes create discord_bot_data --region sea
```

#### Health Checks

```toml
# Add to fly.toml
[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.checks]]
    grace_period = "10s"
    interval = "30s"
    method = "GET"
    path = "/health"
    timeout = "5s"
```

Add health endpoint to your app:

```typescript
// Add to your Express/HTTP server
app.get('/health', (req, res) => {
  res.status(200).json({ status: 'healthy', timestamp: new Date().toISOString() })
})
```

## Other Deployment Options

### Docker Compose (Self-hosted)

For VPS or dedicated servers:

```bash
# Clone repository
git clone <your-repo>
cd discord-anime-bot

# Configure environment
cp .env.docker .env
# Edit .env with your configuration

# Deploy
docker-compose up --build -d

# Monitor
docker-compose logs -f app
```

### Railway

1. Connect your GitHub repository to Railway
2. Set environment variables in Railway dashboard
3. Deploy automatically on push

Required environment variables:
- All variables from prerequisites section
- `PORT=8080` (Railway auto-assigns)

### Heroku

1. **Install Heroku CLI**
2. **Create Heroku app**:
   ```bash
   heroku create discord-anime-bot-yourname
   ```

3. **Add Redis addon**:
   ```bash
   heroku addons:create heroku-redis:mini
   ```

4. **Set environment variables**:
   ```bash
   heroku config:set DISCORD_BOT_TOKEN=your_token
   heroku config:set CHANNEL_ID=your_channel_id
   heroku config:set ANILIST_API=https://graphql.anilist.co
   heroku config:set OPENAI_API_KEY=your_openai_key
   ```

5. **Deploy**:
   ```bash
   git push heroku main
   ```

### DigitalOcean App Platform

1. Connect GitHub repository
2. Configure environment variables
3. Set build command: `bun install`
4. Set run command: `bun run start`
5. Add DigitalOcean Managed Redis database

## Production Considerations

### Security

- Use strong Redis authentication
- Rotate API keys regularly
- Use VPC/private networks when possible
- Enable audit logging
- Set up monitoring and alerting

### Performance

- **Memory**: 512MB minimum recommended
- **CPU**: 1 vCPU sufficient for most servers
- **Scaling**: Set minimum instances to 1 to avoid cold starts
- **Redis**: Use Redis clusters for high availability

### Monitoring

Essential metrics to monitor:
- Discord API rate limits
- Redis connection health
- Memory usage
- Response times
- Error rates

### Backup Strategy

- **Configuration**: Keep environment variables in secure storage
- **Code**: Use Git tags for releases
- **Redis Data**: 
  - Configure Redis persistence (AOF + RDB)
  - Regular backups to cloud storage
  - Test restore procedures

## Troubleshooting

### Common Issues

1. **Discord connection fails**:
   - Check bot token validity
   - Verify bot permissions in Discord server
   - Check network connectivity

2. **Redis connection fails**:
   - Verify Redis URL format
   - Check network/VPC configuration
   - Confirm Redis instance is running

3. **Memory issues**:
   - Increase memory allocation
   - Monitor Redis memory usage
   - Check for memory leaks in logs

4. **Cold start delays** (Cloud Run):
   - Set minimum instances to 1
   - Use gen2 execution environment
   - Optimize startup time

### Getting Help

- Check application logs first
- Monitor Discord API status
- Verify all environment variables
- Test Redis connectivity separately
- Check cloud provider status pages

## Cost Optimization

### Google Cloud Run
- Uses pay-per-request model
- Scales to zero when inactive
- Estimated cost: $5-20/month for typical usage

### Fly.io
- Generous free tier (3 shared-cpu-1x machines)
- Pay-per-use pricing model
- Estimated cost: $0-15/month for typical usage
- Global deployment included

### Redis Options
- **Memory Store**: $30-100/month depending on size
- **Redis Cloud**: $0-30/month with free tier
- **Fly.io Redis (Upstash)**: $0-25/month with free tier
- **Self-hosted**: VPS costs + maintenance time

### Tips
- Use minimum instance counts wisely
- Monitor usage patterns
- Set up billing alerts
- Use resource quotas