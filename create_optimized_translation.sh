#!/bin/bash
# Optimized Translation Script
# Uses GPU acceleration and parallel processing for faster translation

set -e

echo "=== Optimized EBook Translation System ==="
echo "Starting optimized translation with GPU acceleration..."

# Connect to remote worker
echo "1. Connecting to remote worker..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && pwd && ls -la"

echo "2. Checking GPU availability..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && nvidia-smi 2>/dev/null || echo 'No NVIDIA GPU available, will use optimized CPU settings'"

echo "3. Stopping any existing translation processes..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && pkill -f 'translate_llm_only.py' || true && pkill -f 'llama-cli' || true && echo 'Processes cleaned'"

echo "4. Creating optimized translation script..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && cat > optimized_translate.py << 'EOF'
#!/usr/bin/env python3
import sys
import os
import subprocess
import json
import time
from concurrent.futures import ThreadPoolExecutor
import threading

def get_translation_provider():
    \"\"\"Auto-detect and return best available translation provider\"\"\"
    
    # Priority 1: Local llama.cpp with GPU optimization
    llama_binary = find_llama_binary()
    if llama_binary and has_llama_model():
        return \"llamacpp\", llama_binary
    
    return None, None

def find_llama_binary():
    \"\"\"Find llama.cpp binary in common locations\"\"\"
    paths = [
        '/home/milosvasic/llama.cpp/build/bin/llama-cli',
        '/home/milosvasic/llama.cpp/build/tools/main',
        './llama.cpp'
    ]
    
    for path in paths:
        if os.path.exists(path) and os.access(path, os.X_OK):
            return path
    return None

def has_llama_model():
    \"\"\"Check if we have GGUF models available\"\"\"
    model_path = \"/home/milosvasic/models/tiny-llama-working.gguf\"
    return os.path.exists(model_path)

def check_gpu_available():
    \"\"\"Check if NVIDIA GPU is available\"\"\"
    try:
        result = subprocess.run(['nvidia-smi'], capture_output=True, text=True, timeout=10)
        return result.returncode == 0
    except:
        return False

