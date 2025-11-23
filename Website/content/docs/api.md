---
title: "API Documentation"
date: "2024-01-15"
weight: 30
---

# API Documentation

Comprehensive API documentation for integrating the Universal Multi-Format Multi-Language Ebook Translation System into your applications.

## Overview

The Translation API provides RESTful endpoints for translation services, file processing, quality verification, and distributed system management.

### Base URL

```
Production: https://api.translator.digital
Development: http://localhost:8080
```

### Authentication

All API requests require authentication using JWT tokens:

```bash
curl -X POST https://api.translator.digital/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "your-username", "password": "your-password"}'
```

Include the token in subsequent requests:

```bash
curl -X GET https://api.translator.digital/api/v1/providers \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Authentication Endpoints

### POST /api/v1/auth/login

Authenticate user and receive JWT token.

**Request:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "id": 1,
    "username": "admin",
    "role": "user"
  }
}
```

### POST /api/v1/auth/refresh

Refresh authentication token.

**Headers:** `Authorization: Bearer <token>`

**Response:**
```json
{
  "token": "new-jwt-token",
  "expires_in": 3600
}
```

## Translation Endpoints

### POST /api/v1/translate/translate

Translate text between languages.

**Request:**
```json
{
  "text": "Hello, world!",
  "source_lang": "en",
  "target_lang": "sr",
  "provider": "openai",
  "options": {
    "model": "gpt-4",
    "temperature": 0.7,
    "max_tokens": 1000
  }
}
```

**Response:**
```json
{
  "translated_text": "Здраво, свете!",
  "source_lang": "en",
  "target_lang": "sr",
  "provider": "openai",
  "model": "gpt-4",
  "confidence": 0.95,
  "processing_time_ms": 1200,
  "cost": 0.002,
  "tokens_used": 25
}
```

### POST /api/v1/translate/batch

Translate multiple texts in batch.

**Request:**
```json
{
  "texts": [
    "Hello, world!",
    "How are you?",
    "Good morning"
  ],
  "source_lang": "en",
  "target_lang": "sr",
  "provider": "openai",
  "options": {
    "model": "gpt-4"
  }
}
```

**Response:**
```json
{
  "results": [
    {
      "translated_text": "Здраво, свете!",
      "confidence": 0.95,
      "tokens_used": 25
    },
    {
      "translated_text": "Како си?",
      "confidence": 0.92,
      "tokens_used": 20
    },
    {
      "translated_text": "Добро јутро",
      "confidence": 0.94,
      "tokens_used": 18
    }
  ],
  "total_processed": 3,
  "processing_time_ms": 3500,
  "total_tokens_used": 63,
  "total_cost": 0.006
}
```

## Ebook Processing Endpoints

### POST /api/v1/ebook/upload

Upload and process ebook files.

**Request:** `multipart/form-data`
- `file`: The ebook file (FB2, EPUB, TXT, HTML, PDF, DOCX)
- `source_lang`: Source language code
- `target_lang`: Target language code
- `provider`: Translation provider
- `options`: Additional options (JSON string)
- `script`: Serbian script (cyrillic/latin)

