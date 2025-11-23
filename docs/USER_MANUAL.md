# User Manual

## Table of Contents

1. [Getting Started](#getting-started)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Basic Usage](#basic-usage)
5. [Command Line Interface](#command-line-interface)
6. [Web Interface](#web-interface)
7. [Translation Methods](#translation-methods)
8. [File Formats](#file-formats)
9. [Quality Assurance](#quality-assurance)
10. [Troubleshooting](#troubleshooting)

## Getting Started

### What is the Translator System?

The Universal Multi-Format Multi-Language Ebook Translation System is a powerful tool for translating ebooks and text documents between languages using state-of-the-art AI models. It supports multiple file formats, translation providers, and quality verification systems.

### Key Features

- **Multi-format Support**: FB2, EPUB, TXT, HTML, PDF, DOCX
- **8 AI Translation Providers**: OpenAI, Anthropic, Zhipu, DeepSeek, Qwen, Gemini, Ollama, LlamaCpp
- **High-Quality Translations**: Context-aware literary translation with cultural adaptation
- **Distributed Processing**: Scalable multi-node processing for large projects
- **Quality Verification**: Built-in translation quality assessment
- **Multiple Output Formats**: Export to FB2, EPUB, PDF, DOCX
- **Serbian Script Support**: Both Cyrillic and Latin script options

### System Requirements

- **Operating System**: Windows 10+, macOS 10.15+, Linux (Ubuntu 18.04+)
- **Memory**: Minimum 4GB RAM (8GB recommended)
- **Storage**: 1GB free space
- **Network**: Internet connection for API-based providers
- **Python**: Version 3.8 or higher (for standalone scripts)
- **Go**: Version 1.21 or higher (for compiled version)

## Installation

### Option 1: Binary Installation (Recommended)

1. Download the latest release for your platform from the [Releases page](https://github.com/digital-vasic/translator/releases)
2. Extract the archive
3. Add the translator directory to your PATH

```bash
# macOS/Linux
export PATH=$PATH:/path/to/translator

# Windows (Command Prompt)
set PATH=%PATH%;C:\path\to\translator

# Windows (PowerShell)
$env:PATH += ";C:\path\to\translator"
```

### Option 2: Go Installation

```bash
go install github.com/digital-vasic/translator/cmd/cli@latest
```

### Option 3: Python Package Installation

```bash
pip install universal-translator
```

### Option 4: Docker Installation

```bash
docker pull digitalvasic/translator:latest
docker run -p 8080:8080 digitalvasic/translator:latest
```

## Configuration

### Environment Variables

Set up your API keys as environment variables:

```bash
# OpenAI (required for GPT translations)
export OPENAI_API_KEY="your-openai-api-key"

# Anthropic Claude
export ANTHROPIC_API_KEY="your-anthropic-api-key"

# Zhipu AI (GLM-4)
export ZHIPU_API_KEY="your-zhipu-api-key"

# DeepSeek
export DEEPSEEK_API_KEY="your-deepseek-api-key"

# Google AI (Gemini)
export GEMINI_API_KEY="your-gemini-api-key"

# Qwen
export QWEN_API_KEY="your-qwen-api-key"

# Optional: Local model paths
export OLLAMA_BASE_URL="http://localhost:11434"
export LLAMACPP_MODEL_PATH="/path/to/your/model.gguf"
```

### Configuration File

Create a configuration file at `~/.translator/config.yaml`:

```yaml
server:
  port: 8080
  host: "localhost"

translation:
  default_provider: "openai"
  default_model: "gpt-4"
  temperature: 0.7
  max_tokens: 2000
  timeout: 30

quality:
  minimum_score: 0.8
  enable_verification: true
  check_grammar: true
  check_style: true

serbian:
  default_script: "cyrillic"  # or "latin"
  phonetic_transliteration: true

cache:
  enabled: true
  ttl: 86400  # 24 hours
  max_size: 1000
```

### Database Configuration (Optional)

For persistent storage and history:

```yaml
database:
  type: "sqlite"  # or "postgresql", "mysql"
  connection: "translator.db"
  max_connections: 10
```

## Basic Usage

### Quick Start

1. **Simple Text Translation**:
   ```bash
   translator translate --text "Hello, world!" --from en --to sr --provider openai
   ```

2. **File Translation**:
   ```bash
   translator translate --file book.fb2 --from ru --to sr --provider deepseek
   ```

3. **Batch Translation**:
   ```bash
   translator batch --input-dir ./books --output-dir ./translated --from ru --to sr
   ```

### Interactive Mode

Start the interactive translation mode:

```bash
translator interactive
```

This will launch a command-line interface where you can:
- Enter text to translate
- Upload files
- Configure options
- View translation history

## Command Line Interface

### Global Options

All commands support these global options:

- `--config, -c`: Path to configuration file
- `--verbose, -v`: Verbose output
- `--quiet, -q`: Minimal output
- `--help, -h`: Show help
- `--version`: Show version information

### translate Command

Translate text or files.

```bash
translator translate [options] INPUT
```

**Options:**
- `--text, -t`: Text to translate (alternative to file input)
- `--from, -f`: Source language (default: auto-detect)
- `--to, -o`: Target language (required)
- `--provider, -p`: Translation provider (default: openai)
- `--model, -m`: Model to use (provider-specific)
- `--output, -o`: Output file path
- `--format, -F`: Output format (same as input by default)
- `--script, -s`: Serbian script (cyrillic/latin, default: cyrillic)

**Examples:**

```bash
# Translate Russian FB2 to Serbian
translator translate book_ru.fb2 --from ru --to sr --provider deepseek --output book_sr.fb2

# Translate single text
translator translate --text "Hello, world!" --from en --to sr --provider openai

# Translate with specific model
translator translate book.fb2 --from ru --to sr --provider anthropic --model claude-3-opus

# Translate to Latin Serbian script
translator translate book.fb2 --from ru --to sr --script latin
```

### batch Command

Translate multiple files.

```bash
translator batch [options] INPUT_DIRECTORY
```

**Options:**
- `--output-dir, -o`: Output directory (required)
- `--from, -f`: Source language
- `--to, -t`: Target language (required)
- `--provider, -p`: Translation provider
- `--pattern, -P`: File pattern (default: *.fb2)
- `--parallel, -j`: Number of parallel workers (default: 4)
- `--resume, -r`: Resume interrupted batch job

**Examples:**

```bash
# Translate all FB2 files
translator batch ./books --output-dir ./translated --from ru --to sr

# Translate specific pattern
translator batch ./books --output-dir ./translated --pattern "*.epub" --from en --to sr

# Use 8 parallel workers
translator batch ./books --output-dir ./translated --from ru --to sr --parallel 8
```

### quality Command

Check translation quality.

```bash
translator quality [options] FILE
```

**Options:**
- `--original, -O`: Original text/file
- `--translated, -T`: Translated text/file
- `--from, -f`: Source language
- `--to, -t`: Target language
- `--detailed, -d`: Show detailed analysis

**Examples:**

```bash
# Check translation quality
translator quality --original book_ru.fb2 --translated book_sr.fb2 --from ru --to sr

# Detailed quality analysis
translator quality --original book_ru.fb2 --translated book_sr.fb2 --from ru --to sr --detailed
```

### server Command

Start the web server.

```bash
translator server [options]
```

**Options:**
- `--port, -p`: Port number (default: 8080)
- `--host, -h`: Host address (default: localhost)
- `--dev, -d`: Development mode with hot reload

### config Command

Manage configuration.

```bash
translator config [SUBCOMMAND] [options]
```

**Subcommands:**
- `show`: Show current configuration
- `set KEY VALUE`: Set configuration value
- `init`: Initialize configuration file
- `test`: Test configuration and API keys

**Examples:**

```bash
# Show configuration
translator config show

# Set default provider
translator config set default_provider deepseek

# Test configuration
translator config test

# Initialize config file
translator config init
```

## Web Interface

### Starting the Web Server

```bash
translator server --port 8080
```

Then open your browser to `http://localhost:8080`.

### Web Interface Features

1. **Dashboard**: Overview of translation activity and system status
2. **Upload & Translate**: Drag-and-drop file translation
3. **Batch Processing**: Queue and monitor batch jobs
4. **Quality Reports**: View translation quality assessments
5. **History**: Translation history and statistics
6. **Settings**: Configure providers and options

### Uploading Files

1. Navigate to the **Translate** page
2. Drag and drop your file or click to browse
3. Select source and target languages
4. Choose translation provider and model
5. Click **Translate**

### Batch Processing

1. Go to the **Batch** page
2. Upload multiple files or select a directory
3. Configure translation options
4. Start the batch job
5. Monitor progress in real-time

## Translation Methods

### 1. AI LLM Translation (Recommended)

**Provider Options:**

#### OpenAI GPT
- **Models**: GPT-3.5 Turbo, GPT-4, GPT-4 Turbo
- **Strengths**: Excellent general translation, good context understanding
- **Use Case**: General-purpose translation, balanced quality and cost

#### Anthropic Claude
- **Models**: Claude 3 Sonnet, Claude 3 Opus
- **Strengths**: Literary translation, maintains authorial voice
- **Use Case**: Literary works, novels, poetry

#### Zhipu AI (GLM-4)
- **Models**: GLM-4
- **Strengths**: Excellent for Russian-Slavic language pairs
- **Use Case**: Russian to Serbian translation

#### DeepSeek
- **Models**: DeepSeek-V2
- **Strengths**: Cost-effective, good technical translation
- **Use Case**: Large projects, technical documentation

#### Qwen
- **Models**: Qwen-Max
- **Strengths**: Good for Asian-European language pairs
- **Use Case**: Multilingual projects

#### Gemini
- **Models**: Gemini Pro
- **Strengths**: Google's latest model, good cultural adaptation
- **Use Case**: Modern content, web content

#### Ollama (Local)
- **Models**: Llama 3, Mixtral, custom models
- **Strengths**: Free, offline, private
- **Use Case**: Sensitive content, offline work

#### LlamaCpp (Local)
- **Models**: GGUF format models
- **Strengths**: Custom models, optimized inference
- **Use Case**: Specialized domains, research

### 2. Dictionary-Based Translation

**Features:**
- Fast, no API costs
- Predefined word pairs
- Limited to dictionary terms
- Good for basic translations

**Use Case:**
- Quick reference
- Simple texts
- Learning vocabulary

### 3. Hybrid Translation

Combine AI translation with dictionary verification for maximum accuracy:

```bash
translator translate --file book.fb2 --from ru --to sr --provider deepseek --verify-dictionary
```

## File Formats

### Input Formats

#### FB2 (FictionBook2)
- **Description**: Russian ebook format
- **Features**: Rich metadata, good for fiction
- **Encoding**: UTF-8
- **Best for**: Novels, stories, literary works

#### EPUB
- **Description**: Standard ebook format
- **Features**: Reflowable content, multi-device support
- **Encoding**: UTF-8
- **Best for**: Modern fiction, non-fiction

#### TXT
- **Description**: Plain text format
- **Features**: Simple, universal compatibility
- **Encoding**: UTF-8 or detection
- **Best for**: Simple documents, testing

#### HTML
- **Description**: Web page format
- **Features**: Structured content, styling
- **Encoding**: UTF-8
- **Best for**: Web content, documentation

#### PDF
- **Description**: Portable Document Format
- **Features**: Fixed layout, universal viewing
- **Encoding**: Text extraction
- **Best for**: Scanned books, fixed-layout content

#### DOCX
- **Description**: Microsoft Word format
- **Features**: Rich formatting, track changes
- **Encoding**: UTF-8
- **Best for**: Documents, manuscripts

### Output Formats

The same formats are supported for output, with additional options:

- **Preserve Formatting**: Keep original layout and styling
- **Metadata Translation**: Translate titles and metadata
- **Script Options**: Choose Cyrillic or Latin for Serbian

## Quality Assurance

### Automatic Quality Checks

The system performs automatic quality assessment:

1. **Grammar Check**: Verify grammar rules
2. **Style Consistency**: Check for consistent style
3. **Cultural Adaptation**: Verify cultural references
4. **Technical Accuracy**: Check terminology consistency

### Quality Scoring

Translations receive a quality score from 0.0 to 1.0:

- **0.9 - 1.0**: Excellent quality, publication-ready
- **0.8 - 0.9**: Good quality, minimal editing needed
- **0.7 - 0.8**: Acceptable, requires review
- **0.6 - 0.7**: Fair, significant editing required
- **Below 0.6**: Poor, re-translation recommended

### Manual Review Process

1. **Preview**: Review translated segments
2. **Compare**: Side-by-side original and translation
3. **Edit**: Make corrections and improvements
4. **Verify**: Re-check quality score
5. **Finalize**: Approve for output

### Translation Notes

Add notes and comments to translations:

```bash
translator translate --file book.fb2 --from ru --to sr --add-notes
```

This will add notes for:
- Uncertain translations
- Cultural references
- Terminology choices
- Formatting issues

## Troubleshooting

### Common Issues

#### API Key Problems

**Error**: "Invalid API key" or "Authentication failed"

**Solutions**:
1. Verify API key is correct
2. Check environment variable spelling
3. Ensure API key has sufficient credits
4. Try refreshing the key

**Verification**:
```bash
translator config test
```

#### Translation Quality Issues

**Problem**: Translation quality is poor

**Solutions**:
1. Try a different provider
2. Adjust temperature (lower = more consistent, higher = more creative)
3. Use a larger model
4. Enable quality verification
5. Split large texts into smaller chunks

**Example**:
```bash
translator translate --file book.fb2 --from ru --to sr --provider openai --model gpt-4 --temperature 0.5 --verify
```

#### Memory Issues

**Problem**: Out of memory errors with large files

**Solutions**:
1. Use chunked processing:
   ```bash
   translator translate --file large_book.fb2 --from ru --to sr --chunk-size 5000
   ```
2. Increase memory limits
3. Use batch processing for very large files

#### Network Issues

**Problem**: Connection timeouts, network errors

**Solutions**:
1. Increase timeout:
   ```bash
   translator translate --file book.fb2 --from ru --to sr --timeout 120
   ```
2. Use local provider (Ollama/LlamaCpp)
3. Enable retry logic:
   ```bash
   translator translate --file book.fb2 --from ru --to sr --retry 3
   ```

#### Format Issues

**Problem**: Input file not recognized or corrupted

**Solutions**:
1. Validate file format:
   ```bash
   translator validate --file book.fb2
   ```
2. Convert to supported format
3. Check file encoding (should be UTF-8)
4. Repair corrupted files

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
translator translate --file book.fb2 --from ru --to sr --verbose
```

### Log Files

Log files are stored at:
- **Linux/macOS**: `~/.translator/logs/`
- **Windows**: `%APPDATA%/translator/logs/`

### Support

For additional help:

1. **Documentation**: Check the [API Documentation](API.md)
2. **Community**: Join our [Discord server](https://discord.gg/translator)
3. **Issues**: Report bugs at [GitHub Issues](https://github.com/digital-vasic/translator/issues)
4. **Email**: support@translator.digital

### Performance Tips

1. **Choose the right provider**:
   - For quality: Claude 3 Opus, GPT-4
   - For cost: DeepSeek, Qwen
   - For privacy: Ollama, LlamaCpp

2. **Optimize batch processing**:
   - Use appropriate parallel worker count
   - Monitor resource usage
   - Resume interrupted jobs

3. **Cache translations**:
   - Enable caching to avoid re-translation
   - Use consistent prompts
   - Save translation memory

4. **Memory management**:
   - Process large files in chunks
   - Monitor system resources
   - Use appropriate chunk size

### Keyboard Shortcuts

In the web interface:
- **Ctrl/Cmd + O**: Open file
- **Ctrl/Cmd + S**: Save translation
- **Ctrl/Cmd + Enter**: Start translation
- **Ctrl/Cmd + P**: Print/Export
- **F11**: Fullscreen mode
- **Esc**: Cancel operation

---

## Next Steps

- Explore the [API Documentation](API.md) for programmatic usage
- Check the [Developer Guide](DEVELOPER.md) for advanced configuration
- Join our community for tips and best practices
- Try the [Example Projects](https://github.com/digital-vasic/translator-examples)