# Distributed Work Documentation

This document describes the distributed work functionality that allows the Russian-Serbian FB2 translation system to utilize remote worker machines for increased translation capacity and LLM provider diversity.

## Overview

The distributed work system enables a "root" translation server to discover, pair with, and utilize remote worker instances via secure HTTP/3 connections. This provides:

- **Horizontal scaling**: Distribute translation workload across multiple machines
- **Provider diversity**: Access different LLM providers and local models on remote machines
- **Resource optimization**: Utilize remote Ollama/llama.cpp instances and API keys
- **Fault tolerance**: Automatic failover and load balancing across workers

## Architecture

### Components

1. **SSH Pool Manager** (`pkg/distributed/ssh_pool.go`)
   - Manages SSH connections to remote worker machines
   - Handles connection pooling, retries, and cleanup
   - Supports key-based and password authentication

2. **Pairing Manager** (`pkg/distributed/pairing.go`)
   - Discovers remote translator services via SSH
   - Performs health checks and capability detection
   - Manages pairing/unpairing lifecycle

3. **Distributed Coordinator** (`pkg/distributed/coordinator.go`)
   - Coordinates translation requests across remote instances
   - Implements load balancing and failover logic
   - Tracks remote instance availability and capacity

4. **Distributed Manager** (`pkg/distributed/manager.go`)
   - High-level orchestration of distributed work
   - Integrates with existing event system and WebSocket hub
   - Provides unified API for distributed operations

### Communication Flow

```
Root Server ←──HTTP/3/QUIC──→ Remote Worker
     │                              │
     ├──SSH Discovery               ├──Service Health Check
     ├──Capability Query            ├──LLM Instance Management
     └──Translation Requests        └──Result Streaming
```

## Configuration

### Root Server Configuration

```json
{
  "distributed": {
    "enabled": true,
    "workers": {
      "worker-1": {
        "name": "Production Worker 1",
        "host": "worker1.example.com",
        "port": 22,
        "user": "translator",
        "key_file": "/path/to/ssh/key",
        "max_capacity": 10,
        "enabled": true,
        "tags": ["gpu", "high-memory"]
      }
    },
    "ssh_timeout": 30,
    "ssh_max_retries": 3,
    "health_check_interval": 30,
    "max_remote_instances": 20
  }
}
```

### Worker Configuration

Workers run the same translator server but with distributed features disabled:

```json
{
  "distributed": {
    "enabled": false
  },
  "translation": {
    "max_concurrent": 10
  }
}
```

### Environment Variables

```bash
# SSH Authentication (for root server)
export WORKER_SSH_KEY_PATH="/path/to/ssh/private/key"

# API Keys (configure on workers as needed)
export OPENAI_API_KEY="..."
export ANTHROPIC_API_KEY="..."
export ZHIPU_API_KEY="..."
export DEEPSEEK_API_KEY="..."

# Local LLM (on workers)
export OLLAMA_ENABLED=true
export OLLAMA_MODEL="llama3:8b"
```

## API Endpoints

### Distributed Status
```http
GET /api/v1/distributed/status
```

Returns comprehensive status of all workers and remote instances.

**Response:**
```json
{
  "initialized": true,
  "enabled": true,
  "workers": {
    "worker-1": {
      "name": "Production Worker 1",
      "enabled": true,
      "status": "paired",
      "capacity": 10
    }
  },
  "active_connections": 2,
  "remote_instances": 8,
  "paired_workers": 1
}
```

### Discover Workers
```http
POST /api/v1/distributed/workers/discover
```

Automatically discovers and pairs with all configured workers.

### Pair/Unpair Workers
```http
POST   /api/v1/distributed/workers/{worker_id}/pair
DELETE /api/v1/distributed/workers/{worker_id}/pair
```

Manually pair or unpair specific workers.

### Distributed Translation
```http
POST /api/v1/distributed/translate
```

Perform translation using distributed resources.

**Request:**
```json
{
  "text": "Hello world",
  "context_hint": "greeting",
  "session_id": "optional-session-id"
}
```

**Response:**
```json
{
  "translated_text": "Zdravo svete",
  "session_id": "generated-or-provided-session-id"
}
```

## Worker Management

### Adding Workers

Workers can be added dynamically via API or configuration:

```bash
# Via API (POST /api/v1/distributed/workers)
curl -X POST https://your-server:8443/api/v1/distributed/workers \
  -H "Content-Type: application/json" \
  -d '{
    "id": "new-worker",
    "name": "New GPU Worker",
    "host": "gpu-worker.example.com",
    "port": 22,
    "user": "translator",
    "max_capacity": 20,
    "tags": ["gpu", "high-memory"]
  }'
```

### Worker Requirements

