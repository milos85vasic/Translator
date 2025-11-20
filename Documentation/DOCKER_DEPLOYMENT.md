# Docker Deployment Guide

This guide covers deploying the Universal Ebook Translator using Docker and Docker Compose with full PostgreSQL, Redis, and production-ready configuration.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Services](#services)
- [Management Scripts](#management-scripts)
- [Storage Backends](#storage-backends)
- [Security](#security)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## 🔧 Prerequisites

- Docker 20.10+ installed
- Docker Compose 2.0+ installed
- Minimum 2GB RAM available
- Minimum 10GB disk space

### Install Docker

**macOS:**
```bash
brew install --cask docker
```

**Ubuntu/Debian:**
```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

**Verify Installation:**
```bash
docker --version
docker-compose --version
```

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd Translate

# Copy environment template
cp .env.example .env

# Edit with your configuration
vim .env
```

### 2. Generate TLS Certificates

```bash
make generate-certs
```

### 3. Start Services

```bash
# Start all services
./scripts/start.sh

# Or with admin tools (Adminer, Redis Commander)
./scripts/start.sh --admin

# Or build and start
./scripts/start.sh --build
```

### 4. Verify Services

```bash
# Check service status
docker-compose ps

# View logs
./scripts/logs.sh -f

# Test API
curl -k https://localhost:8443/health
```

## ⚙️ Configuration

### Environment Variables

Edit `.env` file with your configuration:

```bash
# API Configuration
API_PORT=8443
JWT_SECRET=your-super-secure-secret

# Database
DB_TYPE=postgres
POSTGRES_USER=translator
POSTGRES_PASSWORD=secure_password
POSTGRES_DB=translator

# Redis
REDIS_PASSWORD=redis_password

# LLM API Keys
OPENAI_API_KEY=sk-...
DEEPSEEK_API_KEY=sk-...
```

**Important:** Never commit `.env` file to version control!

### Storage Backend Selection

The system supports three storage backends:

#### 1. PostgreSQL (Recommended for Production)

```env
DB_TYPE=postgres
DB_HOST=postgres
DB_PORT=5432
DB_USER=translator
DB_PASSWORD=secure_password
DB_NAME=translator
DB_SSLMODE=disable
```

#### 2. SQLite with SQLCipher (Standalone/Development)

```env
DB_TYPE=sqlite
SQLITE_PATH=./data/translator.db
SQLITE_ENCRYPTION_KEY=your-encryption-key
```

#### 3. Redis (High-Performance Caching)

```env
DB_TYPE=redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redis_password
```

## 🏗️ Services

### Service Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Docker Network                        │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │              │  │              │  │              │ │
│  │ Translator   │  │  PostgreSQL  │  │    Redis     │ │
│  │   API        │◄─┤   Database   │  │    Cache     │ │
│  │              │  │              │  │              │ │
│  └───────┬──────┘  └──────────────┘  └──────────────┘ │
│          │                                             │
│          │         ┌──────────────┐  ┌──────────────┐ │
│          │         │   Adminer    │  │    Redis     │ │
│          └────────►│  (optional)  │  │  Commander   │ │
│                    │              │  │  (optional)  │ │
│                    └──────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### Individual Services

#### 1. Translator API

**Container:** `translator-api`
**Ports:** 8443 (HTTPS), 8080 (HTTP)
**Features:**
- REST API with HTTP/3 support
- WebSocket real-time events
- Translation progress tracking
- Multiple LLM providers

```bash
# Access API
curl -k https://localhost:8443/health

# View logs
docker logs -f translator-api

# Execute CLI inside container
./scripts/exec.sh translator -input Books/book.epub -locale de
```

#### 2. PostgreSQL Database

**Container:** `translator-postgres`
**Port:** 5432
**Features:**
- Session persistence
- Translation caching
- Statistics tracking

```bash
# Connect to database
./scripts/exec.sh postgres psql -U translator

# Backup database
docker exec translator-postgres pg_dump -U translator translator > backup.sql

# Restore database
cat backup.sql | docker exec -i translator-postgres psql -U translator translator
```

#### 3. Redis Cache

**Container:** `translator-redis`
**Port:** 6379
**Features:**
- High-performance caching
- Automatic TTL expiration
- Session storage

```bash
# Connect to Redis
./scripts/exec.sh redis redis-cli -a redis_password

# Monitor operations
./scripts/exec.sh redis redis-cli -a redis_password MONITOR

# Get cache stats
./scripts/exec.sh redis redis-cli -a redis_password INFO stats
```

#### 4. Adminer (Optional Admin Tool)

**Container:** `translator-adminer`
**Port:** 8081
**Features:**
- Web-based database management
- SQL query execution
- Data export/import

```bash
# Start with admin tools
./scripts/start.sh --admin

# Access Adminer
open http://localhost:8081
```

#### 5. Redis Commander (Optional Admin Tool)

**Container:** `translator-redis-commander`
**Port:** 8082
**Features:**
- Redis key browser
- Value editor
- TTL management

```bash
# Access Redis Commander
open http://localhost:8082
```

## 🛠️ Management Scripts

### start.sh - Start Services

```bash
# Start all services in background
./scripts/start.sh

# Start with admin tools
./scripts/start.sh --admin

# Build images before starting
./scripts/start.sh --build

# Start in foreground (see logs)
./scripts/start.sh --foreground
```

### stop.sh - Stop Services

```bash
# Stop all services (preserve data)
./scripts/stop.sh

# Stop and remove volumes (WARNING: deletes all data!)
./scripts/stop.sh --volumes

# Stop and remove orphaned containers
./scripts/stop.sh --all
```

### restart.sh - Restart Services

```bash
# Restart all services
./scripts/restart.sh

# Restart specific service
./scripts/restart.sh api
./scripts/restart.sh postgres
./scripts/restart.sh redis

# Rebuild and restart
./scripts/restart.sh --build
```

### logs.sh - View Logs

```bash
# Show last 100 lines of all services
./scripts/logs.sh

# Follow API logs in real-time
./scripts/logs.sh -f api

# Show last 500 lines of PostgreSQL
./scripts/logs.sh -t 500 postgres

# Follow all services
./scripts/logs.sh -f all
```

### exec.sh - Execute Commands

```bash
# Run translator CLI
./scripts/exec.sh translator -input Books/book.epub -locale de

# Open shell in API container
./scripts/exec.sh api /bin/sh

# Run PostgreSQL query
./scripts/exec.sh postgres psql -U translator -c "SELECT COUNT(*) FROM translation_sessions;"

# Access Redis CLI
./scripts/exec.sh redis redis-cli -a redis_password
```

## 🔐 Security

### TLS/SSL Configuration

The system uses self-signed certificates by default. For production:

```bash
# Generate production certificates
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=your-domain.com"

# Or use Let's Encrypt
# Place your certificates in:
# - certs/server.crt
# - certs/server.key
```

### API Key Security

```bash
# Set API keys in .env
OPENAI_API_KEY=sk-...
DEEPSEEK_API_KEY=sk-...

# Never commit .env to version control
echo ".env" >> .gitignore
```

### Database Encryption

**SQLCipher for SQLite:**
```env
SQLITE_ENCRYPTION_KEY=your-strong-encryption-key-32-chars
```

**PostgreSQL SSL:**
```env
DB_SSLMODE=require
```

### JWT Authentication

```env
JWT_SECRET=your-very-long-random-secret-key
```

## 📊 Monitoring

### Health Checks

```bash
# API health
curl -k https://localhost:8443/health

# PostgreSQL health
docker exec translator-postgres pg_isready -U translator

# Redis health
docker exec translator-redis redis-cli -a redis_password PING
```

### Statistics

```bash
# Get translation statistics
curl -k https://localhost:8443/api/v1/stats

# Database statistics
./scripts/exec.sh postgres psql -U translator -c "
  SELECT status, COUNT(*)
  FROM translation_sessions
  GROUP BY status;
"

# Redis statistics
./scripts/exec.sh redis redis-cli -a redis_password INFO stats
```

### Resource Usage

```bash
# Container resource usage
docker stats

# Disk usage
docker system df

# Volume usage
docker volume ls
du -sh $(docker volume inspect translator_postgres_data --format '{{.Mountpoint}}')
```

## 🐛 Troubleshooting

### Service Won't Start

```bash
# Check Docker is running
docker info

# Check port conflicts
lsof -i :8443
lsof -i :5432
lsof -i :6379

# View detailed logs
docker-compose logs --tail=100
```

### Database Connection Issues

```bash
# Test PostgreSQL connection
docker exec translator-postgres pg_isready -U translator

# Check PostgreSQL logs
docker logs translator-postgres

# Reset database
./scripts/stop.sh --volumes
./scripts/start.sh
```

### Translation Failures

```bash
# Check API logs
./scripts/logs.sh -f api

# Verify API keys
docker exec translator-api env | grep API_KEY

# Test LLM connection
curl -k https://localhost:8443/api/v1/providers
```

### Performance Issues

```bash
# Check resource usage
docker stats

# Increase memory (docker-compose.yml)
services:
  translator-api:
    deploy:
      resources:
        limits:
          memory: 2G

# Clear Redis cache
./scripts/exec.sh redis redis-cli -a redis_password FLUSHDB
```

### Clean Reset

```bash
# Complete clean reset
./scripts/stop.sh --volumes
docker system prune -a
make clean
./scripts/start.sh --build
```

## 📚 Additional Resources

- [README.md](../README.md) - Main documentation
- [API.md](API.md) - REST API reference
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [V2_RELEASE_NOTES.md](V2_RELEASE_NOTES.md) - Version 2.0 features

## 🤝 Support

For issues and questions:
- GitHub Issues: [Create Issue]
- Documentation: `/Documentation`
- Examples: `/api/examples`
