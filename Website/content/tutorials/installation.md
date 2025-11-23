---
title: "Installation and Setup Tutorial"
date: "2024-01-15"
weight: 10
---

# Installation and Setup Tutorial

This comprehensive tutorial will guide you through installing and setting up the Universal Multi-Format Multi-Language Ebook Translation System on your system.

## Prerequisites

### System Requirements

Before installing, ensure your system meets these requirements:

- **Operating System**: Windows 10+, macOS 10.15+, or Linux (Ubuntu 18.04+)
- **Memory**: Minimum 4GB RAM (8GB recommended for large files)
- **Storage**: 1GB free disk space
- **Network**: Internet connection for cloud providers

### Required Software

- **Go**: Version 1.21 or higher (if installing from source)
- **Docker**: Optional, for containerized installation
- **Git**: For cloning the repository (if needed)

## Installation Methods

### Method 1: Binary Installation (Recommended)

The easiest way to install is to download pre-compiled binaries for your platform.

#### Step 1: Download the Binary

**For Linux (AMD64):**
```bash
curl -L -o translator.tar.gz https://github.com/digital-vasic/translator/releases/latest/download/translator-linux-amd64.tar.gz
tar -xzf translator.tar.gz
cd translator-linux-amd64
```

**For macOS (Intel):**
```bash
curl -L -o translator.tar.gz https://github.com/digital-vasic/translator/releases/latest/download/translator-darwin-amd64.tar.gz
tar -xzf translator.tar.gz
cd translator-darwin-amd64
```

**For macOS (Apple Silicon):**
```bash
curl -L -o translator.tar.gz https://github.com/digital-vasic/translator/releases/latest/download/translator-darwin-arm64.tar.gz
tar -xzf translator.tar.gz
cd translator-darwin-arm64
```

**For Windows:**
```powershell
Invoke-WebRequest -Uri "https://github.com/digital-vasic/translator/releases/latest/download/translator-windows-amd64.zip" -OutFile "translator.zip"
Expand-Archive -Path translator.zip -DestinationPath .
cd translator-windows-amd64
```

#### Step 2: Add to PATH

**Linux/macOS:**
```bash
# Add to shell profile
echo 'export PATH=$PATH:'$(pwd) >> ~/.bashrc
# or for Zsh
echo 'export PATH=$PATH:'$(pwd) >> ~/.zshrc

# Reload shell
source ~/.bashrc  # or ~/.zshrc
```

**Windows (Command Prompt):**
```cmd
set PATH=%PATH%;%CD%
```

**Windows (PowerShell):**
```powershell
$env:PATH += ";$PWD"
```

#### Step 3: Verify Installation

```bash
translator version
```

You should see version information printed.

### Method 2: Go Installation

If you have Go installed, you can install directly from the source:

```bash
go install github.com/digital-vasic/translator/cmd/cli@latest
```

The binary will be installed to `$GOPATH/bin` (usually `~/go/bin`). Add this to your PATH if needed.

### Method 3: Docker Installation

For containerized deployment:

#### Option A: Run Directly

```bash
docker run -p 8080:8080 \
  -e OPENAI_API_KEY="your-api-key" \
  -v $(pwd)/books:/books \
  digitalvasic/translator:latest
```

#### Option B: Using Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  translator:
    image: digitalvasic/translator:latest
    ports:
      - "8080:8080"
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
    volumes:
      - ./books:/books
      - ./config:/config
    restart: unless-stopped
```

Then run:

```bash
docker-compose up -d
```

### Method 4: From Source

For development or if you want the latest features:

```bash
git clone https://github.com/digital-vasic/translator.git
cd translator
go mod download
make build
sudo make install  # optional, installs system-wide
```

## Configuration

### Setting Up API Keys

Most translation providers require API keys. Set them as environment variables:

#### OpenAI (GPT)
```bash
export OPENAI_API_KEY="sk-your-openai-api-key"
```

#### Anthropic (Claude)
```bash
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-api-key"
```

#### Zhipu AI (GLM-4)
```bash
export ZHIPU_API_KEY="your-zhipu-api-key"
```

#### DeepSeek
```bash
export DEEPSEEK_API_KEY="your-deepseek-api-key"
```

#### Google AI (Gemini)
```bash
export GEMINI_API_KEY="your-gemini-api-key"
```

#### Qwen
```bash
export QWEN_API_KEY="your-qwen-api-key"
```

### Creating a Configuration File

Create a configuration file at `~/.translator/config.yaml`:

```yaml
# Server configuration
server:
  port: 8080
  host: "localhost"