Remote workers must:

1. **Run translator server**: Same binary with worker configuration
2. **Be SSH accessible**: From root server with proper authentication
3. **Have TLS certificates**: For HTTPS communication
4. **Expose API endpoints**: Standard REST API must be accessible
5. **Have LLM providers**: API keys or local models configured

### Capacity Management

The system automatically determines how many LLM instances to run based on:

- **Provider priority**: API keys (10) > OAuth (5) > free/local (1)
- **Host capacity**: `max_capacity` setting per worker
- **Resource availability**: Automatic detection of available providers

## Security Considerations

### SSH Security
- Use key-based authentication (avoid passwords)
- Restrict SSH user permissions on workers
- Regularly rotate SSH keys
- Use SSH bastion hosts for complex networks

### Network Security
- All communication uses HTTPS with TLS 1.3
- HTTP/3 provides additional security benefits
- Implement proper firewall rules
- Use VPNs for untrusted networks

### API Security
- Workers can disable authentication for simplified deployment
- Root server handles authentication and authorization
- API keys should be properly scoped
- Implement rate limiting per worker

## Monitoring and Events

### WebSocket Events

All distributed operations emit WebSocket events:

```javascript
// Worker discovery and pairing
{
  "type": "distributed_worker_discovered",
  "session_id": "system",
  "message": "Discovered and paired with worker worker-1",
  "data": {
    "worker_id": "worker-1",
    "capabilities": {...}
  }
}

// Translation operations
{
  "type": "distributed_translation_attempt",
  "session_id": "user-session-123",
  "message": "Attempting distributed translation with remote-instance-1",
  "data": {
    "instance_id": "remote-instance-1",
    "worker_id": "worker-1"
  }
}

// Health and status
{
  "type": "distributed_worker_offline",
  "session_id": "system",
  "message": "Remote worker worker-1 went offline",
  "data": {
    "worker_id": "worker-1",
    "error": "connection timeout"
  }
}
```

### Health Checks

- **SSH connectivity**: Verified during discovery
- **Service health**: HTTP health checks every 30 seconds
- **Instance availability**: Monitored during translation attempts
- **Capacity limits**: Automatically adjusted based on failures

## Docker Deployment

### Development Setup

```bash
# Start distributed test environment
docker-compose -f docker-compose.distributed.yml up -d

# View logs
docker-compose -f docker-compose.distributed.yml logs -f

# Run distributed tests
docker-compose -f docker-compose.distributed.yml exec distributed-tests go test -v ./test/distributed/...
```

### Production Deployment

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  translator-main:
    image: translator:latest
    environment:
      - CONFIG_FILE=config.distributed.json
    volumes:
      - ./config:/app/config:ro
      - ./ssh-keys:/app/ssh:ro
    ports:
      - "8443:8443"

  translator-worker-1:
    image: translator:latest
    environment:
      - CONFIG_FILE=config.worker.json
    deploy:
      placement:
        constraints:
          - node.labels.gpu == true
```

## Troubleshooting

### Common Issues

1. **SSH Connection Failures**
   - Verify SSH key permissions (`chmod 600 key_file`)
   - Check SSH user permissions on worker
   - Ensure SSH service is running on worker

2. **Service Discovery Failures**
   - Verify translator service is running on worker
   - Check HTTPS certificates are valid
   - Ensure firewall allows HTTPS traffic

3. **Translation Failures**
   - Check worker LLM provider configuration
   - Verify API keys are set on workers
   - Monitor worker resource usage

4. **WebSocket Event Issues**
   - Ensure root server can reach worker WebSocket endpoints
   - Check event bus configuration
   - Verify session ID consistency

### Debugging Commands

```bash
# Test SSH connection
ssh -i /path/to/key translator@worker.example.com

# Test service health
curl -k https://worker.example.com:8443/health

# Check worker logs
docker-compose logs translator-worker-1

# Monitor distributed status
curl https://main-server:8443/api/v1/distributed/status
```

## Performance Optimization

### Load Balancing
- Round-robin distribution across remote instances
- Priority-based instance selection (API keys first)
- Automatic failover on instance failures

### Caching Strategy
- Translation results cached across distributed instances
- Worker capability information cached
- SSH connection pooling for reduced latency

### Resource Management
- Automatic instance scaling based on demand
- Capacity limits prevent resource exhaustion
- Health-based instance deactivation

## Future Enhancements

- **Kubernetes integration**: Native K8s operator for worker management
- **Auto-scaling**: Dynamic worker provisioning based on load
- **Advanced routing**: Content-based routing to specialized workers
- **Federated learning**: Distributed model training and fine-tuning
- **Multi-region support**: Geographic distribution for reduced latency</content>
</xai:function_call">