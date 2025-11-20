# API Documentation

Complete API documentation for the Universal Ebook Translator service.

**Version:** 2.2.0
**Last Updated:** 2025-11-20

## Overview

The Universal Ebook Translator provides REST API and WebSocket endpoints for translating ebooks between any language pair, with support for multiple formats (EPUB, FB2, TXT, HTML) and translation providers.

## Base URL

- **HTTP/3 (QUIC)**: `https://localhost:8443`
- **HTTP/2 (TLS)**: `https://localhost:8443`
- **HTTP/1.1**: `http://localhost:8080` (if HTTP/3 disabled)

## Authentication

The API supports optional authentication via JWT tokens or API keys.

### JWT Authentication

Include the JWT token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### API Key Authentication

Include the API key in the header:

```http
X-API-Key: <your-api-key>
```

## Endpoints

### Health & Status

#### `GET /health`

Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "time": "2025-01-15T10:30:00Z"
}
```

#### `GET /`

API information and available endpoints.

#### `GET /api/v1/providers`

List all available translation providers.

**Response:**
```json
{
  "providers": [
    {
      "name": "dictionary",
      "description": "Simple dictionary-based translation",
      "requires_api_key": false
    },
    {
      "name": "openai",
      "description": "OpenAI GPT models",
      "requires_api_key": true,
      "models": ["gpt-4", "gpt-3.5-turbo"]
    }
  ]
}
```

### Translation

#### `POST /api/v1/translate`

Translate Russian text to Serbian.

**Request:**
```json
{
  "text": "Привет, мир!",
  "provider": "dictionary",
  "model": "gpt-4",
  "context": "Literary text",
  "script": "cyrillic"
}
```

**Response:**
```json
{
  "original": "Привет, мир!",
  "translated": "Здраво, свете!",
  "provider": "dictionary",
  "session_id": "uuid",
  "stats": {
    "total": 1,
    "translated": 1,
    "cached": 0,
    "errors": 0
  }
}
```

#### `POST /api/v1/translate/fb2`

Translate a complete FB2 e-book file.

**Request:**
```http
POST /api/v1/translate/fb2
Content-Type: multipart/form-data

file: <fb2 file>
provider: dictionary
model: gpt-4
script: cyrillic
```

**Response:**
Returns the translated FB2 file as `application/xml`.

#### `POST /api/v1/translate/batch`

Batch translate multiple texts.

**Request:**
```json
{
  "texts": [
    "герой",
    "мир",
    "человек"
  ],
  "provider": "dictionary",
  "context": "Single words"
}
```

**Response:**
```json
{
  "originals": ["герой", "мир", "человек"],
  "translated": ["јунак", "свет", "човек"],
  "provider": "dictionary",
  "session_id": "uuid",
  "stats": {...}
}
```

#### `POST /api/v1/translate/string` (v2.2)

Translate a text string directly without file input. Perfect for quick translations, API integrations, or pipeline workflows.

**Request:**
```json
{
  "text": "Hello, world!",
  "source_language": "en",
  "target_language": "sr",
  "provider": "deepseek",
  "model": "deepseek-chat"
}
```

**Parameters:**
- `text` (required): The text to translate
- `source_language` (optional): Source language code (auto-detected if not provided)
- `target_language` (required): Target language code (e.g., "sr", "en", "de")
- `provider` (optional): Translation provider (default: "dictionary")
- `model` (optional): LLM model to use (provider-specific)

**Response:**
```json
{
  "translated_text": "Здраво, свете!",
  "source_language": "en",
  "target_language": "sr",
  "provider": "deepseek",
  "duration_seconds": 1.23,
  "session_id": "sess_abc123"
}
```

**curl Example:**
```bash
curl -X POST https://localhost:8443/api/v1/translate/string \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, world!",
    "target_language": "sr",
    "provider": "deepseek"
  }'
