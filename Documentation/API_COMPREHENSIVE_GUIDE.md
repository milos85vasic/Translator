# Universal Multi-Format Multi-Language Ebook Translation API
## Comprehensive Guide & Reference

---

## 📚 Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Core Endpoints](#core-endpoints)
4. [Translation Endpoints](#translation-endpoints)
5. [Batch Operations](#batch-operations)
6. [File Operations](#file-operations)
7. [Status & Monitoring](#status--monitoring)
8. [Error Handling](#error-handling)
9. [Rate Limiting](#rate-limiting)
10. [WebSocket API](#websocket-api)
11. [Examples & Use Cases](#examples--use-cases)
12. [SDK & Integration](#sdk--integration)

---

## Overview

The Universal Ebook Translation API provides a comprehensive REST interface for translating ebooks and text content between multiple languages and formats. Built with high performance, security, and scalability in mind.

### Base URL
```
Production: https://api.translator.example.com
Development: http://localhost:8080
```

### API Version
```
Current Version: v1
Base Path: /api/v1
```

### Supported Formats
- **Input**: EPUB, FB2, TXT, MD
- **Output**: EPUB, FB2, TXT, HTML
- **Languages**: 100+ languages including major world languages

### Supported Providers
- OpenAI (GPT models)
- Anthropic (Claude models)
- DeepSeek
- Zhipu AI
- Ollama (local models)
- Llama.cpp (local models)

---

## Authentication

### JWT Token Authentication
All API endpoints require JWT authentication (unless explicitly disabled).

#### Request Headers
```http
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

#### Obtaining a Token
```bash
curl -X POST https://api.translator.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your-username",
    "password": "your-password"
  }'
```

#### Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400,
  "token_type": "Bearer"
}
```

### API Key Authentication (Alternative)
For service-to-service communication:

```http
X-API-Key: your-api-key
```

---

## Core Endpoints

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": 1703123456,
  "version": "1.0.0",
  "uptime": "72h30m15s"
}
```

### API Information
```http
GET /api/v1
```

**Response:**
```json
{
  "name": "Universal Multi-Format Multi-Language Ebook Translation API",
  "version": "1.0.0",
  "description": "High-quality universal ebook translation service",
  "endpoints": {
    "translate": "POST /api/v1/translate",
    "translate_fb2": "POST /api/v1/translate/fb2",
    "batch_translate": "POST /api/v1/translate/batch",
    "validate": "POST /api/v1/translate/validate"
  }
}
```

### Supported Languages
```http
GET /api/v1/languages
```

**Response:**
```json
{
  "languages": [
    {
      "code": "en",
      "name": "English",
      "native_name": "English",
      "direction": "ltr"
    },
    {
      "code": "sr",
      "name": "Serbian",
      "native_name": "Српски",
      "direction": "ltr"
    },
    {
      "code": "zh",
      "name": "Chinese",
      "native_name": "中文",
      "direction": "ltr"
    }
  ]
}
```

### Available Providers
```http
GET /api/v1/providers
```

**Response:**
```json
{
  "providers": [
    {
      "name": "openai",
      "display_name": "OpenAI",
      "models": ["gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"],
      "status": "available"
    },
    {
      "name": "anthropic",
      "display_name": "Anthropic",
      "models": ["claude-3-sonnet", "claude-3-opus"],
      "status": "available"
    }
  ]
}
```

---

## Translation Endpoints

### Text Translation
```http
POST /api/v1/translate
```

**Request Body:**
```json
{
  "text": "Hello, world! How are you today?",
  "provider": "openai",
  "model": "gpt-4",
  "context": "This is a friendly greeting",
  "script": "latin"
}
```

**Response:**
```json
{
  "translated_text": "Здраво, свете! Как сте данес?",
  "provider": "openai",
  "model": "gpt-4",
  "session_id": "uuid-v4-session-id",
  "duration_seconds": 1.23,
  "stats": {
    "tokens_used": 25,
    "cost_estimate": 0.001
  }
}
```

### String Translation (Batch API)
```http
POST /api/v1/translate/string
```

**Request Body:**
```json
{
  "text": "The quick brown fox jumps over the lazy dog",
  "source_language": "en",
  "target_language": "sr",
  "provider": "anthropic",
  "model": "claude-3-sonnet"
}
```

**Response:**
```json
{
  "translated_text": "Брза риђа лисица скаче преко лењг пса",
  "source_language": "en",
  "target_language": "sr",
  "provider": "anthropic",
  "model": "claude-3-sonnet",
  "duration_seconds": 0.89,
  "session_id": "uuid-v4-session-id"
}
```

### Directory Translation
```http
POST /api/v1/translate/directory
```

**Request Body:**
```json
{
  "input_path": "/path/to/books",
  "output_path": "/path/to/translated",
  "source_language": "en",
  "target_language": "es",
  "recursive": true,
  "parallel": true,
  "max_concurrency": 5,
  "output_format": "epub",
  "provider": "openai"
}
```

**Response:**
```json
{
  "session_id": "uuid-v4-session-id",
  "status": "processing",
  "files_found": 12,
  "files_processed": 0,
  "estimated_duration": 300,
  "progress": 0.0
}
```

---

## File Operations

### FB2 File Translation
```http
POST /api/v1/translate/fb2
Content-Type: multipart/form-data
```

**Form Data:**
```
file: [FB2 file]
provider: openai
model: gpt-4
script: cyrillic
```

**Response:**
```json
{
  "session_id": "uuid-v4-session-id",
  "download_url": "/api/v1/download/uuid-v4-session-id.epub",
  "status": "completed",
  "stats": {
    "chapters_translated": 25,
    "words_translated": 15420,
    "duration_seconds": 45.6
  }
}
```

### Ebook Translation
```http
POST /api/v1/translate/ebook
```

**Request Body:**
```json
{
  "input_path": "/path/to/book.epub",
  "output_path": "/path/to/translated.epub",
  "target_language": "fr",
  "provider": "deepseek",
  "model": "deepseek-chat",
  "preserve_formatting": true,
  "quality_preset": "high"
}
```

**Response:**
```json
{
  "session_id": "uuid-v4-session-id",
  "status": "processing",
  "progress": {
    "current_chapter": 5,
    "total_chapters": 20,
    "percentage": 25.0
  }
}
```

---

## Batch Operations

### Batch Text Translation
```http
POST /api/v1/translate/batch
```

**Request Body:**
```json
{
  "texts": [
    "Hello world",
    "How are you?",
    "Good morning"
  ],
  "target_language": "es",
  "provider": "openai",
  "parallel": true
}
```

**Response:**
```json
{
  "session_id": "uuid-v4-session-id",
  "results": [
    {
      "index": 0,
      "original": "Hello world",
      "translated": "Hola mundo",
      "status": "completed"
    },
    {
      "index": 1,
      "original": "How are you?",
      "translated": "¿Cómo estás?",
      "status": "completed"
    },
    {
      "index": 2,
      "original": "Good morning",
      "translated": "Buenos días",
      "status": "completed"
    }
  ],
  "stats": {
    "total_texts": 3,
    "successful": 3,
    "failed": 0,
    "duration_seconds": 2.1
  }
}
```

---

## Status & Monitoring

### Translation Status
```http
GET /api/v1/status/{session_id}
```

**Response:**
```json
{
  "session_id": "uuid-v4-session-id",
  "status": "processing",
  "progress": 65.5,
  "started_at": "2024-01-15T10:30:00Z",
  "estimated_completion": "2024-01-15T10:35:30Z",
  "current_operation": "Translating chapter 15 of 20",
  "stats": {
    "chapters_completed": 15,
    "total_chapters": 20,
    "words_translated": 8750,
    "tokens_used": 12500
  }
}
```

### System Statistics
```http
GET /api/v1/stats
```

**Response:**
```json
{
  "system": {
    "uptime": "72h30m15s",
    "version": "1.0.0",
    "active_sessions": 5,
    "total_translations": 1250
  },
  "performance": {
    "avg_translation_time": 2.3,
    "success_rate": 98.5,
    "requests_per_minute": 45.2
  },
  "resources": {
    "cpu_usage": 25.5,
    "memory_usage": 68.2,
    "disk_usage": 45.8
  }
}
```

### Cancel Translation
```http
POST /api/v1/translate/cancel/{session_id}
```

**Response:**
```json
{
  "session_id": "uuid-v4-session-id",
  "status": "cancelled",
  "message": "Translation cancelled successfully",
  "stats": {
    "chapters_completed": 8,
    "total_chapters": 20,
    "completion_percentage": 40.0
  }
}
```

---

## Error Handling

### Standard Error Response Format
```json
{
  "error": "Human readable error message",
  "error_code": "VALIDATION_ERROR",
  "details": {
    "field": "text",
    "issue": "required field missing"
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "uuid-v4-request-id"
}
```

### Common Error Codes

| Error Code | HTTP Status | Description |
|------------|--------------|-------------|
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `AUTHENTICATION_ERROR` | 401 | Invalid or missing authentication |
| `AUTHORIZATION_ERROR` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server internal error |
| `SERVICE_UNAVAILABLE` | 503 | Translation service unavailable |
| `FILE_TOO_LARGE` | 413 | File exceeds size limit |
| `UNSUPPORTED_FORMAT` | 415 | File format not supported |

### Validation Errors
```json
{
  "error": "Validation failed",
  "error_code": "VALIDATION_ERROR",
  "validation_errors": [
    {
      "field": "text",
      "message": "Text is required and cannot be empty"
    },
    {
      "field": "target_language",
      "message": "Invalid language code 'xyz'"
    }
  ]
}
```

---

## Rate Limiting

### Rate Limit Headers
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1703123456
X-RateLimit-Retry-After: 60
```

### Rate Limits by Plan

| Plan | Requests/Minute | Burst | Concurrent |
|------|-----------------|--------|------------|
| Free | 10 | 20 | 2 |
| Basic | 50 | 100 | 5 |
| Pro | 200 | 400 | 20 |
| Enterprise | 1000 | 2000 | 100 |

### Rate Limit Response
```json
{
  "error": "Rate limit exceeded",
  "error_code": "RATE_LIMIT_EXCEEDED",
  "retry_after": 60,
  "limit": 100,
  "window": "60s"
}
```

---

## WebSocket API

### Connection
```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
```

### Authentication
```javascript
ws.onopen = function() {
  ws.send(JSON.stringify({
    type: 'auth',
    token: 'your-jwt-token'
  }));
};
```

### Subscribe to Translation Updates
```javascript
ws.send(JSON.stringify({
  type: 'subscribe',
  session_id: 'uuid-v4-session-id',
  events: ['progress', 'completed', 'error']
}));
```

### Real-time Progress Updates
```javascript
ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  
  switch(data.type) {
    case 'translation_progress':
      console.log(`Progress: ${data.progress}%`);
      console.log(`Current: ${data.current_operation}`);
      break;
      
    case 'translation_completed':
      console.log('Translation completed!');
      console.log(`Download: ${data.download_url}`);
      break;
      
    case 'translation_error':
      console.error('Translation failed:', data.error);
      break;
  }
};
```

### WebSocket Message Types

| Type | Description | Fields |
|------|-------------|---------|
| `auth` | Authentication | `token` |
| `subscribe` | Subscribe to events | `session_id`, `events` |
| `translation_progress` | Progress update | `session_id`, `progress`, `current_operation` |
| `translation_completed` | Completion | `session_id`, `download_url`, `stats` |
| `translation_error` | Error | `session_id`, `error`, `error_code` |

---

## Examples & Use Cases

### Example 1: Simple Text Translation
```bash
curl -X POST http://localhost:8080/api/v1/translate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "text": "Hello, world!",
    "target_language": "es",
    "provider": "openai"
  }'
```

### Example 2: Translate EPUB File
```bash
curl -X POST http://localhost:8080/api/v1/translate/ebook \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "input_path": "/path/to/book.epub",
    "output_path": "/path/to/translated.epub",
    "target_language": "fr",
    "provider": "anthropic"
  }'
```

### Example 3: Batch Translation with Progress
```javascript
// Start translation
const response = await fetch('/api/v1/translate/batch', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer ' + token
  },
  body: JSON.stringify({
    texts: ['Hello', 'World', 'How are you?'],
    target_language: 'es'
  })
});

const { session_id } = await response.json();

// Monitor progress via WebSocket
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onopen = () => {
  ws.send(JSON.stringify({
    type: 'subscribe',
    session_id: session_id
  }));
};
```

### Example 4: Directory Translation
```python
import requests

# Start directory translation
response = requests.post('http://localhost:8080/api/v1/translate/directory', 
    headers={
        'Authorization': 'Bearer your-token',
        'Content-Type': 'application/json'
    },
    json={
        'input_path': '/books/english',
        'output_path': '/books/spanish',
        'target_language': 'es',
        'recursive': True,
        'parallel': True
    })

session_data = response.json()
session_id = session_data['session_id']

# Poll for completion
while True:
    status_response = requests.get(
        f'http://localhost:8080/api/v1/status/{session_id}',
        headers={'Authorization': 'Bearer your-token'}
    )
    
    status = status_response.json()
    print(f"Progress: {status['progress']}%")
    
    if status['status'] in ['completed', 'failed']:
        break
    
    time.sleep(5)
```

---

## SDK & Integration

### Official SDKs

#### Node.js
```bash
npm install @translator/universal-sdk
```

```javascript
import { Translator } from '@translator/universal-sdk';

const translator = new Translator({
  apiKey: 'your-api-key',
  baseURL: 'https://api.translator.example.com'
});

const result = await translator.translate({
  text: 'Hello, world!',
  targetLanguage: 'es',
  provider: 'openai'
});
```

#### Python
```bash
pip install universal-translator-sdk
```

```python
from universal_translator import Translator

translator = Translator(
    api_key='your-api-key',
    base_url='https://api.translator.example.com'
)

result = translator.translate(
    text='Hello, world!',
    target_language='es',
    provider='openai'
)
```

#### Go
```bash
go get github.com/translator/universal-go-sdk
```

```go
import "github.com/translator/universal-go-sdk"

client := sdk.NewClient("your-api-key", "https://api.translator.example.com")

result, err := client.Translate(sdk.TranslateRequest{
    Text: "Hello, world!",
    TargetLanguage: "es",
    Provider: "openai",
})
```

### Integration Patterns

#### Microservice Integration
```yaml
# docker-compose.yml
version: '3.8'
services:
  app:
    image: your-app:latest
    environment:
      - TRANSLATOR_API_URL=http://translator:8080
      - TRANSLATOR_API_KEY=${API_KEY}
    depends_on:
      - translator
      
  translator:
    image: universal-translator:latest
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis://redis:6379
      - DATABASE_URL=postgresql://user:pass@postgres:5432/translator
    depends_on:
      - redis
      - postgres
```

#### Webhook Integration
```json
{
  "webhook_url": "https://your-app.com/webhooks/translation",
  "events": ["completed", "failed"],
  "secret": "webhook-secret"
}
```

---

## Best Practices

### Performance Optimization
1. **Use Batch Operations**: Translate multiple texts at once
2. **Enable Caching**: Cache frequently translated content
3. **Choose Appropriate Models**: Balance quality vs speed
4. **Monitor Progress**: Use WebSocket for real-time updates
5. **Handle Timeouts**: Set appropriate timeout values

### Security Considerations
1. **Secure API Keys**: Never expose keys in client-side code
2. **Use HTTPS**: Always use TLS in production
3. **Validate Inputs**: Sanitize all user inputs
4. **Rate Limiting**: Implement client-side rate limiting
5. **Audit Logs**: Monitor API usage patterns

### Error Handling
1. **Retry Logic**: Implement exponential backoff
2. **Graceful Degradation**: Fallback to alternative providers
3. **User Feedback**: Provide clear error messages
4. **Logging**: Log errors for debugging
5. **Monitoring**: Set up alerts for failures

---

## Support & Resources

### Documentation
- [API Reference](./API.md)
- [Configuration Guide](./CONFIGURATION.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [Troubleshooting](./TROUBLESHOOTING.md)

### Community
- [GitHub Issues](https://github.com/translator/universal/issues)
- [Discord Community](https://discord.gg/translator)
- [Stack Overflow](https://stackoverflow.com/questions/tagged/universal-translator)

### Support
- Email: support@translator.example.com
- Documentation: https://docs.translator.example.com
- Status Page: https://status.translator.example.com

---

*Last Updated: January 2024*
*Version: 1.0.0*