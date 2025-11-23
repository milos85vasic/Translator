# Configuration Reference - Universal Multi-Format Multi-Language Ebook Translation System

## Table of Contents

1. [Overview](#overview)
2. [Configuration File Structure](#configuration-file-structure)
3. [Server Configuration](#server-configuration)
4. [Database Configuration](#database-configuration)
5. [Cache Configuration](#cache-configuration)
6. [LLM Provider Configuration](#llm-provider-configuration)
7. [Security Configuration](#security-configuration)
8. [Logging Configuration](#logging-configuration)
9. [Worker Configuration](#worker-configuration)
10. [Storage Configuration](#storage-configuration)
11. [Queue Configuration](#queue-configuration)
12. [Environment Variables](#environment-variables)
13. [Configuration Examples](#configuration-examples)
14. [Validation and Testing](#validation-and-testing)

## Overview

The Universal Multi-Format Multi-Language Ebook Translation System uses a comprehensive configuration system that supports JSON configuration files, environment variables, and command-line arguments. This reference covers all available configuration options.

### Configuration Precedence

1. Command-line arguments (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

### Configuration File Locations

The system searches for configuration files in the following order:

1. `--config` command-line argument
2. `TRANSLATOR_CONFIG` environment variable
3. `./config.json`
4. `/etc/translator/config.json`
5. `$HOME/.translator/config.json`

## Configuration File Structure

```json
{
  "server": {},
  "database": {},
  "cache": {},
  "llm": {},
  "security": {},
  "logging": {},
  "worker": {},
  "storage": {},
  "queue": {},
  "coordination": {},
  "monitoring": {},
  "features": {}
}
```

## Server Configuration

### Basic Server Settings

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "tls_port": 8081,
    "enable_tls": true,
    "read_timeout": "30s",
    "write_timeout": "30s",
    "idle_timeout": "60s",
    "max_header_bytes": 1048576,
    "max_request_size": 104857600
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `host` | string | `"0.0.0.0"` | Server bind address |
| `port` | int | `8080` | HTTP port |
| `tls_port` | int | `8081` | HTTPS port |
| `enable_tls` | bool | `false` | Enable TLS/SSL |
| `read_timeout` | duration | `"30s"` | Maximum time to read request |
| `write_timeout` | duration | `"30s"` | Maximum time to write response |
| `idle_timeout` | duration | `"60s"` | Maximum idle time |
| `max_header_bytes` | int | `1048576` | Maximum header size in bytes |
| `max_request_size` | int | `104857600` | Maximum request size in bytes |

### WebSocket Configuration

```json
{
  "server": {
    "websocket": {
      "enabled": true,
      "port": 8082,
      "path": "/ws",
      "ping_interval": "30s",
      "pong_wait": "60s",
      "write_wait": "10s",
      "max_message_size": 1048576
    }
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `enabled` | bool | `true` | Enable WebSocket support |
| `port` | int | `8082` | WebSocket port |
| `path` | string | `"/ws"` | WebSocket endpoint path |
| `ping_interval` | duration | `"30s"` | Ping interval |
| `pong_wait` | duration | `"60s"` | Wait time for pong |
| `write_wait` | duration | `"10s"` | Write wait time |
| `max_message_size` | int | `1048576` | Maximum message size |

### Coordination Configuration

```json
{
  "server": {
    "coordination": {
      "enabled": false,
      "worker_timeout": "30s",
      "heartbeat_interval": "10s",
      "max_workers": 100,
      "load_balancing": "round_robin"
    }
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `enabled` | bool | `false` | Enable coordination mode |
| `worker_timeout` | duration | `"30s"` | Worker timeout |
| `heartbeat_interval` | duration | `"10s"` | Heartbeat interval |
| `max_workers` | int | `100` | Maximum workers |
| `load_balancing` | string | `"round_robin"` | Load balancing strategy |

## Database Configuration

### SQLite Configuration

```json
{
  "database": {
    "type": "sqlite",
    "connection_string": "/var/lib/translator/translator.db",
    "max_open_connections": 25,
    "max_idle_connections": 5,
    "connection_max_lifetime": "5m",
    "auto_migrate": true
  }
}
```

### PostgreSQL Configuration

```json
{
  "database": {
    "type": "postgresql",
    "connection_string": "postgres://user:password@localhost:5432/translator?sslmode=disable",
    "max_open_connections": 100,
    "max_idle_connections": 10,
    "connection_max_lifetime": "1h",
    "auto_migrate": true
  }
}
```

### MySQL Configuration

```json
{
  "database": {
    "type": "mysql",
    "connection_string": "user:password@tcp(localhost:3306)/translator?charset=utf8mb4&parseTime=True&loc=Local",
    "max_open_connections": 100,
    "max_idle_connections": 10,
    "connection_max_lifetime": "1h",
    "auto_migrate": true
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `type` | string | `"sqlite"` | Database type |
| `connection_string` | string | - | Database connection string |
| `max_open_connections` | int | `25` | Maximum open connections |
| `max_idle_connections` | int | `5` | Maximum idle connections |
| `connection_max_lifetime` | duration | `"5m"` | Connection lifetime |
| `auto_migrate` | bool | `true` | Auto-migrate schema |

## Cache Configuration

### Memory Cache

```json
{
  "cache": {
    "type": "memory",
    "ttl": "1h",
    "max_size": 1000,
    "cleanup_interval": "10m"
  }
}
```

### Redis Cache

```json
{
  "cache": {
    "type": "redis",
    "connection_string": "redis://localhost:6379",
    "ttl": "1h",
    "max_memory": "512mb",
    "eviction_policy": "allkeys-lru",
    "pool_size": 10
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `type` | string | `"memory"` | Cache type |
| `connection_string` | string | - | Cache connection string |
| `ttl` | duration | `"1h"` | Default TTL |
| `max_size` | int | `1000` | Maximum entries (memory) |
| `max_memory` | string | `"512mb"` | Maximum memory (Redis) |
| `eviction_policy` | string | `"allkeys-lru"` | Eviction policy |
| `cleanup_interval` | duration | `"10m"` | Cleanup interval |
| `pool_size` | int | `10` | Connection pool size |

## LLM Provider Configuration

### Anthropic Claude

```json
{
  "llm": {
    "provider": "anthropic",
    "api_key": "${ANTHROPIC_API_KEY}",
    "model": "claude-3-sonnet-20240229",
    "max_tokens": 4096,
    "temperature": 0.3,
    "top_p": 0.9,
    "timeout": "60s",
    "retry_attempts": 3,
    "retry_delay": "1s"
  }
}
```

### OpenAI

```json
{
  "llm": {
    "provider": "openai",
    "api_key": "${OPENAI_API_KEY}",
    "model": "gpt-4",
    "max_tokens": 4096,
    "temperature": 0.3,
    "top_p": 0.9,
    "timeout": "60s",
    "retry_attempts": 3,
    "retry_delay": "1s"
  }
}
```

### DeepSeek

```json
{
  "llm": {
    "provider": "deepseek",
    "api_key": "${DEEPSEEK_API_KEY}",
    "model": "deepseek-chat",
    "max_tokens": 4096,
    "temperature": 0.3,
    "top_p": 0.9,
    "timeout": "60s",
    "retry_attempts": 3,
    "retry_delay": "1s"
  }
}
```

### Ollama (Local)

```json
{
  "llm": {
    "provider": "ollama",
    "base_url": "http://localhost:11434",
    "model": "llama2",
    "max_tokens": 4096,
    "temperature": 0.3,
    "top_p": 0.9,
    "timeout": "120s",
    "retry_attempts": 3,
    "retry_delay": "2s"
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `provider` | string | `"anthropic"` | LLM provider |
| `api_key` | string | - | API key |
| `base_url` | string | - | Base URL (Ollama) |
| `model` | string | - | Model name |
| `max_tokens` | int | `4096` | Maximum tokens |
| `temperature` | float | `0.3` | Temperature (0.0-1.0) |
| `top_p` | float | `0.9` | Top P (0.0-1.0) |
| `timeout` | duration | `"60s"` | Request timeout |
| `retry_attempts` | int | `3` | Retry attempts |
| `retry_delay` | duration | `"1s"` | Retry delay |

## Security Configuration

### API Security

```json
{
  "security": {
    "api_key_required": true,
    "api_key_header": "X-API-Key",
    "api_key_query_param": "api_key",
    "jwt_secret": "${JWT_SECRET}",
    "jwt_expiry": "24h",
    "cors": {
      "enabled": true,
      "allowed_origins": ["*"],
      "allowed_methods": ["GET", "POST", "PUT", "DELETE"],
      "allowed_headers": ["*"],
      "max_age": "86400"
    }
  }
}
```

### Rate Limiting

```json
{
  "security": {
    "rate_limit": {
      "enabled": true,
      "requests_per_minute": 60,
      "burst": 10,
      "cleanup_interval": "1m"
    }
  }
}
```

### TLS Configuration

```json
{
  "security": {
    "tls": {
      "cert_file": "/etc/translator/certs/server.crt",
      "key_file": "/etc/translator/certs/server.key",
      "ca_file": "",
      "min_version": "1.2",
      "max_version": "1.3",
      "cipher_suites": [
        "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
        "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
      ]
    }
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `api_key_required` | bool | `true` | Require API key |
| `api_key_header` | string | `"X-API-Key"` | API key header |
| `api_key_query_param` | string | `"api_key"` | API key query param |
| `jwt_secret` | string | - | JWT secret |
| `jwt_expiry` | duration | `"24h"` | JWT expiry |
| `cors.enabled` | bool | `true` | Enable CORS |
| `cors.allowed_origins` | array | `["*"]` | Allowed origins |
| `cors.allowed_methods` | array | `["GET","POST"]` | Allowed methods |
| `cors.allowed_headers` | array | `["*"]` | Allowed headers |
| `cors.max_age` | int | `86400` | CORS max age |
| `rate_limit.enabled` | bool | `true` | Enable rate limiting |
| `rate_limit.requests_per_minute` | int | `60` | Requests per minute |
| `rate_limit.burst` | int | `10` | Burst size |
| `rate_limit.cleanup_interval` | duration | `"1m"` | Cleanup interval |

## Logging Configuration

### Basic Logging

```json
{
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout",
    "file": "/var/log/translator/translator.log",
    "max_size": "100MB",
    "max_backups": 10,
    "max_age": 30,
    "compress": true
  }
}
```

### Structured Logging

```json
{
  "logging": {
    "level": "info",
    "format": "json",
    "output": "file",
    "file": "/var/log/translator/translator.log",
    "fields": {
      "service": "translator",
      "version": "1.0.0",
      "environment": "production"
    },
    "error_stack": true,
    "caller": true
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `level` | string | `"info"` | Log level |
| `format` | string | `"json"` | Log format |
| `output` | string | `"stdout"` | Output destination |
| `file` | string | - | Log file path |
| `max_size` | string | `"100MB"` | Max file size |
| `max_backups` | int | `10` | Max backup files |
| `max_age` | int | `30` | Max age in days |
| `compress` | bool | `true` | Compress old files |
| `fields` | object | - | Additional fields |
| `error_stack` | bool | `true` | Include error stack |
| `caller` | bool | `true` | Include caller info |

## Worker Configuration

### Basic Worker

```json
{
  "worker": {
    "id": "worker-1",
    "coordination_server": "http://localhost:8080",
    "heartbeat_interval": "10s",
    "max_concurrent_jobs": 5,
    "job_timeout": "30m",
    "retry_attempts": 3,
    "retry_delay": "5s"
  }
}
```

### Advanced Worker

```json
{
  "worker": {
    "id": "worker-1",
    "coordination_server": "http://localhost:8080",
    "heartbeat_interval": "10s",
    "max_concurrent_jobs": 5,
    "job_timeout": "30m",
    "retry_attempts": 3,
    "retry_delay": "5s",
    "capabilities": ["epub", "fb2", "txt", "markdown"],
    "priority": "normal",
    "resources": {
      "cpu_cores": 4,
      "memory_mb": 8192,
      "gpu": false
    },
    "health_check": {
      "enabled": true,
      "interval": "30s",
      "timeout": "5s"
    }
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `id` | string | - | Worker ID |
| `coordination_server` | string | - | Coordination server URL |
| `heartbeat_interval` | duration | `"10s"` | Heartbeat interval |
| `max_concurrent_jobs` | int | `5` | Max concurrent jobs |
| `job_timeout` | duration | `"30m"` | Job timeout |
| `retry_attempts` | int | `3` | Retry attempts |
| `retry_delay` | duration | `"5s"` | Retry delay |
| `capabilities` | array | `["epub","fb2"]` | Supported formats |
| `priority` | string | `"normal"` | Worker priority |
| `resources.cpu_cores` | int | `4` | CPU cores |
| `resources.memory_mb` | int | `8192` | Memory in MB |
| `resources.gpu` | bool | `false` | GPU available |
| `health_check.enabled` | bool | `true` | Enable health checks |
| `health_check.interval` | duration | `"30s"` | Health check interval |
| `health_check.timeout` | duration | `"5s"` | Health check timeout |

## Storage Configuration

### Local Storage

```json
{
  "storage": {
    "type": "local",
    "base_path": "/var/lib/translator/files",
    "temp_path": "/tmp/translator",
    "max_file_size": "100MB",
    "cleanup_interval": "1h"
  }
}
```

### S3 Storage

```json
{
  "storage": {
    "type": "s3",
    "bucket": "translator-files",
    "region": "us-west-2",
    "access_key": "${AWS_ACCESS_KEY_ID}",
    "secret_key": "${AWS_SECRET_ACCESS_KEY}",
    "endpoint": "",
    "force_path_style": false,
    "max_file_size": "100MB",
    "multipart_threshold": "64MB",
    "multipart_chunk_size": "16MB"
  }
}
```

### Google Cloud Storage

```json
{
  "storage": {
    "type": "gcs",
    "bucket": "translator-files",
    "credentials_file": "/path/to/credentials.json",
    "project_id": "my-project",
    "max_file_size": "100MB"
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `type` | string | `"local"` | Storage type |
| `base_path` | string | - | Base path (local) |
| `temp_path` | string | `"/tmp/translator"` | Temporary path |
| `bucket` | string | - | Bucket name (S3/GCS) |
| `region` | string | - | Region (S3) |
| `access_key` | string | - | Access key (S3) |
| `secret_key` | string | - | Secret key (S3) |
| `endpoint` | string | - | Custom endpoint |
| `force_path_style` | bool | `false` | Force path style |
| `credentials_file` | string | - | Credentials file (GCS) |
| `project_id` | string | - | Project ID (GCS) |
| `max_file_size` | string | `"100MB"` | Max file size |
| `multipart_threshold` | string | `"64MB"` | Multipart threshold |
| `multipart_chunk_size` | string | `"16MB"` | Multipart chunk size |
| `cleanup_interval` | duration | `"1h"` | Cleanup interval |

## Queue Configuration

### Redis Queue

```json
{
  "queue": {
    "type": "redis",
    "connection_string": "redis://localhost:6379",
    "queue_name": "translator_jobs",
    "priority_queue_name": "translator_priority_jobs",
    "max_retries": 3,
    "retry_delay": "5s",
    "dead_letter_queue": "translator_dead_letter",
    "visibility_timeout": "30m",
    "message_retention": "7d"
  }
}
```

### RabbitMQ Queue

```json
{
  "queue": {
    "type": "rabbitmq",
    "connection_string": "amqp://guest:guest@localhost:5672/",
    "queue_name": "translator_jobs",
    "exchange": "translator_exchange",
    "routing_key": "translation.job",
    "max_retries": 3,
    "retry_delay": "5s",
    "dead_letter_queue": "translator_dead_letter",
    "prefetch_count": 10
  }
}
```

#### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `type` | string | `"redis"` | Queue type |
| `connection_string` | string | - | Connection string |
| `queue_name` | string | - | Queue name |
| `priority_queue_name` | string | - | Priority queue name |
| `exchange` | string | - | Exchange name (RabbitMQ) |
| `routing_key` | string | - | Routing key (RabbitMQ) |
| `max_retries` | int | `3` | Max retries |
| `retry_delay` | duration | `"5s"` | Retry delay |
| `dead_letter_queue` | string | - | Dead letter queue |
| `visibility_timeout` | duration | `"30m"` | Visibility timeout |
| `message_retention` | duration | `"7d"` | Message retention |
| `prefetch_count` | int | `10` | Prefetch count |

## Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `ANTHROPIC_API_KEY` | Anthropic API key | `sk-ant-...` |
| `OPENAI_API_KEY` | OpenAI API key | `sk-...` |
| `DEEPSEEK_API_KEY` | DeepSeek API key | `sk-...` |
| `JWT_SECRET` | JWT signing secret | `your-secret-key` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TRANSLATOR_CONFIG` | Config file path | `./config.json` |
| `TRANSLATOR_LOG_LEVEL` | Log level | `info` |
| `TRANSLATOR_PORT` | Server port | `8080` |
| `TRANSLATOR_HOST` | Server host | `0.0.0.0` |
| `TRANSLATOR_ENV` | Environment | `development` |
| `DATABASE_URL` | Database URL | - |
| `REDIS_URL` | Redis URL | - |
| `AWS_ACCESS_KEY_ID` | AWS access key | - |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | - |
| `GOOGLE_APPLICATION_CREDENTIALS` | GCP credentials | - |

### Database URLs

#### PostgreSQL
```
postgres://user:password@localhost:5432/translator?sslmode=disable
```

#### MySQL
```
user:password@tcp(localhost:3306)/translator?charset=utf8mb4&parseTime=True&loc=Local
```

#### SQLite
```
file:/path/to/database.db?cache=shared&mode=rwc
```

## Configuration Examples

### Development Configuration

```json
{
  "server": {
    "host": "localhost",
    "port": 8080,
    "enable_tls": false
  },
  "database": {
    "type": "sqlite",
    "connection_string": "./translator.db"
  },
  "cache": {
    "type": "memory",
    "ttl": "10m"
  },
  "llm": {
    "provider": "anthropic",
    "api_key": "${ANTHROPIC_API_KEY}",
    "model": "claude-3-sonnet-20240229",
    "max_tokens": 4096,
    "temperature": 0.3
  },
  "logging": {
    "level": "debug",
    "format": "text",
    "output": "stdout"
  },
  "security": {
    "api_key_required": false,
    "cors": {
      "enabled": true,
      "allowed_origins": ["*"]
    }
  }
}
```

### Production Configuration

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "tls_port": 8081,
    "enable_tls": true,
    "cert_file": "/etc/translator/certs/server.crt",
    "key_file": "/etc/translator/certs/server.key"
  },
  "database": {
    "type": "postgresql",
    "connection_string": "${DATABASE_URL}",
    "max_open_connections": 100,
    "max_idle_connections": 10
  },
  "cache": {
    "type": "redis",
    "connection_string": "${REDIS_URL}",
    "ttl": "1h"
  },
  "llm": {
    "provider": "anthropic",
    "api_key": "${ANTHROPIC_API_KEY}",
    "model": "claude-3-sonnet-20240229",
    "max_tokens": 4096,
    "temperature": 0.3,
    "retry_attempts": 3
  },
  "storage": {
    "type": "s3",
    "bucket": "translator-files",
    "region": "us-west-2",
    "access_key": "${AWS_ACCESS_KEY_ID}",
    "secret_key": "${AWS_SECRET_ACCESS_KEY}"
  },
  "security": {
    "api_key_required": true,
    "jwt_secret": "${JWT_SECRET}",
    "rate_limit": {
      "enabled": true,
      "requests_per_minute": 60,
      "burst": 10
    },
    "cors": {
      "enabled": true,
      "allowed_origins": ["https://translator.example.com"]
    }
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "file",
    "file": "/var/log/translator/translator.log",
    "max_size": "100MB",
    "max_backups": 10
  }
}
```

### Distributed Worker Configuration

```json
{
  "worker": {
    "id": "worker-${HOSTNAME}",
    "coordination_server": "http://coordination.example.com:8080",
    "heartbeat_interval": "10s",
    "max_concurrent_jobs": 10,
    "capabilities": ["epub", "fb2", "txt", "markdown"],
    "resources": {
      "cpu_cores": 8,
      "memory_mb": 16384
    }
  },
  "llm": {
    "provider": "anthropic",
    "api_key": "${ANTHROPIC_API_KEY}",
    "model": "claude-3-sonnet-20240229"
  },
  "storage": {
    "type": "s3",
    "bucket": "translator-files",
    "region": "us-west-2"
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  }
}
```

## Validation and Testing

### Configuration Validation

```bash
# Validate configuration file
translator-server --config config.json --validate

# Check configuration syntax
translator-server --config config.json --check-syntax

# Show effective configuration
translator-server --config config.json --show-config
```

### Environment Variable Validation

```bash
# Check required environment variables
translator-server --check-env

# Show environment variables
translator-server --show-env
```

### Configuration Testing

```bash
# Test database connection
translator-server --test-db

# Test cache connection
translator-server --test-cache

# Test LLM provider
translator-server --test-llm

# Test storage
translator-server --test-storage
```

### Configuration Schema

The system uses JSON Schema for configuration validation. The schema is available at:

```bash
# Get configuration schema
translator-server --schema > config-schema.json

# Validate against schema
ajv validate -s config-schema.json -d config.json
```

## Best Practices

### Security

1. **Use environment variables** for sensitive data (API keys, passwords)
2. **Enable TLS** in production environments
3. **Use strong JWT secrets** and rotate them regularly
4. **Enable rate limiting** to prevent abuse
5. **Restrict CORS origins** to specific domains

### Performance

1. **Configure appropriate connection pools** for database and cache
2. **Use Redis cache** for production deployments
3. **Enable compression** for large files
4. **Configure timeouts** appropriately
5. **Monitor resource usage** and adjust limits

### Reliability

1. **Enable auto-migration** for database schema
2. **Configure retry logic** for external services
3. **Use dead letter queues** for failed jobs
4. **Enable health checks** for monitoring
5. **Configure logging** for troubleshooting

### Scalability

1. **Use distributed mode** for high throughput
2. **Configure load balancing** for multiple instances
3. **Use external storage** (S3, GCS) for file storage
4. **Enable horizontal scaling** with container orchestration
5. **Monitor queue depth** and scale workers accordingly

## Conclusion

This configuration reference covers all available options for the Universal Multi-Format Multi-Language Ebook Translation System. Use the examples as starting points and customize them according to your specific requirements.

For additional help, refer to the deployment guide or contact the support team.