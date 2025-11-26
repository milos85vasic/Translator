#!/bin/bash

echo "=== Testing Single Paragraph Translation ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Killing all translation processes..."
pkill -f "python.*translate"

echo
echo "2. Creating a simple test with just one paragraph..."
cat > test_single.py << 'PYEOF'
#!/usr/bin/env python3
import os
import subprocess

# Test single paragraph translation
text = "Это простая тестовая книга на русском языке."

llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
model_path = "/home/milosvasic/models/tiny-llama-working.gguf"

if not os.path.exists(model_path):
    model_path = "/home/milosvasic/models/tiny-llama.gguf"

prompt = f"""You are a professional translator from Russian to Serbian. 
Translate the following text. Provide ONLY the Serbian translation, no explanations.

Original text:
{text}

Serbian translation:"""

cmd = [
    llama_binary,
    '-m', model_path,
    '--n-gpu-layers', '0',
    '-p', prompt,
    '--ctx-size', '2048',
    '--temp', '0.1',
    '-n', '100'
]

print(f"Running: {' '.join(cmd[:3])}...")
print(f"Model: {model_path}")
print(f"Text to translate: {text}")
print("Starting translation...")

try:
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
    if result.returncode != 0:
        print(f"llama.cpp failed with code {result.returncode}")
        print(f"Stderr: {result.stderr}")
        exit(1)
    
    # Extract translation from output
    lines = result.stdout.split('\n')
    translation_started = False
    result_lines = []
    
    for line in lines:
        if "Serbian translation:" in line:
            translation_started = True
            parts = line.split("Serbian translation:", 1)
            if len(parts) > 1 and parts[1].strip():
                result_lines.append(parts[1].strip())
            continue
        elif translation_started and line.strip():
            result_lines.append(line.strip())
    
    translation = '\n'.join(result_lines).strip()
    
    print("\n=== Translation Result ===")
    print(f"Original: {text}")
    print(f"Translated: {translation}")
    
    # Check if it contains Cyrillic
    if any('\u0400' <= char <= '\u04FF' for char in translation):
        print("✓ Translation contains Cyrillic characters")
    else:
        print("✗ Translation does not contain Cyrillic characters")
    
    if translation and translation != text:
        print("✓ Translation successful and different from original")
    else:
        print("✗ Translation failed or same as original")
        
except subprocess.TimeoutExpired:
    print("Translation timed out")
except Exception as e:
    print(f"Translation failed: {e}")

PYEOF

chmod +x test_single.py

echo
echo "3. Running single paragraph test..."
python3 test_single.py

EOF