# Universal Ebook Translator CLI
## Complete User Manual

---

## 📚 Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Command Reference](#command-reference)
4. [Configuration](#configuration)
5. [Advanced Usage](#advanced-usage)
6. [Examples](#examples)
7. [Troubleshooting](#troubleshooting)
8. [Tips & Best Practices](#tips--best-practices)

---

## Installation

### Prerequisites
- Go 1.25.2 or later
- Git
- 2GB+ RAM
- 1GB+ disk space

### Install from Source
```bash
# Clone the repository
git clone https://github.com/translator/universal.git
cd universal

# Build the CLI
make build

# Install globally
sudo cp build/translator /usr/local/bin/
```

### Install Binary Release
```bash
# Download latest release
wget https://github.com/translator/universal/releases/latest/translator-linux-amd64

# Make executable
chmod +x translator-linux-amd64

# Install globally
sudo mv translator-linux-amd64 /usr/local/bin/translator
```

### Verify Installation
```bash
translator --version
translator --help
```

---

## Quick Start

### Basic Text Translation
```bash
# Translate a simple text
translator translate "Hello, world!" --target es --provider openai
```

### Translate an Ebook
```bash
# Translate EPUB to Spanish
translator translate book.epub --target es --output book_es.epub

# Use specific provider and model
translator translate book.epub --target fr --provider anthropic --model claude-3-sonnet
```

### Batch Directory Translation
```bash
# Translate all books in a directory
translator translate /path/to/books --target de --recursive --parallel
```

### Generate Configuration
```bash
# Create default config file
translator config --init

# Edit configuration
nano ~/.translator/config.json
```

---

## Command Reference

### Global Options

| Option | Short | Description |
|--------|--------|-------------|
| `--config` | `-c` | Path to configuration file |
| `--verbose` | `-v` | Enable verbose output |
| `--quiet` | `-q` | Suppress non-error output |
| `--help` | `-h` | Show help information |
| `--version` | `-V` | Show version information |

### Main Commands

#### `translate`
Translate text or files between languages.

**Syntax:**
```bash
translator translate [input] [options]
```

**Options:**
| Option | Short | Type | Default | Description |
|--------|--------|-------|---------|-------------|
| `--target` | `-t` | string | Required | Target language code |
| `--source` | `-s` | string | auto | Source language code |
| `--provider` | `-p` | string | openai | Translation provider |
| `--model` | `-m` | string | default | Model name |
| `--output` | `-o` | string | auto | Output file path |
| `--format` | `-f` | string | auto | Output format |
| `--script` | | string | latin | Script conversion |
| `--context` | | string | | Translation context |
| `--recursive` | `-r` | flag | false | Process directories recursively |
| `--parallel` | | flag | false | Enable parallel processing |
| `--max-concurrency` | | int | 5 | Maximum concurrent operations |
| `--quality` | | string | standard | Quality preset (low/standard/high) |
| `--preserve-formatting` | | flag | true | Preserve original formatting |
| `--progress` | | flag | false | Show progress bar |

**Examples:**
```bash
# Simple text translation
translator translate "Hello world" --target es

# File translation with specific provider
translator translate book.epub --target fr --provider anthropic --model claude-3-sonnet

# Directory translation
translator translate ./books/ --target de --recursive --parallel --max-concurrency 10

# With custom output
translator translate input.txt --target ja --output translated_ja.txt

# High quality translation
translator translate novel.epub --target ru --quality high --preserve-formatting
```

#### `config`
Manage configuration settings.

**Subcommands:**
- `config --init` - Create default configuration
- `config --show` - Display current configuration
- `config --set key=value` - Set configuration value
- `config --get key` - Get configuration value
- `config --validate` - Validate configuration file

**Examples:**
```bash
# Initialize configuration
translator config --init

# Show current config
translator config --show

# Set default provider
translator config --set translation.default_provider=anthropic

# Get specific setting
translator config --get translation.default_provider

# Validate configuration
translator config --validate
```

#### `languages`
List supported languages.

**Syntax:**
```bash
translator languages [options]
```

**Options:**
| Option | Short | Description |
|--------|--------|-------------|
| `--search` | `-s` | Search languages by name |
| `--code` | `-c` | Filter by language code |
| `--format` | `-f` | Output format (table/json) |

**Examples:**
```bash
# List all languages
translator languages

# Search for languages
translator languages --search "Chinese"

# Get specific language info
translator languages --code es

# JSON output
translator languages --format json
```

#### `providers`
List available translation providers.

**Syntax:**
```bash
translator providers [options]
```

**Options:**
| Option | Short | Description |
|--------|--------|-------------|
| `--detailed` | `-d` | Show detailed provider information |
| `--check` | | Check provider availability |

**Examples:**
```bash
# List all providers
translator providers

# Detailed information
translator providers --detailed

# Check availability
translator providers --check
```

#### `version`
Show version information.

**Syntax:**
```bash
translator version [options]
```

**Options:**
| Option | Short | Description |
|--------|--------|-------------|
| `--detailed` | `-d` | Show detailed version information |
| `--check-update` | | Check for updates |

**Examples:**
```bash
# Basic version
translator version

# Detailed version
translator version --detailed

# Check for updates
translator version --check-update
```

---

## Configuration

### Configuration File Location
- **Linux/macOS**: `~/.translator/config.json`
- **Windows**: `%APPDATA%\Translator\config.json`
- **Custom**: Use `--config` flag

### Configuration Structure
```json
{
  "translation": {
    "default_provider": "openai",
    "default_model": "",
    "cache_enabled": true,
    "cache_ttl": 3600,
    "max_concurrent": 5,
    "quality_preset": "standard",
    "preserve_formatting": true
  },
  "providers": {
    "openai": {
      "api_key": "",
      "base_url": "",
      "model": "gpt-4",
      "max_tokens": 4096
    },
    "anthropic": {
      "api_key": "",
      "base_url": "",
      "model": "claude-3-sonnet",
      "max_tokens": 4096
    },
    "deepseek": {
      "api_key": "",
      "base_url": "",
      "model": "deepseek-chat",
      "max_tokens": 4096
    }
  },
  "languages": {
    "default_source": "auto",
    "default_target": "en",
    "auto_detect": true
  },
  "output": {
    "default_format": "auto",
    "preserve_structure": true,
    "create_backup": true,
    "compression": "none"
  },
  "logging": {
    "level": "info",
    "file": "",
    "format": "text"
  }
}
```

### Environment Variables
```bash
# API Keys (alternative to config file)
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export DEEPSEEK_API_KEY="sk-..."

# Configuration overrides
export TRANSLATOR_CONFIG_PATH="/path/to/config.json"
export TRANSLATOR_CACHE_DIR="/path/to/cache"
export TRANSLATOR_LOG_LEVEL="debug"
```

### Provider Configuration

#### OpenAI
```json
{
  "providers": {
    "openai": {
      "api_key": "sk-...",
      "base_url": "https://api.openai.com/v1",
      "model": "gpt-4",
      "max_tokens": 4096,
      "temperature": 0.3,
      "timeout": 30
    }
  }
}
```

#### Anthropic
```json
{
  "providers": {
    "anthropic": {
      "api_key": "sk-ant-...",
      "base_url": "https://api.anthropic.com",
      "model": "claude-3-sonnet",
      "max_tokens": 4096,
      "timeout": 30
    }
  }
}
```

#### DeepSeek
```json
{
  "providers": {
    "deepseek": {
      "api_key": "sk-...",
      "base_url": "https://api.deepseek.com",
      "model": "deepseek-chat",
      "max_tokens": 4096,
      "timeout": 30
    }
  }
}
```

---

## Advanced Usage

### Batch Processing
```bash
# Process multiple files
translator translate *.epub --target es --parallel --max-concurrency 8

# Directory with specific pattern
translator translate ./books/*.fb2 --target fr --recursive

# Custom output naming
translator translate ./input/ --target de --output "./translated/{original_name}_{target_lang}.{ext}"
```

### Quality Presets
```bash
# Low quality (fast)
translator translate book.epub --target es --quality low

# Standard quality (balanced)
translator translate book.epub --target es --quality standard

# High quality (slow)
translator translate book.epub --target es --quality high
```

### Script Conversion
```bash
# Convert to Cyrillic script
translator translate book.txt --target sr --script cyrillic

# Convert to Latin script
translator translate book.txt --target ru --script latin

# Auto-detect script
translator translate book.txt --target sr --script auto
```

### Context-Aware Translation
```bash
# Provide context for better translation
translator translate "The bill is due" --target es --context "Financial document"

# Technical context
translator translate "API endpoint" --target ja --context "Software development"
```

### Progress Monitoring
```bash
# Show progress bar
translator translate large_book.epub --target fr --progress

# Verbose output
translator translate book.epub --target de --verbose

# Quiet mode (errors only)
translator translate book.epub --target it --quiet
```

---

## Examples

### Example 1: Academic Paper Translation
```bash
# Translate academic paper with high quality
translator translate research_paper.pdf --target en \
  --provider anthropic \
  --model claude-3-opus \
  --quality high \
  --context "Academic research paper" \
  --preserve-formatting
```

### Example 2: Novel Translation Series
```bash
# Translate entire novel series
translator translate ./fantasy_series/ --target es \
  --recursive \
  --parallel \
  --max-concurrency 4 \
  --quality standard \
  --progress

# Process with custom naming
translator translate ./fantasy_series/book1.epub --target es \
  --output "./translated/book1_es.epub"
translator translate ./fantasy_series/book2.epub --target es \
  --output "./translated/book2_es.epub"
```

### Example 3: Technical Documentation
```bash
# Translate technical documentation
translator translate ./docs/ --target ja \
  --recursive \
  --provider openai \
  --model gpt-4 \
  --context "Technical documentation" \
  --preserve-formatting
```

### Example 4: Quick Translation
```bash
# Quick single sentence
translator translate "Where is the nearest restaurant?" --target es

# Multiple phrases
translator translate "Hello world" "How are you?" "Thank you" --target fr
```

### Example 5: Batch Processing with Custom Settings
```bash
# Process with configuration file
translator translate ./books/ --target de \
  --config ./production_config.json \
  --recursive \
  --parallel \
  --max-concurrency 10 \
  --quality high
```

---

## Troubleshooting

### Common Issues

#### Authentication Errors
```
Error: API key not found or invalid
```

**Solutions:**
1. Check API key in configuration file
2. Set environment variable: `export OPENAI_API_KEY="sk-..."`
3. Use `--config` flag to specify correct config file
4. Verify API key is valid and active

#### Network Issues
```
Error: Connection timeout
Error: Unable to reach provider
```

**Solutions:**
1. Check internet connection
2. Verify firewall settings
3. Try different provider
4. Set custom base URL if using proxy

#### File Format Issues
```
Error: Unsupported file format
Error: Unable to parse file
```

**Solutions:**
1. Check supported formats: EPUB, FB2, TXT, MD
2. Verify file is not corrupted
3. Try converting to supported format first
4. Use `--format` flag to specify output format

#### Memory Issues
```
Error: Out of memory
Error: File too large
```

**Solutions:**
1. Reduce `--max-concurrency`
2. Use `--quality low` for faster processing
3. Split large files into smaller chunks
4. Increase system RAM or use machine with more memory

#### Rate Limiting
```
Error: Rate limit exceeded
Error: Too many requests
```

**Solutions:**
1. Wait and retry later
2. Use different provider
3. Reduce `--max-concurrency`
4. Upgrade API plan for higher limits

### Debug Mode
```bash
# Enable verbose logging
translator translate book.epub --target es --verbose

# Debug configuration
translator config --show --verbose

# Test provider connectivity
translator providers --check --verbose
```

### Log Files
```bash
# Default log location
~/.translator/logs/
# Custom log location
export TRANSLATOR_LOG_DIR="/path/to/logs"

# Log levels
export TRANSLATOR_LOG_LEVEL="debug"  # debug, info, warn, error
```

---

## Tips & Best Practices

### Performance Optimization
1. **Use Parallel Processing**: Enable `--parallel` for multiple files
2. **Adjust Concurrency**: Set `--max-concurrency` based on your system
3. **Choose Right Quality**: Use `standard` quality for balance
4. **Cache Results**: Enable caching for repeated translations
5. **Batch Operations**: Process multiple files at once

### Quality Improvement
1. **Provide Context**: Use `--context` for domain-specific translations
2. **Choose Appropriate Model**: Use higher-end models for complex content
3. **Review Output**: Always review important translations
4. **Iterative Refinement**: Translate in stages for long documents
5. **Script Conversion**: Use correct script for target language

### Cost Management
1. **Monitor Usage**: Track token usage and costs
2. **Choose Efficient Models**: Balance quality vs cost
3. **Use Caching**: Avoid re-translating same content
4. **Batch Processing**: More efficient than individual requests
5. **Provider Selection**: Compare costs across providers

### Security Best Practices
1. **Protect API Keys**: Never share or commit API keys
2. **Use Environment Variables**: Store sensitive data in environment
3. **Regular Rotation**: Rotate API keys periodically
4. **Access Control**: Limit API key permissions
5. **Audit Logs**: Monitor usage for unauthorized access

### Workflow Automation
```bash
# Create translation script
#!/bin/bash
# translate_books.sh

SOURCE_DIR="./books/english"
TARGET_DIR="./books/spanish"
TARGET_LANG="es"

mkdir -p "$TARGET_DIR"

for file in "$SOURCE_DIR"/*.epub; do
    filename=$(basename "$file")
    translator translate "$file" \
        --target "$TARGET_LANG" \
        --output "$TARGET_DIR/${filename%.epub}_es.epub" \
        --progress \
        --quality standard
done

echo "Translation completed!"
```

### Integration with Other Tools
```bash
# Use with find command
find ./docs/ -name "*.md" -exec translator translate {} --target es \;

# Use with xargs
ls *.txt | xargs -I {} translator translate {} --target fr

# Use in makefile
translate-all:
    translator translate ./books/ --target de --recursive --parallel

clean:
    rm -rf ./translated/
```

---

## Keyboard Shortcuts (Interactive Mode)

When using the interactive mode (`translator translate --interactive`):

| Shortcut | Action |
|----------|--------|
| `Ctrl+C` | Cancel current operation |
| `Ctrl+Z` | Pause operation |
| `Space` | Show/hide progress |
| `V` | Toggle verbose mode |
| `Q` | Quit interactive mode |
| `H` | Show help |
| `S` | Save current session |

---

## Support & Resources

### Getting Help
```bash
# General help
translator --help

# Command-specific help
translator translate --help
translator config --help

# Version information
translator version --detailed
```

### Community Resources
- **Documentation**: https://docs.translator.example.com
- **GitHub**: https://github.com/translator/universal
- **Discord**: https://discord.gg/translator
- **Stack Overflow**: https://stackoverflow.com/questions/tagged/universal-translator

### Reporting Issues
1. Check existing issues on GitHub
2. Create new issue with:
   - CLI version (`translator version --detailed`)
   - Operating system
   - Error message
   - Steps to reproduce
   - Configuration file (sanitized)

### Contributing
1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Submit pull request
5. Follow contribution guidelines

---

*Last Updated: January 2024*
*Version: 1.0.0*