**Response:**
```json
{
  "upload_id": "uuid-string",
  "filename": "example.fb2",
  "size": 1024000,
  "format": "fb2",
  "status": "uploaded",
  "source_lang": "ru",
  "target_lang": "sr",
  "provider": "deepseek",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### GET /api/v1/ebook/{upload_id}/status

Get processing status.

**Response:**
```json
{
  "upload_id": "uuid-string",
  "status": "processing",
  "progress": 75,
  "current_stage": "translating",
  "estimated_completion": "2024-01-15T10:35:00Z",
  "pages_processed": 150,
  "total_pages": 200,
  "current_page": 150,
  "processing_time_elapsed": 180000,
  "quality_score": 0.92
}
```

### GET /api/v1/ebook/{upload_id}/download

Download processed ebook.

**Query Parameters:**
- `format`: Output format (fb2, epub, pdf, docx)
- `script`: Serbian script (cyrillic, latin)

**Response:** Binary file download with appropriate headers.

## Provider Endpoints

### GET /api/v1/providers

List available translation providers.

**Response:**
```json
{
  "providers": [
    {
      "name": "openai",
      "display_name": "OpenAI GPT",
      "models": ["gpt-3.5-turbo", "gpt-4", "gpt-4-turbo"],
      "features": ["translation", "context-aware", "high-quality"],
      "status": "available",
      "pricing": {
        "input_tokens": 0.00001,
        "output_tokens": 0.00003
      }
    },
    {
      "name": "deepseek",
      "display_name": "DeepSeek",
      "models": ["deepseek-chat"],
      "features": ["translation", "cost-effective", "consistent"],
      "status": "available",
      "pricing": {
        "input_tokens": 0.000001,
        "output_tokens": 0.000002
      }
    }
  ]
}
```

### GET /api/v1/providers/{provider}/models

Get available models for specific provider.

**Response:**
```json
{
  "provider": "openai",
  "models": [
    {
      "name": "gpt-4",
      "display_name": "GPT-4",
      "max_tokens": 8192,
      "cost_per_1k_tokens": 0.03,
      "supported_languages": ["en", "ru", "sr", "zh", "es", "fr"],
      "capabilities": ["translation", "context-aware", "literary"]
    },
    {
      "name": "gpt-3.5-turbo",
      "display_name": "GPT-3.5 Turbo",
      "max_tokens": 4096,
      "cost_per_1k_tokens": 0.002,
      "supported_languages": ["en", "ru", "sr", "zh", "es", "fr"],
      "capabilities": ["translation", "fast", "cost-effective"]
    }
  ]
}
```

## Language Endpoints

### GET /api/v1/languages

Get supported languages.

**Response:**
```json
{
  "languages": [
    {
      "code": "en",
      "name": "English",
      "native_name": "English",
      "direction": "ltr",
      "quality_scores": {
        "sr": 0.94,
        "ru": 0.95,
        "fr": 0.96
      }
    },
    {
      "code": "ru",
      "name": "Russian",
      "native_name": "Русский",
      "direction": "ltr",
      "quality_scores": {
        "sr": 0.95,
        "en": 0.93,
        "fr": 0.85
      }
    },
    {
      "code": "sr",
      "name": "Serbian",
      "native_name": "Српски",
      "direction": "ltr",
      "scripts": ["cyrillic", "latin"],
      "quality_scores": {
        "ru": 0.95,
        "en": 0.92,
        "fr": 0.88
      }
    }
  ]
}
```

### GET /api/v1/languages/{source}/targets

Get available target languages for source language.

**Response:**
```json
{
  "source_language": {
    "code": "ru",
    "name": "Russian",
    "native_name": "Русский"
  },
  "target_languages": [
    {
      "code": "sr",
      "name": "Serbian",
      "native_name": "Српски",
      "quality_score": 0.95,
      "best_providers": ["zhipu", "deepseek"]
    },
    {
      "code": "en",
      "name": "English",
      "native_name": "English",
      "quality_score": 0.93,
      "best_providers": ["openai", "anthropic"]
    }
  ]
}
```

## Quality Verification Endpoints

### POST /api/v1/verification/check

Check translation quality.

**Request:**
```json
{
  "original_text": "Hello, world!",
  "translated_text": "Здраво, свете!",
  "source_lang": "en",
  "target_lang": "sr",
  "provider": "openai",
  "options": {
    "detailed_analysis": true,
    "check_grammar": true,
    "check_style": true
  }
}
```

**Response:**
```json
{
  "quality_score": 0.92,
  "confidence": 0.95,
  "verification_passed": true,
  "analysis": {
    "grammar_score": 0.95,
    "style_score": 0.88,
    "terminology_score": 0.94,
    "cultural_adaptation": 0.90
  },
  "issues": [],
  "suggestions": [
    {
      "type": "style",
      "message": "Translation style is appropriate for the context",
      "severity": "info"
    }
  ]
}
```

## Distributed Processing Endpoints

### GET /api/v1/distributed/status

Get distributed system status.

**Response:**
```json
{
  "cluster_status": "healthy",
  "total_nodes": 3,
  "active_nodes": 3,
  "total_translations": 1500,
  "processing_rate": 45.2,
  "average_response_time": 1200,
  "nodes": [
    {
      "id": "node-1",
      "status": "active",
      "load": 0.75,
      "translations_processed": 500,
      "last_seen": "2024-01-15T10:30:00Z",
      "capabilities": ["openai", "anthropic", "deepseek"]
    },
    {
      "id": "node-2",
      "status": "active",
      "load": 0.45,
      "translations_processed": 450,
      "last_seen": "2024-01-15T10:30:05Z",
      "capabilities": ["openai", "zhipu"]
    },
    {
      "id": "node-3",
      "status": "active",
      "load": 0.60,
      "translations_processed": 550,
      "last_seen": "2024-01-15T10:30:02Z",
      "capabilities": ["ollama", "llamacpp"]
    }
  ]
}
```

### POST /api/v1/distributed/translate

Submit translation request to distributed system.

**Request:**
```json
{
  "text": "Large text for distributed processing...",
  "source_lang": "en",
  "target_lang": "sr",
  "priority": "normal",
  "callback_url": "https://your-app.com/callback",
  "options": {
    "provider": "openai",
    "quality_threshold": 0.8
  }
}
```

**Response:**
```json
{
  "request_id": "uuid-string",
  "status": "queued",
  "estimated_completion": "2024-01-15T10:32:00Z",
  "queue_position": 5,
  "assigned_node": "node-1"
}
```

## Monitoring Endpoints

### GET /api/v1/metrics/health

System health check.

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": 86400,
  "services": {
    "database": "healthy",
    "cache": "healthy",
    "translation_providers": "healthy",
    "distributed_system": "healthy",
    "file_processing": "healthy"
  },
  "system_resources": {
    "cpu_usage": 0.35,
    "memory_usage": 0.60,
    "disk_usage": 0.45,
    "network_latency": 25
  }
}
```

