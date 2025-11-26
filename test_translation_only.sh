#!/bin/bash

echo "=== Testing Translation Step Only ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Setting up llama.cpp..."
# Check if llama.cpp binary exists and works
if [ -f "./llama.cpp" ]; then
    chmod +x ./llama.cpp
    echo "Testing llama.cpp binary..."
    ./llama.cpp --help | head -3
else
    echo "llama.cpp binary not found in current directory"
    # Try the existing one
    if [ -f "/home/milosvasic/llama.cpp" ]; then
        echo "Found llama.cpp at /home/milosvasic/llama.cpp"
        cp /home/milosvasic/llama.cpp ./llama.cpp
        chmod +x ./llama.cpp
        ./llama.cpp --help | head -3
    fi
fi

echo
echo "2. Testing llama.cpp with model..."
# Find model
MODEL_PATH="/home/milosvasic/models/tiny-llama-working.gguf"
if [ ! -f "$MODEL_PATH" ]; then
    MODEL_PATH="/home/milosvasic/models/tiny-llama.gguf"
fi

if [ -f "$MODEL_PATH" ]; then
    echo "Found model at: $MODEL_PATH"
    
    # Test llama.cpp with a simple prompt
    echo "Testing llama.cpp with a simple prompt..."
    echo "Test prompt: 'Hello'" | timeout 30 ./llama.cpp -m "$MODEL_PATH" --n-gpu-layers 0 -p "Hello" -n 10 --temp 0.1 || echo "llama.cpp test failed or timed out"
else
    echo "No model found!"
fi

echo
echo "3. Setting up translation script..."
# Create a simplified translation script
cat > simple_translate.py << 'PYEOF'
#!/usr/bin/env python3
import sys
import subprocess
import time

def translate_text(text):
    """Translate using llama.cpp"""
    llama_binary = "./llama.cpp"
    model_path = "/home/milosvasic/models/tiny-llama-working.gguf"
    
    # Try alternative model if first one doesn't exist
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
        '-n', '500'
    ]
    
    print(f"Running command: {' '.join(cmd)}")
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        print(f"Return code: {result.returncode}")
        if result.stdout:
            print(f"Output length: {len(result.stdout)}")
            print("First 500 chars of output:")
            print(result.stdout[:500])
        if result.stderr:
            print(f"Stderr: {result.stderr}")
        
        return result.stdout
    except subprocess.TimeoutExpired:
        print("Translation timed out")
        return ""
    except Exception as e:
        print(f"Translation failed: {e}")
        return ""

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 simple_translate.py <input.md> <output.md>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    # Read input
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    print(f"Input content: {content[:100]}...")
    
    # Translate just the first paragraph for testing
    paragraphs = content.split('\n\n')
    if paragraphs:
        first_paragraph = paragraphs[0]
        print(f"Translating: {first_paragraph}")
        translated = translate_text(first_paragraph)
        
        # Write result
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(f"# Translated Content\n\n{translated}\n\n[Translation test - first paragraph only]")
        
        print(f"Translation written to {output_file}")
    else:
        print("No content to translate")
PYEOF

chmod +x simple_translate.py

echo
echo "4. Running translation test..."
python3 simple_translate.py "materials/books/test_book_small_original.md" "materials/books/test_book_small_translated.md"

echo
echo "5. Checking translation result..."
if [ -f "materials/books/test_book_small_translated.md" ]; then
    echo "Translation file created:"
    cat "materials/books/test_book_small_translated.md"
else
    echo "Translation failed"
fi

EOF