# Translation settings
translation:
  default_provider: "openai"
  default_model: "gpt-4"
  temperature: 0.7
  max_tokens: 2000
  timeout: 30

# Quality settings
quality:
  minimum_score: 0.8
  enable_verification: true
  check_grammar: true
  check_style: true

# Serbian language settings
serbian:
  default_script: "cyrillic"  # or "latin"
  phonetic_transliteration: true

# Cache settings
cache:
  enabled: true
  ttl: 86400  # 24 hours
  max_size: 1000

# Provider configurations
providers:
  openai:
    api_key: "${OPENAI_API_KEY}"
    model: "gpt-4"
    base_url: "https://api.openai.com/v1"
  
  anthropic:
    api_key: "${ANTHROPIC_API_KEY}"
    model: "claude-3-sonnet"
    base_url: "https://api.anthropic.com"
  
  deepseek:
    api_key: "${DEEPSEEK_API_KEY}"
    model: "deepseek-chat"
    base_url: "https://api.deepseek.com"
```

### Setting Up Local Providers

For offline translation with Ollama:

```bash
# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama3:8b

# Configure translator
export OLLAMA_BASE_URL="http://localhost:11434"
```

For LlamaCpp:

```bash
# Download a model
wget https://huggingface.co/TheBloke/Mistral-7B-Instruct-v0.2-GGUF/resolve/main/mistral-7b-instruct-v0.2.Q4_K_M.gguf

# Set model path
export LLAMACPP_MODEL_PATH="/path/to/your/model.gguf"
```

## Testing Your Installation

### 1. Test Basic Functionality

```bash
translator translate --text "Hello, world!" --from en --to sr
```

### 2. Test File Translation

Create a test file:

```bash
echo "Hello, world! This is a test." > test.txt
translator translate test.txt --from en --to sr
```

### 3. Test Web Interface

Start the web server:

```bash
translator server --port 8080
```

Open your browser to `http://localhost:8080` and test the interface.

### 4. Test Configuration

Verify your configuration:

```bash
translator config show
translator config test
```

## Troubleshooting Installation Issues

### Common Problems

#### Permission Denied (Linux/macOS)

```bash
chmod +x translator
# or install to user directory
mkdir -p ~/bin
cp translator ~/bin/
export PATH=$HOME/bin:$PATH
```

#### API Key Not Found

```bash
# Check if environment variable is set
echo $OPENAI_API_KEY

# Set it permanently
echo 'export OPENAI_API_KEY="your-key"' >> ~/.bashrc
source ~/.bashrc
```

#### Port Already in Use

```bash
# Check what's using the port
lsof -i :8080  # Linux/macOS
netstat -ano | findstr :8080  # Windows

# Use a different port
translator server --port 8081
```

#### Memory Issues

For large files, adjust memory limits:

```bash
# Increase Go memory limit
export GOMEMLIMIT=4GiB
translator translate large_book.fb2 --from ru --to sr
```

### Verification Commands

Run these commands to verify your installation:

```bash
# Check version
translator version

# Check available providers
translator providers list

# Check API keys
translator config test

# Test translation with different providers
translator translate --text "test" --from en --to sr --provider openai
translator translate --text "test" --from en --to sr --provider deepseek
```

## Next Steps

Congratulations! You now have the translator system installed and configured. Here are the recommended next steps:

1. [Your First Translation Tutorial](/tutorials/first-translation) - Learn the basics
2. [Web Interface Tutorial](/tutorials/web-interface) - Master the GUI
3. [File Formats Tutorial](/tutorials/file-formats) - Understand format support
4. [Provider Selection Tutorial](/tutorials/provider-selection) - Choose the right provider

## Getting Help

If you encounter issues during installation:

- **Documentation**: [Full User Manual](/docs/user-manual)
- **GitHub Issues**: [Report problems](https://github.com/digital-vasic/translator/issues)
- **Community**: [Join Discord](https://discord.gg/translator)
- **Email**: [support@translator.digital](mailto:support@translator.digital)

## Video Tutorial

For a visual guide, watch our installation video:

[![Installation Video](https://img.youtube.com/vi/YOUTUBE_VIDEO_ID/maxresdefault.jpg)](https://www.youtube.com/watch?v=YOUTUBE_VIDEO_ID)

---

Now that you have the system installed, let's move on to [Your First Translation](/tutorials/first-translation)!