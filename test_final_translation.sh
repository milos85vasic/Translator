#!/bin/bash

echo "=== Testing Translation with Correct llama.cpp Binary ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Testing the correct llama.cpp binary..."
LLAMA_BINARY="/home/milosvasic/llama.cpp/build/bin/llama-cli"
timeout 10 $LLAMA_BINARY --help | head -5

echo
echo "2. Testing llama.cpp with model..."
MODEL_PATH="/home/milosvasic/models/tiny-llama-working.gguf"
if [ ! -f "$MODEL_PATH" ]; then
    MODEL_PATH="/home/milosvasic/models/tiny-llama.gguf"
fi

echo "Testing with model: $MODEL_PATH"
timeout 30 $LLAMA_BINARY -m "$MODEL_PATH" -p "Hello" -n 10 --temp 0.1 || echo "Test timed out"

echo
echo "3. Running translation test..."
# Create final translation script
cat > final_translate.py << 'PYEOF'
#!/usr/bin/env python3
import sys
import os
import subprocess

def translate_text(text):
    """Translate using llama.cpp"""
    # Use the correct llama.cpp binary
    llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
    
    # Find model
    model_path = "/home/milosvasic/models/tiny-llama-working.gguf"
    if not os.path.exists(model_path):
        model_path = "/home/milosvasic/models/tiny-llama.gguf"
    
    if not os.path.exists(model_path):
        print("No model found!")
        return ""
    
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
        '-n', '500'
    ]
    
    print(f"Running translation...")
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        if result.returncode != 0:
            print(f"llama.cpp failed with code {result.returncode}")
            if result.stderr:
                print(f"Stderr: {result.stderr[:200]}")
            return ""
        
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
        
        # Verify we got a translation
        if not translation:
            print("No translation extracted from output")
            print("Full output:")
            print(result.stdout[:500])
        
        return translation
        
    except subprocess.TimeoutExpired:
        print("Translation timed out")
        return ""
    except Exception as e:
        print(f"Translation failed: {e}")
        return ""

def translate_markdown_file(input_file, output_file):
    """Translate entire markdown file"""
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Process content paragraph by paragraph
    paragraphs = [p for p in content.split('\n\n') if p.strip()]
    translated_paragraphs = []
    
    for i, paragraph in enumerate(paragraphs):
        print(f"Translating paragraph {i+1}/{len(paragraphs)}")
        
        # Skip title headers
        if paragraph.startswith('#'):
            translated_paragraphs.append(paragraph)
            continue
        
        # Translate the paragraph
        translated = translate_text(paragraph)
        
        if translated:
            translated_paragraphs.append(translated)
        else:
            print(f"Failed to translate paragraph {i+1}")
            translated_paragraphs.append(f"[Translation failed for: {paragraph[:50]}...]")
    
    # Write translated content
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write('\n\n'.join(translated_paragraphs))
    
    print(f"Translation completed: {output_file}")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 final_translate.py <input.md> <output.md>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    translate_markdown_file(input_file, output_file)
PYEOF

chmod +x final_translate.py

echo
echo "4. Running full translation..."
python3 final_translate.py "materials/books/test_book_small_original.md" "materials/books/test_book_small_translated.md"

echo
echo "5. Checking final result..."
if [ -f "materials/books/test_book_small_translated.md" ]; then
    echo "Translation successful!"
    echo "--- Translation Result ---"
    cat "materials/books/test_book_small_translated.md"
    echo "-------------------------"
    
    # Check for Cyrillic characters
    if grep -q '[\u0400-\u04FF]' "materials/books/test_book_small_translated.md"; then
        echo "✓ Translation contains Cyrillic characters"
    else
        echo "✗ Translation does not contain Cyrillic characters"
    fi
else
    echo "Translation failed"
fi

EOF