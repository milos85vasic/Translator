#!/bin/bash
# Quick GPU-Optimized Translation Deployment
# Uses the existing RTX 3060 GPU for faster translation

set -e

echo "=== Quick GPU Translation Deployment ==="

# First, let's check if the remote worker has GPU (we saw RTX 3060 available)
echo "1. Uploading optimized translation script..."

cat > /tmp/gpu_translate.py << 'EOF'
#!/usr/bin/env python3
import sys
import os
import subprocess
import time

def translate_with_gpu():
    """Use llama.cpp with GPU acceleration"""
    
    # Check for model
    model_path = "/home/milosvasic/models/tiny-llama-working.gguf"
    if not os.path.exists(model_path):
        print("Model not found, checking alternative paths...")
        # Find any .gguf model
        for root, dirs, files in os.walk("/home/milosvasic"):
            for file in files:
                if file.endswith(".gguf"):
                    model_path = os.path.join(root, file)
                    print(f"Found model: {model_path}")
                    break
            if os.path.exists(model_path):
                break
    
    if not os.path.exists(model_path):
        print("No GGUF model found!")
        return False
    
    llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
    
    # Test text for translation
    test_text = "Это тестовый абзац для проверки перевода."
    
    # Create optimized prompt
    prompt = f"""Translate this Russian text to Serbian Cyrillic. Return ONLY the Serbian translation.

Original: {test_text}

Serbian:"""
    
    # Use GPU acceleration (we have RTX 3060 available)
    cmd = [
        llama_binary,
        '-m', model_path,
        '--n-gpu-layers', '99',  # Use GPU for all layers
        '-p', prompt,
        '--ctx-size', '2048',
        '--temp', '0.1',
        '-n', '512',
        '--threads', '4'
    ]
    
    print("Starting GPU-accelerated translation...")
    print(f"Command: {' '.join(cmd)}")
    
    start_time = time.time()
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=120)
        elapsed = time.time() - start_time
        
        if result.returncode == 0:
            output = result.stdout.strip()
            print(f"✓ Translation completed in {elapsed:.1f} seconds")
            print(f"Output: {output}")
            return True
        else:
            print(f"✗ Translation failed: {result.stderr}")
            return False
            
    except subprocess.TimeoutExpired:
        print("✗ Translation timed out")
        return False
    except Exception as e:
        print(f"✗ Error: {e}")
        return False

def translate_book(input_file, output_file):
    """Translate entire book with GPU"""
    
    # Read the book
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split into paragraphs
    paragraphs = [p.strip() for p in content.split('\n\n') if p.strip()]
    print(f"Found {len(paragraphs)} paragraphs to translate")
    
    model_path = "/home/milosvasic/models/tiny-llama-working.gguf"
    llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
    
    translated_paragraphs = []
    total_time = 0
    
    for i, paragraph in enumerate(paragraphs):
        if not paragraph or len(paragraph) < 10:
            translated_paragraphs.append(paragraph)
            continue
        
        print(f"Translating paragraph {i+1}/{len(paragraphs)}...")
        
        prompt = f"""Translate to Serbian Cyrillic:
{paragraph}

Translation:"""
        
        cmd = [
            llama_binary,
            '-m', model_path,
            '--n-gpu-layers', '99',
            '-p', prompt,
            '--ctx-size', '2048',
            '--temp', '0.1',
            '-n', '1024'
        ]
        
        start_time = time.time()
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=180)
            elapsed = time.time() - start_time
            total_time += elapsed
            
            if result.returncode == 0:
                output = result.stdout.strip()
                if "Translation:" in output:
                    translation = output.split("Translation:")[-1].strip()
                else:
                    translation = output
                translated_paragraphs.append(translation)
                print(f"  ✓ Completed in {elapsed:.1f}s")
            else:
                print(f"  ✗ Failed: {result.stderr[:100]}")
                translated_paragraphs.append(paragraph)
                
        except subprocess.TimeoutExpired:
            print(f"  ✗ Timed out")
            translated_paragraphs.append(paragraph)
        
        # Show progress
        if i > 0 and i % 5 == 0:
            avg_time = total_time / (i + 1)
            remaining = (len(paragraphs) - i - 1) * avg_time
            print(f"  Progress: {i+1}/{len(paragraphs)}, ETA: {remaining/60:.1f} min")
    
    # Write translated book
    translated_content = '\n\n'.join(translated_paragraphs)
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(translated_content)
    
    print(f"\n✓ Book translation completed!")
    print(f"Total time: {total_time/60:.1f} minutes")
    print(f"Average per paragraph: {total_time/len(paragraphs):.1f} seconds")
    print(f"Output: {output_file}")
    
    return True

if __name__ == "__main__":
    if len(sys.argv) == 1:
        # Test mode
        print("Testing GPU translation...")
        if translate_with_gpu():
            print("✓ GPU translation test successful!")
        else:
            print("✗ GPU translation test failed!")
    elif len(sys.argv) == 3:
        # Book translation mode
        input_file = sys.argv[1]
        output_file = sys.argv[2]
        print(f"Translating book: {input_file} -> {output_file}")
        translate_book(input_file, output_file)
    else:
        print("Usage:")
        print("  python3 gpu_translate.py                    # Test GPU translation")
        print("  python3 gpu_translate.py input.md output.md # Translate book")
EOF

echo "2. Deploying to remote worker..."
scp /tmp/gpu_translate.py milosvasic@thinker.local:/tmp/translate-ssh/
ssh milosvasic@thinker.local "chmod +x /tmp/translate-ssh/gpu_translate.py"

echo "3. Testing GPU translation..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && python3 gpu_translate.py"

echo "4. Starting full book translation with GPU..."
# Use the existing markdown file on the remote worker
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && python3 gpu_translate.py book1_original.md book1_gpu_translated.md 2>&1 | tee gpu_translation.log &"

echo "5. Checking initial progress..."
sleep 10
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && echo 'Recent log:' && tail -10 gpu_translation.log"

echo "✅ GPU-optimized translation started!"
echo "To monitor: ssh milosvasic@thinker.local 'cd /tmp/translate-ssh && tail -f gpu_translation.log'"