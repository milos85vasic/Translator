# Troubleshooting Guide - Universal Multi-Format Multi-Language Ebook Translation System

## Table of Contents

1. [Overview](#overview)
2. [Common Issues](#common-issues)
3. [Installation Problems](#installation-problems)
4. [Configuration Issues](#configuration-issues)
5. [Database Problems](#database-problems)
6. [LLM Provider Issues](#llm-provider-issues)
7. [Performance Problems](#performance-problems)
8. [Network Issues](#network-issues)
9. [File Processing Issues](#file-processing-issues)
10. [Distributed System Issues](#distributed-system-issues)
11. [Security Issues](#security-issues)
12. [Monitoring and Debugging](#monitoring-and-debugging)
13. [Recovery Procedures](#recovery-procedures)

## Overview

This troubleshooting guide covers common issues, their causes, and solutions for the Universal Multi-Format Multi-Language Ebook Translation System. Each section includes symptoms, diagnosis steps, and resolution procedures.

### Quick Diagnosis Checklist

- [ ] Check system logs for error messages
- [ ] Verify configuration files are valid
- [ ] Test network connectivity
- [ ] Check resource usage (CPU, memory, disk)
- [ ] Verify external service availability
- [ ] Run health checks

## Common Issues

### Service Won't Start

**Symptoms:**
- Service fails to start
- Immediate exit after startup
- No response on configured ports

**Diagnosis:**
```bash
# Check service status
systemctl status translator-server
journalctl -u translator-server -f

# Check configuration
translator-server --config /etc/translator/config.json --validate

# Check port availability
netstat -tlnp | grep :8080
ss -tuln | grep :8080

# Check permissions
ls -la /etc/translator/
ls -la /var/log/translator/
```

**Common Causes and Solutions:**

1. **Invalid Configuration**
   ```bash
   # Validate configuration
   translator-server --config config.json --validate
   
   # Check syntax
   cat config.json | jq .
   ```

2. **Port Already in Use**
   ```bash
   # Find process using port
   lsof -i :8080
   
   # Kill process
   kill -9 <PID>
   
   # Or change port in config
   ```

3. **Permission Issues**
   ```bash
   # Fix ownership
   sudo chown -R translator:translator /etc/translator /var/log/translator /var/lib/translator
   
   # Fix permissions
   sudo chmod 755 /etc/translator
   sudo chmod 644 /etc/translator/config.json
   ```

4. **Missing Dependencies**
   ```bash
   # Check required files
   ls -la /etc/translator/certs/
   ls -la /var/lib/translator/
   
   # Create missing directories
   sudo mkdir -p /var/lib/translator
   sudo chown translator:translator /var/lib/translator
   ```

### High Memory Usage

**Symptoms:**
- System becomes unresponsive
- Out of memory errors
- Service restarts frequently

**Diagnosis:**
```bash
# Check memory usage
top -p $(pgrep translator-server)
htop
free -h

# Check for memory leaks
valgrind --tool=memcheck --leak-check=full translator-server

# Monitor over time
watch -n 5 'ps aux | grep translator-server'
```

**Solutions:**

1. **Reduce Connection Limits**
   ```json
   {
     "database": {
       "max_open_connections": 25,
       "max_idle_connections": 5
     },
     "cache": {
       "max_size": 500
     }
   }
   ```

2. **Enable Garbage Collection Tuning**
   ```bash
   # Set GOGC environment variable
   export GOGC=50
   ```

3. **Increase System Memory**
   ```bash
   # Add swap space
   sudo fallocate -l 2G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

### High CPU Usage

**Symptoms:**
- System slow response
- High CPU utilization
- Fan running constantly

**Diagnosis:**
```bash
# Check CPU usage
top -p $(pgrep translator-server)
htop
iostat 1

# Profile CPU usage
go tool pprof http://localhost:6060/debug/pprof/profile
```

**Solutions:**

1. **Optimize Configuration**
   ```json
   {
     "server": {
       "max_connections": 1000,
       "read_timeout": "30s",
       "write_timeout": "30s"
     }
   }
   ```

2. **Enable Caching**
   ```json
   {
     "cache": {
       "type": "redis",
       "ttl": "1h"
     }
   }
   ```

3. **Scale Horizontally**
   ```bash
   # Deploy multiple instances
   docker-compose up -d --scale translator-server=3
   ```

## Installation Problems

### Build Failures

**Symptoms:**
- Compilation errors
- Missing dependencies
- Version conflicts

**Diagnosis:**
```bash
# Check Go version
go version

# Check module dependencies
go mod tidy
go mod verify

# Check for issues
go vet ./...
```

**Solutions:**

1. **Update Go Version**
   ```bash
   # Install required Go version
   go install golang.org/dl/go1.25.2@latest
   go1.25.2 download
   ```

2. **Clean and Rebuild**
   ```bash
   # Clean build cache
   go clean -cache
   go clean -modcache
   
   # Rebuild
   go build ./cmd/cli
   go build ./cmd/server
   ```

3. **Fix Module Issues**
   ```bash
   # Update dependencies
   go get -u ./...
   go mod tidy
   ```

### Docker Issues

**Symptoms:**
- Container won't start
- Build failures
- Network issues

**Diagnosis:**
```bash
# Check Docker status
docker info
docker version

# Check container logs
docker logs translator-server
docker-compose logs translator-server

# Check network
docker network ls
docker network inspect translator_default
```

**Solutions:**

1. **Fix Docker Build**
   ```bash
   # Clean Docker cache
   docker system prune -a
   
   # Rebuild
   docker-compose build --no-cache
   ```

2. **Fix Network Issues**
   ```bash
   # Recreate network
   docker network rm translator_default
   docker-compose up -d
   ```

3. **Fix Volume Permissions**
   ```bash
   # Check volume permissions
   docker exec translator-server ls -la /var/lib/translator
   
   # Fix permissions
   sudo chown -R 1000:1000 ./data
   ```

## Configuration Issues

### Invalid Configuration

**Symptoms:**
- Service won't start
- Configuration validation errors
- Unexpected behavior

**Diagnosis:**
```bash
# Validate configuration
translator-server --config config.json --validate

# Check syntax
cat config.json | jq .

# Show effective configuration
translator-server --config config.json --show-config
```

**Solutions:**

1. **Fix JSON Syntax**
   ```bash
   # Use jq to validate
   cat config.json | jq .
   
   # Fix common issues
   # - Missing commas
   # - Trailing commas
   # - Incorrect quotes
   ```

2. **Use Configuration Template**
   ```bash
   # Generate default config
   translator-server --generate-config > config.json
   
   # Edit as needed
   ```

3. **Check Environment Variables**
   ```bash
   # Show environment variables
   env | grep TRANSLATOR_
   
   # Test variable substitution
   translator-server --test-env
   ```

### Missing Required Settings

**Symptoms:**
- API key errors
- Database connection failures
- Authentication issues

**Diagnosis:**
```bash
# Check required environment variables
translator-server --check-env

# Test external connections
translator-server --test-db
translator-server --test-cache
translator-server --test-llm
```

**Solutions:**

1. **Set API Keys**
   ```bash
   # Set environment variables
   export ANTHROPIC_API_KEY="your-key-here"
   export OPENAI_API_KEY="your-key-here"
   
   # Or use .env file
   echo "ANTHROPIC_API_KEY=your-key-here" >> .env
   ```

2. **Configure Database**
   ```json
   {
     "database": {
       "type": "postgresql",
       "connection_string": "postgres://user:password@localhost:5432/translator"
     }
   }
   ```

3. **Fix File Paths**
   ```bash
   # Check file existence
   ls -la /etc/translator/certs/server.crt
   ls -la /var/lib/translator/
   
   # Create missing files/directories
   sudo mkdir -p /etc/translator/certs
   sudo touch /etc/translator/certs/server.crt
   ```

## Database Problems

### Connection Failures

**Symptoms:**
- Database connection errors
- Timeouts
- Authentication failures

**Diagnosis:**
```bash
# Test database connection
translator-server --test-db

# Check database status
pg_isready -h localhost -p 5432
mysql -h localhost -u root -p

# Check network connectivity
telnet localhost 5432
nc -zv localhost 5432
```

**Solutions:**

1. **Fix Connection String**
   ```json
   {
     "database": {
       "connection_string": "postgres://user:password@localhost:5432/translator?sslmode=disable"
     }
   }
   ```

2. **Check Database Service**
   ```bash
   # PostgreSQL
   sudo systemctl status postgresql
   sudo systemctl start postgresql
   
   # MySQL
   sudo systemctl status mysql
   sudo systemctl start mysql
   ```

3. **Fix Firewall Rules**
   ```bash
   # Allow database port
   sudo ufw allow 5432/tcp
   sudo iptables -A INPUT -p tcp --dport 5432 -j ACCEPT
   ```

### Performance Issues

**Symptoms:**
- Slow queries
- High database load
- Timeouts

**Diagnosis:**
```bash
# Check database connections
SELECT * FROM pg_stat_activity WHERE datname = 'translator';

# Check slow queries
SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;

# Check database size
SELECT pg_size_pretty(pg_database_size('translator'));
```

**Solutions:**

1. **Optimize Database Configuration**
   ```sql
   -- PostgreSQL
   ALTER SYSTEM SET shared_buffers = '256MB';
   ALTER SYSTEM SET effective_cache_size = '1GB';
   ALTER SYSTEM SET maintenance_work_mem = '64MB';
   SELECT pg_reload_conf();
   ```

2. **Create Indexes**
   ```sql
   CREATE INDEX idx_translations_created_at ON translations(created_at);
   CREATE INDEX idx_files_status ON files(status);
   CREATE INDEX idx_jobs_status ON jobs(status);
   ```

3. **Configure Connection Pool**
   ```json
   {
     "database": {
       "max_open_connections": 100,
       "max_idle_connections": 10,
       "connection_max_lifetime": "1h"
     }
   }
   ```

### Migration Issues

**Symptoms:**
- Migration failures
- Schema conflicts
- Data corruption

**Diagnosis:**
```bash
# Check migration status
translator-server --migrate-status

# Run migrations manually
translator-server --migrate-up

# Check database schema
\dt translator.*
```

**Solutions:**

1. **Fix Migration Conflicts**
   ```bash
   # Rollback migrations
   translator-server --migrate-down
   
   # Re-run migrations
   translator-server --migrate-up
   ```

2. **Backup and Restore**
   ```bash
   # Backup database
   pg_dump translator > backup.sql
   
   # Restore if needed
   psql translator < backup.sql
   ```

## LLM Provider Issues

### API Key Problems

**Symptoms:**
- Authentication errors
- 401 Unauthorized
- API quota exceeded

**Diagnosis:**
```bash
# Test API key
curl -H "Authorization: Bearer $ANTHROPIC_API_KEY" \
  https://api.anthropic.com/v1/messages

# Check API usage
translator-server --test-llm
```

**Solutions:**

1. **Verify API Key**
   ```bash
   # Check environment variable
   echo $ANTHROPIC_API_KEY
   
   # Test with curl
   curl -H "Authorization: Bearer $ANTHROPIC_API_KEY" \
     https://api.anthropic.com/v1/messages \
     -d '{"model": "claude-3-sonnet-20240229", "max_tokens": 10, "messages": [{"role": "user", "content": "test"}]}'
   ```

2. **Check Quota and Billing**
   ```bash
   # Check Anthropic usage
   # Visit https://console.anthropic.com/
   
   # Check OpenAI usage
   # Visit https://platform.openai.com/usage
   ```

3. **Use Multiple Providers**
   ```json
   {
     "llm": {
       "provider": "anthropic",
       "fallback_providers": ["openai", "deepseek"]
     }
   }
   ```

### Rate Limiting

**Symptoms:**
- 429 Too Many Requests
- Slow responses
- Request timeouts

**Diagnosis:**
```bash
# Check rate limits
curl -I https://api.anthropic.com/v1/messages

# Monitor response times
curl -w "@curl-format.txt" -H "Authorization: Bearer $ANTHROPIC_API_KEY" \
  https://api.anthropic.com/v1/messages
```

**Solutions:**

1. **Configure Retry Logic**
   ```json
   {
     "llm": {
       "retry_attempts": 5,
       "retry_delay": "2s",
       "exponential_backoff": true
     }
   }
   ```

2. **Implement Rate Limiting**
   ```json
   {
     "llm": {
       "rate_limit": {
         "requests_per_minute": 50,
         "burst": 10
       }
     }
   }
   ```

3. **Use Multiple API Keys**
   ```json
   {
     "llm": {
       "api_keys": [
         "${ANTHROPIC_API_KEY_1}",
         "${ANTHROPIC_API_KEY_2}",
         "${ANTHROPIC_API_KEY_3}"
       ],
       "key_rotation": true
     }
   }
   ```

### Model Availability

**Symptoms:**
- Model not found errors
- Unsupported model
- Model deprecated

**Diagnosis:**
```bash
# List available models
curl -H "Authorization: Bearer $ANTHROPIC_API_KEY" \
  https://api.anthropic.com/v1/models

# Test specific model
translator-server --test-llm --model claude-3-sonnet-20240229
```

**Solutions:**

1. **Update Model Name**
   ```json
   {
     "llm": {
       "model": "claude-3-5-sonnet-20241022"
     }
   }
   ```

2. **Use Model Fallbacks**
   ```json
   {
     "llm": {
       "model": "claude-3-sonnet-20240229",
       "fallback_models": [
         "claude-3-haiku-20240307",
         "gpt-4",
         "gpt-3.5-turbo"
       ]
     }
   }
   ```

## Performance Problems

### Slow Translation

**Symptoms:**
- Long processing times
- Poor throughput
- User complaints

**Diagnosis:**
```bash
# Check system resources
top
htop
iostat 1

# Profile application
go tool pprof http://localhost:6060/debug/pprof/profile

# Check queue depth
curl http://localhost:8080/metrics | grep queue_depth
```

**Solutions:**

1. **Optimize LLM Parameters**
   ```json
   {
     "llm": {
       "max_tokens": 2048,
       "temperature": 0.1,
       "timeout": "30s"
     }
   }
   ```

2. **Enable Caching**
   ```json
   {
     "cache": {
       "type": "redis",
       "ttl": "24h"
     }
   }
   ```

3. **Scale Workers**
   ```bash
   # Add more workers
   docker-compose up -d --scale worker=5
   
   # Or configure distributed mode
   ```

### Memory Leaks

**Symptoms:**
- Memory usage increases over time
- Service crashes
- Out of memory errors

**Diagnosis:**
```bash
# Monitor memory usage
watch -n 5 'ps aux | grep translator-server'

# Check for leaks
valgrind --tool=memcheck --leak-check=full translator-server

# Go memory profiling
go tool pprof http://localhost:6060/debug/pprof/heap
```

**Solutions:**

1. **Update Dependencies**
   ```bash
   go get -u ./...
   go mod tidy
   ```

2. **Configure GC**
   ```bash
   export GOGC=50
   export GOMEMLIMIT=1GiB
   ```

3. **Restart Service Periodically**
   ```bash
   # Add to systemd service
   Restart=always
   RestartSec=3600
   ```

## Network Issues

### Connection Timeouts

**Symptoms:**
- Request timeouts
- Connection refused
- Slow responses

**Diagnosis:**
```bash
# Test connectivity
curl -v http://localhost:8080/health
telnet localhost 8080

# Check network configuration
ip addr show
netstat -tlnp

# Check firewall
sudo ufw status
sudo iptables -L
```

**Solutions:**

1. **Adjust Timeouts**
   ```json
   {
     "server": {
       "read_timeout": "60s",
       "write_timeout": "60s",
       "idle_timeout": "120s"
     }
   }
   ```

2. **Fix Firewall Rules**
   ```bash
   # Allow required ports
   sudo ufw allow 8080/tcp
   sudo ufw allow 8081/tcp
   sudo ufw allow 8082/tcp
   ```

3. **Configure Load Balancer**
   ```nginx
   upstream translator {
     server translator1:8080 max_fails=3 fail_timeout=30s;
     server translator2:8080 max_fails=3 fail_timeout=30s;
   }
   ```

### SSL/TLS Issues

**Symptoms:**
- Certificate errors
- Handshake failures
- HTTPS not working

**Diagnosis:**
```bash
# Test SSL certificate
openssl s_client -connect localhost:8081 -servername localhost

# Check certificate validity
openssl x509 -in /etc/translator/certs/server.crt -text -noout

# Test with curl
curl -v https://localhost:8081/health
```

**Solutions:**

1. **Generate Valid Certificate**
   ```bash
   # Self-signed for testing
   openssl req -x509 -newkey rsa:4096 -keyout server.key \
     -out server.crt -days 365 -nodes
   
   # Let's Encrypt for production
   sudo certbot certonly --standalone -d translator.example.com
   ```

2. **Fix Certificate Paths**
   ```json
   {
     "server": {
       "cert_file": "/etc/translator/certs/server.crt",
       "key_file": "/etc/translator/certs/server.key"
     }
   }
   ```

3. **Configure TLS Settings**
   ```json
   {
     "server": {
       "tls": {
         "min_version": "1.2",
         "cipher_suites": [
           "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
         ]
       }
     }
   }
   ```

## File Processing Issues

### Unsupported Formats

**Symptoms:**
- File format not supported
- Parsing errors
- Corrupted output

**Diagnosis:**
```bash
# Check file format
file input.epub
file input.fb2

# Test with CLI
translator translate --input input.epub --output output.epub --dry-run

# Check logs
grep -i "format" /var/log/translator/translator.log
```

**Solutions:**

1. **Check Supported Formats**
   ```bash
   # List supported formats
   translator --help | grep -i format
   
   # Test specific format
   translator translate --input test.epub --format epub
   ```

2. **Convert File Format**
   ```bash
   # Convert to supported format
   pandoc -input.docx -o epub.epub
   
   # Or use calibre
   ebook-convert input.docx output.epub
   ```

3. **Update Format Handlers**
   ```bash
   # Check for updates
   go get -u ./...
   
   # Rebuild
   go build ./cmd/cli
   ```

### Large File Issues

**Symptoms:**
- Memory exhaustion
- Processing timeouts
- Upload failures

**Diagnosis:**
```bash
# Check file size
ls -lh large_file.epub

# Monitor memory during processing
top -p $(pgrep translator)

# Check upload limits
curl -I http://localhost:8080/upload
```

**Solutions:**

1. **Increase Limits**
   ```json
   {
     "server": {
       "max_request_size": 1048576000,
       "read_timeout": "300s"
     },
     "storage": {
       "max_file_size": "1GB"
     }
   }
   ```

2. **Stream Processing**
   ```json
   {
     "processing": {
       "stream_large_files": true,
       "chunk_size": "10MB"
     }
   }
   ```

3. **Use External Storage**
   ```json
   {
     "storage": {
       "type": "s3",
       "multipart_threshold": "64MB"
     }
   }
   ```

## Distributed System Issues

### Worker Registration Failures

**Symptoms:**
- Workers not connecting
- Coordination server not seeing workers
- Load balancing not working

**Diagnosis:**
```bash
# Check coordination server
curl http://coordination:8080/workers

# Check worker logs
docker logs worker-1

# Test connectivity
telnet coordination 8080
```

**Solutions:**

1. **Fix Network Configuration**
   ```json
   {
     "worker": {
       "coordination_server": "http://coordination:8080"
     }
   }
   ```

2. **Check Firewall Rules**
   ```bash
   # Allow coordination port
   sudo ufw allow 8080/tcp
   ```

3. **Verify Worker Configuration**
   ```json
   {
     "worker": {
       "id": "worker-1",
       "heartbeat_interval": "10s",
       "capabilities": ["epub", "fb2"]
     }
   }
   ```

### Job Distribution Issues

**Symptoms:**
- Jobs not distributed
- Uneven load
- Jobs stuck in queue

**Diagnosis:**
```bash
# Check queue depth
curl http://localhost:8080/metrics | grep queue

# Check worker status
curl http://localhost:8080/workers

# Monitor job processing
tail -f /var/log/translator/translator.log | grep job
```

**Solutions:**

1. **Configure Load Balancing**
   ```json
   {
     "coordination": {
       "load_balancing": "least_connections"
     }
   }
   ```

2. **Adjust Worker Capacity**
   ```json
   {
     "worker": {
       "max_concurrent_jobs": 10
     }
   }
   ```

3. **Monitor Queue Health**
   ```bash
   # Check Redis queue
   redis-cli llen translator_jobs
   
   # Check for stuck jobs
   redis-cli lrange translator_jobs 0 10
   ```

## Security Issues

### Authentication Failures

**Symptoms:**
- 401 Unauthorized errors
- API key not accepted
- JWT token invalid

**Diagnosis:**
```bash
# Test API key
curl -H "X-API-Key: your-key" http://localhost:8080/health

# Check JWT token
curl -H "Authorization: Bearer token" http://localhost:8080/health

# Verify configuration
translator-server --show-config | grep -A 5 security
```

**Solutions:**

1. **Verify API Key**
   ```bash
   # Check environment variable
   echo $TRANSLATOR_API_KEY
   
   # Generate new key
   openssl rand -hex 32
   ```

2. **Fix JWT Configuration**
   ```json
   {
     "security": {
       "jwt_secret": "${JWT_SECRET}",
       "jwt_expiry": "24h"
     }
   }
   ```

3. **Disable Authentication (Development Only)**
   ```json
   {
     "security": {
       "api_key_required": false
     }
   }
   ```

### CORS Issues

**Symptoms:**
- Cross-origin requests blocked
- Browser console errors
- Frontend not working

**Diagnosis:**
```bash
# Test CORS
curl -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -X OPTIONS http://localhost:8080/api/translate

# Check response headers
curl -I http://localhost:8080/api/translate
```

**Solutions:**

1. **Configure CORS Properly**
   ```json
   {
     "security": {
       "cors": {
         "enabled": true,
         "allowed_origins": ["https://yourdomain.com"],
         "allowed_methods": ["GET", "POST", "PUT", "DELETE"],
         "allowed_headers": ["*"]
       }
     }
   }
   ```

2. **Use Wildcard for Development**
   ```json
   {
     "security": {
       "cors": {
         "allowed_origins": ["*"]
       }
     }
   }
   ```

## Monitoring and Debugging

### Health Checks

**Basic Health Check**
```bash
curl http://localhost:8080/health

# Detailed health check
curl http://localhost:8080/health/deep

# Component status
curl http://localhost:8080/health/components
```

### Metrics Collection

**Prometheus Metrics**
```bash
# Get metrics
curl http://localhost:9090/metrics

# Key metrics to watch
curl http://localhost:9090/metrics | grep -E "(http_requests|job_queue|translation_duration)"
```

### Log Analysis

**Common Log Patterns**
```bash
# Error logs
grep -i error /var/log/translator/translator.log

# Slow requests
grep "slow request" /var/log/translator/translator.log

# Failed translations
grep "translation failed" /var/log/translator/translator.log

# Database issues
grep -i database /var/log/translator/translator.log
```

### Debug Mode

**Enable Debug Logging**
```json
{
  "logging": {
    "level": "debug",
    "format": "text"
  }
}
```

**Enable Profiling**
```bash
# Enable pprof
export TRANSLATOR_ENABLE_PPROF=true

# Get CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile

# Get memory profile
go tool pprof http://localhost:6060/debug/pprof/heap
```

## Recovery Procedures

### Service Recovery

**Automatic Recovery**
```bash
# Systemd will automatically restart
sudo systemctl status translator-server

# Check restart count
systemctl show translator-server -p RestartCount
```

**Manual Recovery**
```bash
# Restart service
sudo systemctl restart translator-server

# Clear cache if needed
redis-cli flushall

# Reset database connections
sudo systemctl restart postgresql
```

### Data Recovery

**Database Recovery**
```bash
# Restore from backup
pg_dump translator > backup_$(date +%Y%m%d).sql
psql translator < backup_20231101.sql

# Point-in-time recovery (PostgreSQL)
pg_basebackup -h localhost -D /backup/base -U replication
```

**File Recovery**
```bash
# Restore from S3 backup
aws s3 sync s3://translator-backups/20231101/ /var/lib/translator/

# Verify file integrity
md5sum /var/lib/translator/*
```

### Disaster Recovery

**Full System Recovery**
```bash
# 1. Restore configuration
sudo cp /backup/config.json /etc/translator/

# 2. Restore database
psql translator < /backup/database.sql

# 3. Restore files
sudo rsync -av /backup/files/ /var/lib/translator/

# 4. Restart services
sudo systemctl restart translator-server
```

## Getting Help

### Support Channels

1. **Documentation**: Check this guide and other documentation
2. **Logs**: Always check application logs first
3. **Community**: GitHub Issues and Discussions
4. **Professional Support**: Contact support team for enterprise customers

### Reporting Issues

When reporting issues, include:

1. **System Information**
   ```bash
   uname -a
   go version
   docker --version
   ```

2. **Configuration**
   ```bash
   translator-server --show-config
   ```

3. **Logs**
   ```bash
   journalctl -u translator-server --since "1 hour ago"
   ```

4. **Error Description**
   - What happened
   - Expected behavior
   - Steps to reproduce

### Performance Tuning

**Quick Performance Checks**
```bash
# Check system resources
free -h
df -h
top

# Check application metrics
curl http://localhost:9090/metrics

# Database performance
SELECT * FROM pg_stat_activity WHERE datname = 'translator';
```

**Common Optimizations**

1. **Enable Redis cache**
2. **Use SSD storage**
3. **Increase memory allocation**
4. **Configure connection pooling**
5. **Enable compression**
6. **Use CDN for static files**

## Conclusion

This troubleshooting guide covers the most common issues and their solutions. For complex problems or issues not covered here, please refer to the additional documentation or contact the support team.

Remember to:
- Always check logs first
- Test configuration changes in development
- Keep backups of important data
- Monitor system performance regularly
- Keep software updated