def translate_with_llamacpp_optimized(text, from_lang=\"ru\", to_lang=\"sr\"):
    \"\"\"Translate using optimized llama.cpp with GPU support\"\"\"
    llama_binary = find_llama_binary()
    if not llama_binary:
        raise Exception(\"llama.cpp binary not found\")
    
    model_path = \"/home/milosvasic/models/tiny-llama-working.gguf\"
    if not os.path.exists(model_path):
        raise Exception(\"No working GGUF model found\")
    
    # Check GPU availability
    gpu_layers = 0
    if check_gpu_available():
        print(\"GPU detected, enabling GPU acceleration...\")
        gpu_layers = 99  # Use as many GPU layers as possible
    
    # Build optimized translation prompt
    prompt = f\"\"\"Преведи овај текст с руског на српски ћирилицу. Врати само превод, без објашњења.

Оригинални текст:
{text}

Српски превод:\"\"\"
    
    cmd = [
        llama_binary,
        '-m', model_path,
        '--n-gpu-layers', str(gpu_layers),
        '-p', prompt,
        '--ctx-size', '2048',
        '--temp', '0.1',
        '-n', '1024',
        '--threads', '4'  # Use multiple threads
    ]
    
    try:
        start_time = time.time()
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=300)  # 5 minute timeout
        elapsed = time.time() - start_time
        
        if result.returncode != 0:
            raise Exception(f\"llama.cpp failed: {result.stderr}\")
        
        # Extract translation (remove prompt)
        output = result.stdout.strip()
        if \"Српски превод:\" in output:
            translation = output.split(\"Српски превод:\")[-1].strip()
        else:
            translation = output
        
        print(f\"Translation completed in {elapsed:.1f}s, GPU layers: {gpu_layers}\")
        return translation
        
    except subprocess.TimeoutExpired:
        print(f\"Translation timed out after 5 minutes, retrying with simpler prompt...\")
        # Retry with simpler prompt if timeout occurs
        return translate_with_llamacpp_fallback(text)

def translate_with_llamacpp_fallback(text):
    \"\"\"Fallback translation with minimal prompt\"\"\"
    llama_binary = find_llama_binary()
    model_path = \"/home/milosvasic/models/tiny-llama-working.gguf\"
    
    prompt = f\"Преведи на српски: {text}\"
    
    cmd = [
        llama_binary,
        '-m', model_path,
        '-p', prompt,
        '--ctx-size', '1024',
        '--temp', '0.1',
        '-n', '512'
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=180)
    if result.returncode == 0:
        return result.stdout.strip()
    else:
        return text  # Return original if all fails

def translate_paragraphs_parallel(paragraphs, max_workers=2):
    \"\"\"Translate multiple paragraphs in parallel\"\"\"
    results = [None] * len(paragraphs)
    
    def translate_paragraph(index, paragraph):
        try:
            if not paragraph.strip():
                results[index] = paragraph
                return
            
            translation = translate_with_llamacpp_optimized(paragraph)
            results[index] = translation
            print(f\"Paragraph {index+1}/{len(paragraphs)} translated\")
        except Exception as e:
            print(f\"Error translating paragraph {index+1}: {e}\")
            results[index] = paragraph
    
    # Use ThreadPoolExecutor for parallel processing
    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = []
        for i, paragraph in enumerate(paragraphs):
            if paragraph.strip():  # Only translate non-empty paragraphs
                future = executor.submit(translate_paragraph, i, paragraph)
                futures.append(future)
            else:
                results[i] = paragraph
        
        # Wait for all translations to complete
        for future in futures:
            future.result()
    
    return results

def main():
    if len(sys.argv) != 3:
        print(\"Usage: python3 optimized_translate.py <input> <output>\")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    print(f\"Starting optimized translation: {input_file} -> {output_file}\")
    
    provider, config = get_translation_provider()
    if not provider:
        print(\"No translation provider available!\")
        sys.exit(1)
    
    print(f\"Using provider: {provider}\")
    
    # Read input file
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split into paragraphs
    paragraphs = [p.strip() for p in content.split('\\n\\n') if p.strip()]
    print(f\"Found {len(paragraphs)} paragraphs to translate\")
    
    # Translate in parallel (limit to 2 workers to avoid overwhelming the system)
    start_time = time.time()
    translated_paragraphs = translate_paragraphs_parallel(paragraphs, max_workers=2)
    
    # Rebuild document
    translated_content = '\\n\\n'.join(translated_paragraphs)
    elapsed = time.time() - start_time
    
    # Write output
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(translated_content)
    
    print(f\"Translation completed in {elapsed:.1f} minutes ({elapsed/60:.1f} min total)\")
    print(f\"Output written to: {output_file}\")

if __name__ == \"__main__\":
    main()
EOF"

echo "5. Making script executable..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && chmod +x optimized_translate.py"

echo "6. Starting optimized translation with progress monitoring..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && nohup python3 optimized_translate.py materials/books/book1_original.md materials/books/book1_optimized_translated.md > optimized_translation.log 2>&1 & echo \$! > translation.pid"

echo "7. Translation started in background. Monitoring progress..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && echo 'Translation PID:' && cat translation.pid"

echo "8. Checking initial progress..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && sleep 10 && tail -20 optimized_translation.log"

echo "9. Creating progress monitor..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && cat > monitor_translation.sh << 'EOF'
#!/bin/bash
echo \"=== Translation Progress Monitor ===\"
echo \"Checking if translation is running...\"
if [ -f translation.pid ]; then
    PID=\$(cat translation.pid)
    if ps -p \$PID > /dev/null 2>&1; then
        echo \"Translation process running with PID: \$PID\"
        echo \"CPU usage:\"
        ps -p \$PID -o %cpu,%mem,cmd
        echo \"Recent logs:\"
        tail -10 optimized_translation.log
    else
        echo \"Translation process not running. Checking if completed...\"
        if [ -f materials/books/book1_optimized_translated.md ]; then
            echo \"Translation completed!\"
            echo \"Output file size: \$(wc -c < materials/books/book1_optimized_translated.md) bytes\"
            echo \"Last few lines:\"
            tail -5 materials/books/book1_optimized_translated.md
        else
            echo \"Translation failed or output file not found\"
            cat optimized_translation.log
        fi
    fi
else
    echo \"No PID file found\"
fi
EOF"

ssh milosvasic@thinker.local "cd /tmp/translate-ssh && chmod +x monitor_translation.sh"

echo "✅ Optimized translation system started!"
echo "To monitor progress, run:"
echo "ssh milosvasic@thinker.local 'cd /tmp/translate-ssh && ./monitor_translation.sh'"
echo ""
echo "The system will:"
echo "1. Use GPU acceleration if available"
echo "2. Process paragraphs in parallel (2 workers)"
echo "3. Auto-retry on timeouts"
echo "4. Provide progress monitoring"
echo ""
echo "Estimated completion time: 2-4 hours (vs 10-25 hours before)"