```

**Error Responses:**
- `400 Bad Request`: Missing required fields or invalid language code
- `500 Internal Server Error`: Translation failed

#### `POST /api/v1/translate/directory` (v2.2)

Translate entire directories with structure preservation. Automatically discovers all supported ebook files and maintains directory hierarchy in output.

**Request:**
```json
{
  "input_path": "Books/",
  "output_path": "Translated/",
  "source_language": "en",
  "target_language": "sr",
  "recursive": true,
  "parallel": true,
  "max_concurrency": 4,
  "provider": "deepseek",
  "model": "deepseek-chat",
  "output_format": "epub"
}
```

**Parameters:**
- `input_path` (required): Directory path to translate
- `output_path` (optional): Output directory (defaults to input_path with suffix)
- `source_language` (optional): Source language code (auto-detected if not provided)
- `target_language` (required): Target language code
- `recursive` (optional): Process subdirectories (default: false)
- `parallel` (optional): Enable parallel processing (default: false)
- `max_concurrency` (optional): Number of parallel workers (default: 4)
- `provider` (optional): Translation provider (default: "dictionary")
- `model` (optional): LLM model to use
- `output_format` (optional): Output format (epub, fb2, txt, html)

**Response:**
```json
{
  "session_id": "sess_xyz789",
  "total_files": 15,
  "successful": 14,
  "failed": 1,
  "duration_seconds": 324.56,
  "results": [
    {
      "input_path": "Books/book1.epub",
      "output_path": "Translated/book1_sr.epub",
      "success": true
    },
    {
      "input_path": "Books/book2.epub",
      "output_path": "Translated/book2_sr.epub",
      "success": false,
      "error": "API rate limit exceeded"
    }
  ]
}
```

**curl Example:**
```bash
curl -X POST https://localhost:8443/api/v1/translate/directory \
  -H "Content-Type: application/json" \
  -d '{
    "input_path": "Books/",
    "output_path": "Translated/",
    "target_language": "sr",
    "recursive": true,
    "parallel": true,
    "max_concurrency": 4,
    "provider": "deepseek"
  }'
```

**Directory Structure Example:**
```
Input:
Books/
├── Fiction/
│   ├── book1.epub
│   └── book2.fb2
└── NonFiction/
    └── book3.txt

Output (with recursive=true):
Translated/
├── Fiction/
│   ├── book1_sr.epub
│   └── book2_sr.epub
└── NonFiction/
    └── book3_sr.epub
```

**Error Responses:**
- `400 Bad Request`: Invalid parameters or directory doesn't exist
- `500 Internal Server Error`: Translation failed for all files

### Script Conversion

#### `POST /api/v1/convert/script`

Convert between Cyrillic and Latin scripts.

**Request:**
```json
{
  "text": "Ратибор је јунак",
  "target": "latin"
}
```

**Response:**
```json
{
  "original": "Ратибор је јунак",
  "converted": "Ratibor je junak",
  "target": "latin"
}
```

### WebSocket

#### `GET /ws?session_id={id}`

WebSocket endpoint for real-time translation progress with enhanced v2.1 progress tracking.

**Connection:**
```javascript
const ws = new WebSocket('wss://localhost:8443/ws?session_id=uuid');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log(data.type, data.message, data.data);
};
```

**Event Types:**
- `translation_started` - Translation begins
- `translation_progress` - Progress updates (includes detailed metrics)
- `translation_completed` - Translation finished
- `translation_error` - Error occurred

**Progress Event Data (v2.1 Enhanced):**
```json
{
  "type": "translation_progress",
  "session_id": "sess_abc123",
  "timestamp": "2025-11-20T10:30:00Z",
  "message": "Translating chapter 15 of 38",
  "data": {
    "book_title": "Son Nad Bezdnoy",
    "total_chapters": 38,
    "current_chapter": 15,
    "chapter_title": "Chapter 15: The Revelation",
    "percent_complete": 39.5,
    "items_completed": 30,
    "items_failed": 2,
    "items_total": 76,
    "elapsed_time": "5 minutes 23 seconds",
    "estimated_eta": "8 minutes 15 seconds",
    "source_language": "ru",
    "target_language": "sr",
    "provider": "deepseek",
    "status": "translating"
  }
}
```

**Full JavaScript Example:**
```javascript
const sessionId = 'your-session-id';
const ws = new WebSocket(`wss://localhost:8443/ws?session_id=${sessionId}`);

