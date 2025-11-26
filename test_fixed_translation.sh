#!/bin/bash

echo "=== Testing Translation with Fixed Parameters ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Testing without --no-cnv flag..."

cat > test_fixed.py << 'PYEOF'
#!/usr/bin/env python3
import os
import subprocess

# Test with simpler prompt format
text = "Это простая тестовая книга на русском языке."

llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
model_path = "/home/milosvasic/models/tiny-llama-working.gguf"

if not os.path.exists(model_path):
    model_path = "/home/milosvasic/models/tiny-llama.gguf"

# Simpler prompt - just give the task without markers
prompt = f"Translate this Russian text to Serbian Cyrillic: {text}"

cmd = [
    llama_binary,
    '-m', model_path,
    '--n-gpu-layers', '0',
    '-p', prompt,
    '--ctx-size', '2048',
    '--temp', '0.1',
    '-n', '100'
]

print(f"Testing with prompt: {prompt}")
print("Starting translation...")

try:
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
    if result.returncode != 0:
        print(f"llama.cpp failed with code {result.returncode}")
        print(f"Stderr: {result.stderr}")
        exit(1)
    
    # Just take the output as translation
    translation = result.stdout.strip()
    
    print("\n=== Translation Result ===")
    print(f"Original: {text}")
    print(f"Translated: {translation}")
    
    # Check if it contains Cyrillic
    has_cyrillic = any('\u0400' <= char <= '\u04FF' for char in translation)
    if has_cyrillic:
        print("✓ Translation contains Cyrillic characters")
    else:
        print("✗ Translation does not contain Cyrillic characters")
    
    if translation and translation != text:
        print("✓ Translation successful")
        
        # Now test with a variation that explicitly requests Cyrillic
        print("\n--- Testing with explicit Cyrillic request ---")
        prompt2 = f"Translate from Russian to Serbian (Cyrillic script): {text}"
        
        cmd2 = [
            llama_binary,
            '-m', model_path,
            '--n-gpu-layers', '0',
            '-p', prompt2,
            '--ctx-size', '2048',
            '--temp', '0.1',
            '-n', '100'
        ]
        
        result2 = subprocess.run(cmd2, capture_output=True, text=True, timeout=60)
        if result2.returncode == 0:
            translation2 = result2.stdout.strip()
            print(f"With explicit request: {translation2}")
            
            has_cyrillic2 = any('\u0400' <= char <= '\u04FF' for char in translation2)
            if has_cyrillic2:
                print("✓ Explicit request produced Cyrillic!")
            else:
                print("✗ Still no Cyrillic with explicit request")
        else:
            print("Explicit request failed")
    else:
        print("✗ Translation failed")
        
except subprocess.TimeoutExpired:
    print("Translation timed out")
except Exception as e:
    print(f"Translation failed: {e}")

PYEOF

chmod +x test_fixed.py

echo
echo "2. Running fixed test..."
python3 test_fixed.py

EOF