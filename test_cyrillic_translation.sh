#!/bin/bash

echo "=== Testing Translation with Cyrillic Output ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Testing with corrected prompt for Cyrillic output..."

cat > test_cyrillic.py << 'PYEOF'
#!/usr/bin/env python3
import os
import subprocess

# Test paragraph translation with specific instructions for Cyrillic
text = "Это простая тестовая книга на русском языке."

llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
model_path = "/home/milosvasic/models/tiny-llama-working.gguf"

if not os.path.exists(model_path):
    model_path = "/home/milosvasic/models/tiny-llama.gguf"

# Fixed prompt that explicitly requests Cyrillic and disables chat mode
prompt = f"""Text: {text}

Translate from Russian to Serbian Cyrillic script:"""

cmd = [
    llama_binary,
    '-m', model_path,
    '--n-gpu-layers', '0',
    '-p', prompt,
    '--ctx-size', '2048',
    '--temp', '0.1',
    '-n', '100',
    '--no-cnv'  # Disable conversation mode
]

print(f"Text to translate: {text}")
print("Starting translation with corrected prompt...")

try:
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
    if result.returncode != 0:
        print(f"llama.cpp failed with code {result.returncode}")
        print(f"Stderr: {result.stderr}")
        exit(1)
    
    # Extract just the translated text (everything after the prompt marker)
    if "Translate from Russian to Serbian Cyrillic script:" in result.stdout:
        translation = result.stdout.split("Translate from Russian to Serbian Cyrillic script:", 1)[1].strip()
    else:
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
        # Test with actual book content
        print("\n--- Testing with real book content ---")
        book_text = "Снег танцевал в свете фонаря, подобно хлопковому пуху."
        
        prompt2 = f"""Text: {book_text}

Translate from Russian to Serbian Cyrillic script:"""
        
        cmd2 = [
            llama_binary,
            '-m', model_path,
            '--n-gpu-layers', '0',
            '-p', prompt2,
            '--ctx-size', '2048',
            '--temp', '0.1',
            '-n', '100',
            '--no-cnv'
        ]
        
        result2 = subprocess.run(cmd2, capture_output=True, text=True, timeout=60)
        if result2.returncode == 0:
            if "Translate from Russian to Serbian Cyrillic script:" in result2.stdout:
                book_translation = result2.stdout.split("Translate from Russian to Serbian Cyrillic script:", 1)[1].strip()
            else:
                book_translation = result2.stdout.strip()
            
            print(f"Book text: {book_text}")
            print(f"Translation: {book_translation}")
        else:
            print("Book translation test failed")
    else:
        print("✗ Translation failed")
        
except subprocess.TimeoutExpired:
    print("Translation timed out")
except Exception as e:
    print(f"Translation failed: {e}")

PYEOF

chmod +x test_cyrillic.py

echo
echo "2. Running Cyrillic test..."
python3 test_cyrillic.py

EOF