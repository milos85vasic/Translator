#!/bin/bash

echo "=== Testing Better Extraction ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Testing with better extraction logic..."

cat > test_extraction.py << 'PYEOF'
#!/usr/bin/env python3
import os
import subprocess

def translate_text(text, max_tokens=200):
    """Translate text using llama.cpp with better extraction"""
    llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
    model_path = "/home/milosvasic/models/tiny-llama-working.gguf"

    if not os.path.exists(model_path):
        model_path = "/home/milosvasic/models/tiny-llama.gguf"

    # Simple direct instruction
    prompt = f"Russian to Serbian Cyrillic: {text}"

    cmd = [
        llama_binary,
        '-m', model_path,
        '--n-gpu-layers', '0',
        '-p', prompt,
        '--ctx-size', '2048',
        '--temp', '0.1',
        '-n', str(max_tokens)
    ]

    print(f"Prompt: {prompt}")
    print("Running translation...")
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        if result.returncode != 0:
            print(f"llama.cpp failed with code {result.returncode}")
            print(f"Stderr: {result.stderr}")
            return ""
        
        output = result.stdout
        
        # Try different extraction strategies
        print("\n--- Raw Output ---")
        print(repr(output))
        
        # Strategy 1: Look for content after the prompt
        if prompt in output:
            translation = output.split(prompt, 1)[1].strip()
        # Strategy 2: Look for content after newline
        elif '\n' in output:
            lines = output.split('\n')
            # Skip first line if it's the prompt echo
            if lines[0].strip().startswith("Russian to Serbian Cyrillic:"):
                translation = '\n'.join(lines[1:]).strip()
            else:
                translation = output.strip()
        else:
            translation = output.strip()
        
        # Clean up any remaining prompt echoes
        translation = translation.replace(prompt, "").strip()
        
        # Remove any "translation:" markers
        if "translation:" in translation.lower():
            parts = translation.lower().split("translation:", 1)
            if len(parts) > 1:
                translation = parts[1].strip()
        
        print("\n--- Extracted Translation ---")
        print(repr(translation))
        
        # Check if it contains Cyrillic
        has_cyrillic = any('\u0400' <= char <= '\u04FF' for char in translation)
        
        # Also check if it's not just echoing Russian
        is_different = translation != text
        
        return {
            'text': translation,
            'has_cyrillic': has_cyrillic,
            'is_different': is_different
        }
        
    except subprocess.TimeoutExpired:
        print("Translation timed out")
        return {'text': '', 'has_cyrillic': False, 'is_different': False}
    except Exception as e:
        print(f"Translation failed: {e}")
        return {'text': '', 'has_cyrillic': False, 'is_different': False}

# Test with multiple examples
test_texts = [
    "Это простая тестовая книга на русском языке.",
    "Снег танцевал в свете фонаря.",
    "Кровь капала с воротничка его рубашки на снег."
]

for i, text in enumerate(test_texts):
    print(f"\n=== Test {i+1} ===")
    result = translate_text(text)
    
    print(f"Original: {text}")
    print(f"Translation: {result['text']}")
    print(f"Contains Cyrillic: {result['has_cyrillic']}")
    print(f"Is different: {result['is_different']}")
    
    if result['has_cyrillic'] and result['is_different']:
        print("✓ SUCCESS!")
    else:
        print("✗ FAILED")
    
    print("-" * 50)

PYEOF

chmod +x test_extraction.py

echo
echo "2. Running extraction test..."
python3 test_extraction.py

EOF