ws.onopen = () => {
  console.log('WebSocket connected');
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);

  switch (msg.type) {
    case 'translation_started':
      console.log('Translation started:', msg.data.book_title);
      break;

    case 'translation_progress':
      const progress = msg.data;
      console.log(`Progress: ${progress.percent_complete.toFixed(1)}%`);
      console.log(`Chapter: ${progress.current_chapter}/${progress.total_chapters}`);
      console.log(`ETA: ${progress.estimated_eta}`);
      console.log(`Elapsed: ${progress.elapsed_time}`);
      break;

    case 'translation_completed':
      console.log('Translation completed!');
      console.log(`Duration: ${msg.data.duration_seconds}s`);
      break;

    case 'translation_error':
      console.error('Translation error:', msg.message);
      break;
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('WebSocket disconnected');
};
```

## Supported Languages (v2.0 Universal Support)

The translator supports translation between any language pair. Common language codes:

| Language | Code | Example |
|----------|------|---------|
| English | `en` | English |
| Serbian | `sr` | Српски |
| Russian | `ru` | Русский |
| German | `de` | Deutsch |
| French | `fr` | Français |
| Spanish | `es` | Español |
| Italian | `it` | Italiano |
| Portuguese | `pt` | Português |
| Chinese | `zh` | 中文 |
| Japanese | `ja` | 日本語 |
| Korean | `ko` | 한국어 |
| Arabic | `ar` | العربية |
| Polish | `pl` | Polski |
| Czech | `cs` | Čeština |
| Croatian | `hr` | Hrvatski |
| Bosnian | `bs` | Bosanski |
| Bulgarian | `bg` | Български |
| Ukrainian | `uk` | Українська |

**Automatic Language Detection:**
If `source_language` is not provided, the system will automatically detect it.

**Serbian Scripts:**
Serbian output supports both Cyrillic (default) and Latin scripts via the `/api/v1/convert/script` endpoint.

## Supported Formats

| Format | Extension | Read | Write | Notes |
|--------|-----------|------|-------|-------|
| EPUB | `.epub` | ✓ | ✓ | Most common ebook format |
| FB2 | `.fb2` | ✓ | ✓ | FictionBook 2.0 XML format |
| TXT | `.txt` | ✓ | ✓ | Plain text |
| HTML | `.html` | ✓ | ✓ | HTML documents |

All formats can be converted to any other supported format during translation.

## Provider Configuration

### Dictionary Provider

No configuration required. Uses built-in Russian-Serbian dictionary.

```json
{
  "provider": "dictionary"
}
```

### OpenAI Provider

Requires API key via environment variable or config:

```bash
export OPENAI_API_KEY="your-key"
```

```json
{
  "provider": "openai",
  "model": "gpt-4"
}
```

### Anthropic Provider

```bash
export ANTHROPIC_API_KEY="your-key"
```

```json
{
  "provider": "anthropic",
  "model": "claude-3-sonnet-20240229"
}
```

### Zhipu AI Provider

```bash
export ZHIPU_API_KEY="your-key"
```

```json
{
  "provider": "zhipu",
  "model": "glm-4"
}
```

### DeepSeek Provider

```bash
export DEEPSEEK_API_KEY="your-key"
```

```json
{
  "provider": "deepseek",
  "model": "deepseek-chat"
}
```

### Ollama Provider (Local)

Requires Ollama running locally:

```bash
ollama pull llama3:8b
```

```json
{
  "provider": "ollama",
  "model": "llama3:8b"
}
```

## Error Handling

The API returns standard HTTP status codes:

- `200 OK` - Request successful
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication required
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

**Error Response:**
```json
{
  "error": "Error message description"
}
```

## Rate Limiting

Default rate limits:
- **10 requests per second** per IP
- **20 burst requests** allowed

Rate limit headers:
```http
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 9
X-RateLimit-Reset: 1642251600
```

## Version History

### v2.2.0 (2025-11-20)
**New Endpoints:**
- `POST /api/v1/translate/string` - Direct string translation
- `POST /api/v1/translate/directory` - Batch directory translation

**Features:**
- String input support (no file required)
- Stdin/pipeline support for Unix workflows
- Recursive directory processing
- Structure-preserving output
- Parallel processing with configurable concurrency
- Enhanced REST API with curl examples

### v2.1.0 (2025-11-18)
**Enhanced Features:**
- WebSocket progress events with detailed metrics:
  - Percentage complete
  - Current chapter/part information
  - Estimated time to completion (ETA)
  - Elapsed time tracking
  - Items completed/failed/total
- Storage backends:
  - PostgreSQL with ACID compliance
  - SQLite with SQLCipher encryption
  - Redis for high-performance caching
- Docker infrastructure with docker-compose
- Management scripts for deployment

### v2.0.0 (2025-11-15)
**Major Release:**
- Universal language support (18+ languages, any pair)
- Universal format support (EPUB, FB2, TXT, HTML)
- Automatic language detection
- Multiple translation providers (OpenAI, Anthropic, Zhipu, DeepSeek, Ollama, Dictionary)
- Format conversion during translation
- Enhanced metadata preservation
- Comprehensive test coverage

## Examples

See the `/api/examples` directory for:
- **curl** scripts
- **HTTP** files
- **Postman** collection
- **WebSocket** test page

## Quick Start

**Translate a single string:**
```bash
curl -X POST https://localhost:8443/api/v1/translate/string \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello", "target_language": "sr"}'
```

**Translate a directory:**
```bash
curl -X POST https://localhost:8443/api/v1/translate/directory \
  -H "Content-Type: application/json" \
  -d '{
    "input_path": "Books/",
    "target_language": "sr",
    "recursive": true,
    "parallel": true
  }'
```

**Monitor progress via WebSocket:**
```javascript
const ws = new WebSocket('wss://localhost:8443/ws?session_id=your-session');
ws.onmessage = (e) => console.log(JSON.parse(e.data));
```
