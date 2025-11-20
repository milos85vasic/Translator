# Storage and Progress Tracking

This document describes the enhanced storage backends and comprehensive progress tracking system in Universal Ebook Translator v2.1+.

## 📋 Table of Contents

- [Storage Backends](#storage-backends)
- [Progress Tracking](#progress-tracking)
- [WebSocket Events](#websocket-events)
- [Session Management](#session-management)
- [Translation Cache](#translation-cache)
- [Statistics](#statistics)

## 💾 Storage Backends

The system supports three storage backends that can be used individually or in combination:

### 1. PostgreSQL (Production Recommended)

**Best for:** Production deployments, high reliability, ACID compliance

**Features:**
- Full ACID transactions
- Concurrent access support
- Advanced indexing
- Query optimization
- Backup and recovery

**Configuration:**
```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=translator
DB_PASSWORD=secure_password
DB_NAME=translator
DB_SSLMODE=require
```

**Connection Pool Settings:**
```go
MaxOpenConns: 25        // Maximum open connections
MaxIdleConns: 5         // Maximum idle connections
ConnMaxLifetime: 5min   // Maximum connection lifetime
```

### 2. SQLite with SQLCipher (Standalone)

**Best for:** Development, single-user deployments, embedded systems

**Features:**
- Zero configuration
- File-based storage
- SQLCipher encryption
- Lightweight
- No separate server needed

**Configuration:**
```env
DB_TYPE=sqlite
SQLITE_PATH=./data/translator.db
SQLITE_ENCRYPTION_KEY=your-32-character-encryption-key
```

**Encryption:**
```bash
# Encrypted database using SQLCipher
# AES-256 encryption
# 4096-byte page size
```

### 3. Redis (High-Performance Caching)

**Best for:** High-speed caching, temporary sessions, real-time data

**Features:**
- In-memory storage
- Sub-millisecond latency
- Automatic TTL expiration
- Pub/Sub messaging
- High throughput

**Configuration:**
```env
DB_TYPE=redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis_password
REDIS_TTL=24h
```

### Hybrid Configuration (Recommended)

Use PostgreSQL for persistence and Redis for caching:

```go
// Primary storage: PostgreSQL
primaryStorage := storage.NewPostgreSQLStorage(pgConfig)

// Cache layer: Redis
cacheStorage := storage.NewRedisStorage(redisConfig, 24*time.Hour)

// Use cache for reads, PostgreSQL for writes
```

## 📊 Progress Tracking

### Overview

The enhanced progress tracker provides real-time detailed information about translation progress.

**Tracked Metrics:**
- Percentage complete
- Current chapter and section
- Elapsed time
- Estimated time to completion (ETA)
- Items completed/failed/total
- Book and translation metadata

### Progress Data Structure

```go
type TranslationProgress struct {
    // Book information
    BookTitle       string    `json:"book_title"`
    TotalChapters   int       `json:"total_chapters"`
    CurrentChapter  int       `json:"current_chapter"`
    ChapterTitle    string    `json:"chapter_title"`
    CurrentSection  int       `json:"current_section"`
    TotalSections   int       `json:"total_sections"`

    // Progress metrics
    PercentComplete float64   `json:"percent_complete"`
    ItemsTotal      int       `json:"items_total"`
    ItemsCompleted  int       `json:"items_completed"`
    ItemsFailed     int       `json:"items_failed"`

    // Time tracking
    StartTime       time.Time `json:"start_time"`
    ElapsedTime     string    `json:"elapsed_time"`        // "5 minutes 32 seconds"
    EstimatedETA    string    `json:"estimated_eta"`       // "15 minutes 48 seconds"

    // Translation details
    SourceLanguage  string    `json:"source_language"`
    TargetLanguage  string    `json:"target_language"`
    Provider        string    `json:"provider"`
    Model           string    `json:"model"`

    // Status
    Status          string    `json:"status"`              // "translating", "completed", etc.
    CurrentTask     string    `json:"current_task"`        // Human-readable current action
    SessionID       string    `json:"session_id"`
}
```

### Usage Example

```go
// Create tracker
tracker := progress.NewTracker(
    sessionID,
    "Crime and Punishment",
    36,  // total chapters
    "ru", // source language
    "sr", // target language
    "deepseek",
    "deepseek-chat",
)

// Update chapter progress
tracker.UpdateChapter(5, "Part 1, Chapter 5", 12) // chapter 5, 12 sections

// Update section
tracker.UpdateSection(3) // currently on section 3

// Increment counters
tracker.IncrementCompleted()
tracker.IncrementFailed()

// Get current progress
progress := tracker.GetProgress()
fmt.Printf("Progress: %.2f%%\n", progress.PercentComplete)
fmt.Printf("Elapsed: %s\n", progress.ElapsedTime)
fmt.Printf("ETA: %s\n", progress.EstimatedETA)

// Mark complete
tracker.Complete()
```

### CLI Progress Display

```
Universal Ebook Translator v2.1.0

Book: Crime and Punishment
Chapters: 36
Progress: 45.8% (Chapter 16/36: "Part 2, Chapter 8")

Status: Translating
Elapsed Time: 25 minutes 14 seconds
Estimated Time Remaining: 29 minutes 32 seconds

Provider: DeepSeek (deepseek-chat)
Source: Russian → Target: Serbian Cyrillic

Items: 1,248 completed, 23 failed, 2,723 total
```

## 🔌 WebSocket Events

### Event Types

The system emits comprehensive WebSocket events for real-time progress tracking:

#### 1. Translation Started

```json
{
  "type": "translation_started",
  "session_id": "sess_123abc",
  "timestamp": "2025-11-20T16:00:00Z",
  "data": {
    "book_title": "Crime and Punishment",
    "total_chapters": 36,
    "source_language": "Russian",
    "target_language": "Serbian",
    "provider": "deepseek",
    "model": "deepseek-chat"
  }
}
```

#### 2. Progress Update

```json
{
  "type": "translation_progress",
  "session_id": "sess_123abc",
  "timestamp": "2025-11-20T16:05:30Z",
  "data": {
    "book_title": "Crime and Punishment",
    "current_chapter": 5,
    "total_chapters": 36,
    "chapter_title": "Part 1, Chapter 5",
    "current_section": 3,
    "total_sections": 12,
    "percent_complete": 13.8,
    "items_completed": 376,
    "items_failed": 5,
    "items_total": 2723,
    "elapsed_time": "5 minutes 30 seconds",
    "estimated_eta": "34 minutes 15 seconds",
    "source_language": "Russian",
    "target_language": "Serbian",
    "provider": "deepseek",
    "model": "deepseek-chat",
    "status": "translating",
    "current_task": "Translating chapter Part 1, Chapter 5"
  }
}
```

#### 3. Translation Completed

```json
{
  "type": "translation_completed",
  "session_id": "sess_123abc",
  "timestamp": "2025-11-20T16:42:18Z",
  "data": {
    "book_title": "Crime and Punishment",
    "total_chapters": 36,
    "items_completed": 2700,
    "items_failed": 23,
    "items_total": 2723,
    "elapsed_time": "42 minutes 18 seconds",
    "output_file": "Crime_and_Punishment_sr.epub"
  }
}
```

#### 4. Translation Error

```json
{
  "type": "translation_error",
  "session_id": "sess_123abc",
  "timestamp": "2025-11-20T16:15:42Z",
  "data": {
    "error": "API rate limit exceeded",
    "chapter": 8,
    "section": 5,
    "recoverable": true
  }
}
```

### WebSocket Client Example

**JavaScript:**
```javascript
const ws = new WebSocket('wss://localhost:8443/ws?session_id=sess_123abc');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);

  switch (data.type) {
    case 'translation_progress':
      console.log(`Progress: ${data.data.percent_complete}%`);
      console.log(`Chapter: ${data.data.current_chapter}/${data.data.total_chapters}`);
      console.log(`ETA: ${data.data.estimated_eta}`);
      break;

    case 'translation_completed':
      console.log('Translation completed!');
      console.log(`Output: ${data.data.output_file}`);
      break;

    case 'translation_error':
      console.error('Error:', data.data.error);
      break;
  }
};
```

**Python:**
```python
import websockets
import json
import asyncio

async def monitor_translation(session_id):
    uri = f"wss://localhost:8443/ws?session_id={session_id}"

    async with websockets.connect(uri, ssl=ssl_context) as websocket:
        async for message in websocket:
            data = json.loads(message)

            if data['type'] == 'translation_progress':
                progress = data['data']
                print(f"Progress: {progress['percent_complete']:.1f}%")
                print(f"Chapter: {progress['current_chapter']}/{progress['total_chapters']}")
                print(f"ETA: {progress['estimated_eta']}")

asyncio.run(monitor_translation("sess_123abc"))
```

## 💼 Session Management

### Session Lifecycle

```
┌──────────────┐
│  Initialize  │
│   Session    │
└──────┬───────┘
       │
       v
┌──────────────┐
│  Translating │◄─┐
│              │  │ Retry on error
└──────┬───────┘  │
       │          │
       │  Error? ─┘
       v
┌──────────────┐
│   Completed  │
│  or Failed   │
└──────────────┘
```

### Session Storage

**Create Session:**
```go
session := &storage.TranslationSession{
    ID:             sessionID,
    BookTitle:      "Crime and Punishment",
    InputFile:      "book.epub",
    OutputFile:     "book_sr.epub",
    SourceLanguage: "ru",
    TargetLanguage: "sr",
    Provider:       "deepseek",
    Model:          "deepseek-chat",
    Status:         "initializing",
    TotalChapters:  36,
    StartTime:      time.Now(),
    CreatedAt:      time.Now(),
    UpdatedAt:      time.Now(),
}

err := storage.CreateSession(ctx, session)
```

**Update Session:**
```go
session.Status = "translating"
session.CurrentChapter = 5
session.PercentComplete = 13.8
session.ItemsCompleted = 376

err := storage.UpdateSession(ctx, session)
```

**Retrieve Session:**
```go
session, err := storage.GetSession(ctx, sessionID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", session.Status)
fmt.Printf("Progress: %.1f%%\n", session.PercentComplete)
```

### REST API Endpoints

```bash
# Get session status
curl -k https://localhost:8443/api/v1/sessions/{session_id}

# List all sessions
curl -k https://localhost:8443/api/v1/sessions?limit=10&offset=0

# Get session statistics
curl -k https://localhost:8443/api/v1/sessions/{session_id}/stats
```

## 🗄️ Translation Cache

### Cache Structure

```go
type TranslationCache struct {
    ID              string
    SourceText      string
    TargetText      string
    SourceLanguage  string
    TargetLanguage  string
    Provider        string
    Model           string
    CreatedAt       time.Time
    AccessCount     int
    LastAccessedAt  time.Time
}
```

### Cache Operations

**Store Translation:**
```go
cache := &storage.TranslationCache{
    ID:             cacheID,
    SourceText:     "Hello, world!",
    TargetText:     "Здраво, свете!",
    SourceLanguage: "en",
    TargetLanguage: "sr",
    Provider:       "deepseek",
    Model:          "deepseek-chat",
    CreatedAt:      time.Now(),
    LastAccessedAt: time.Now(),
}

err := storage.CacheTranslation(ctx, cache)
```

**Retrieve Cached Translation:**
```go
cached, err := storage.GetCachedTranslation(
    ctx,
    "Hello, world!",
    "en", // source
    "sr", // target
    "deepseek",
    "deepseek-chat",
)

if cached != nil {
    fmt.Println("Cache hit:", cached.TargetText)
} else {
    fmt.Println("Cache miss - translating...")
}
```

**Cache Cleanup:**
```go
// Remove entries not accessed in 30 days
err := storage.CleanupOldCache(ctx, 30*24*time.Hour)
```

## 📈 Statistics

### Available Statistics

```go
type Statistics struct {
    TotalSessions      int64
    CompletedSessions  int64
    FailedSessions     int64
    InProgressSessions int64
    TotalTranslations  int64
    CacheHitRate       float64
    AverageDuration    float64
}
```

### Get Statistics

```go
stats, err := storage.GetStatistics(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total Sessions: %d\n", stats.TotalSessions)
fmt.Printf("Completed: %d\n", stats.CompletedSessions)
fmt.Printf("Failed: %d\n", stats.FailedSessions)
fmt.Printf("In Progress: %d\n", stats.InProgressSessions)
fmt.Printf("Cache Hit Rate: %.2f%%\n", stats.CacheHitRate)
fmt.Printf("Average Duration: %.2f seconds\n", stats.AverageDuration)
```

### REST API

```bash
# Get overall statistics
curl -k https://localhost:8443/api/v1/stats

# Response:
{
  "total_sessions": 150,
  "completed_sessions": 142,
  "failed_sessions": 5,
  "in_progress_sessions": 3,
  "total_translations": 125847,
  "cache_hit_rate": 23.5,
  "average_duration_seconds": 1842.3
}
```

## 🔍 Performance Considerations

### PostgreSQL Optimization

```sql
-- Create indexes for faster queries
CREATE INDEX idx_sessions_status ON translation_sessions(status);
CREATE INDEX idx_sessions_created_at ON translation_sessions(created_at DESC);
CREATE INDEX idx_cache_lookup ON translation_cache(source_text, source_language, target_language);
```

### Redis Configuration

```
# Redis tuning for caching
maxmemory 512mb
maxmemory-policy allkeys-lru
save ""  # Disable persistence for pure caching
```

### Cache Strategy

```go
// Try Redis first (fast)
cached, _ := redisStorage.GetCachedTranslation(ctx, text, src, tgt, provider, model)
if cached != nil {
    return cached.TargetText, nil
}

// Fall back to PostgreSQL
cached, _ = pgStorage.GetCachedTranslation(ctx, text, src, tgt, provider, model)
if cached != nil {
    // Warm Redis cache
    _ = redisStorage.CacheTranslation(ctx, cached)
    return cached.TargetText, nil
}

// Translate and cache in both
translation := translate(text)
cache := &TranslationCache{...}
_ = pgStorage.CacheTranslation(ctx, cache)
_ = redisStorage.CacheTranslation(ctx, cache)

return translation, nil
```

## 📚 Additional Resources

- [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md) - Docker setup
- [API.md](API.md) - REST API reference
- [README.md](../README.md) - Main documentation