### GET /api/v1/metrics/stats

Get system statistics.

**Response:**
```json
{
  "total_translations": 10000,
  "translations_today": 500,
  "active_users": 25,
  "success_rate": 0.98,
  "average_processing_time_ms": 1200,
  "quality_score_average": 0.92,
  "most_used_providers": [
    {"name": "openai", "usage_count": 5000, "percentage": 50},
    {"name": "deepseek", "usage_count": 3000, "percentage": 30},
    {"name": "anthropic", "usage_count": 2000, "percentage": 20}
  ],
  "most_used_language_pairs": [
    {"pair": "ru-sr", "usage_count": 4000, "percentage": 40},
    {"pair": "en-sr", "usage_count": 3000, "percentage": 30},
    {"pair": "en-ru", "usage_count": 2000, "percentage": 20}
  ]
}
```

## WebSocket API

### Real-time Translation Updates

Connect to WebSocket for real-time translation updates:

```javascript
const ws = new WebSocket('ws://localhost:8080/ws/translate');

ws.onopen = function(event) {
    // Authenticate
    ws.send(JSON.stringify({
        type: 'auth',
        token: 'your-jwt-token'
    }));
    
    // Subscribe to translation updates
    ws.send(JSON.stringify({
        type: 'subscribe',
        channel: 'translations'
    }));
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    if (data.type === 'translation_update') {
        console.log('Translation update:', data.data);
    }
};
```

**WebSocket Events:**
- `translation_started`: Translation processing started
- `translation_progress`: Translation progress update
- `translation_completed`: Translation completed
- `translation_failed`: Translation failed
- `quality_update`: Quality assessment update

## Error Handling

### Error Response Format

All error responses follow this consistent format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input parameters",
    "details": [
      {
        "field": "source_lang",
        "message": "Invalid language code"
      }
    ],
    "request_id": "uuid-string",
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid request parameters |
| 401 | Unauthorized | Authentication required or failed |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource conflict |
| 422 | Unprocessable Entity | Invalid data format |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 502 | Bad Gateway | Provider error |
| 503 | Service Unavailable | Service temporarily unavailable |

## Rate Limiting

API requests are rate-limited based on your subscription plan:

- **Free**: 100 requests/hour
- **Basic**: 1,000 requests/hour
- **Pro**: 10,000 requests/hour
- **Enterprise**: Unlimited

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642245600
X-RateLimit-Retry-After: 60
```

## SDKs and Client Libraries

### Go SDK

```go
import "github.com/digital-vasic/translator-go"

client := translator.New("your-api-key")
result, err := client.Translate.Translate(
    translator.TranslateRequest{
        Text:       "Hello, world!",
        SourceLang: "en",
        TargetLang: "sr",
        Provider:   "openai",
    },
)
```

### Python SDK

```python
from translator import Translator

client = Translator(api_key="your-api-key")
result = client.translate.translate(
    text="Hello, world!",
    source_lang="en",
    target_lang="sr",
    provider="openai"
)
```

### JavaScript/TypeScript SDK

```typescript
import { Translator } from '@translator/client';

const client = new Translator({ apiKey: 'your-api-key' });
const result = await client.translate.translate({
    text: 'Hello, world!',
    sourceLang: 'en',
    targetLang: 'sr',
    provider: 'openai'
});
```

## Webhooks

### Configure Webhooks

Receive real-time notifications for translation events:

```bash
curl -X POST https://api.translator.digital/api/v1/webhooks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-app.com/webhook",
    "events": ["translation.completed", "translation.failed"],
    "secret": "your-webhook-secret",
    "active": true
  }'
```

### Webhook Payload

```json
{
  "event": "translation.completed",
  "data": {
    "request_id": "uuid-string",
    "translation_id": "uuid-string",
    "status": "completed",
    "result": {
      "translated_text": "Здраво, свете!",
      "provider": "openai",
      "model": "gpt-4",
      "confidence": 0.95,
      "quality_score": 0.92
    },
    "processing_time_ms": 1200,
    "cost": 0.002
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## OpenAPI Specification

The complete OpenAPI 3.0 specification is available at:
```
https://api.translator.digital/openapi.json
```

### Interactive Documentation

Interactive API documentation is available at:
```
https://api.translator.digital/docs
```

This includes:
- Interactive API testing
- Request/response examples
- Parameter documentation
- Schema definitions

## Testing the API

### Using curl

```bash
# Authentication
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "test"}' \
  | jq -r '.token')

# Simple translation
curl -X POST http://localhost:8080/api/v1/translate/translate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, world!",
    "source_lang": "en",
    "target_lang": "sr",
    "provider": "openai"
  }'
```

### Using Postman

1. Import the OpenAPI specification from: `https://api.translator.digital/openapi.json`
2. Set your API key as an environment variable
3. Start making requests

### API Playground

Use our interactive API playground at:
```
https://api.translator.digital/playground
```

This web interface allows you to:
- Test all API endpoints
- Generate code snippets
- Explore documentation
- Monitor requests and responses

## Support and Resources

- **API Documentation**: This guide
- **OpenAPI Spec**: [Download specification](/openapi.json)
- **SDK Documentation**: Language-specific guides
- **Status Page**: [api.translator.digital/status](https://api.translator.digital/status)
- **Support**: [support@translator.digital](mailto:support@translator.digital)
- **Community**: [GitHub Discussions](https://github.com/digital-vasic/translator/discussions)

---

For more information about integrating the API into your applications, check out our [API Integration Tutorial](/tutorials/api-integration).