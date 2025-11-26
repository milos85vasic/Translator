#!/bin/bash

echo "=== Checking Architecture and Fixing llama.cpp ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Checking system architecture..."
arch
uname -m

echo
echo "2. Checking uploaded llama.cpp binary..."
file ./llama.cpp 2>/dev/null || echo "No llama.cpp file found"
file /home/milosvasic/llama.cpp 2>/dev/null || echo "No local llama.cpp found"

echo
echo "3. Checking if we need to compile llama.cpp..."
# Check if we can use the existing llama.cpp
if [ -f "/home/milosvasic/llama.cpp" ]; then
    echo "Found existing llama.cpp binary, testing if it works..."
    timeout 10 /home/milosvasic/llama.cpp --help || echo "Existing llama.cpp doesn't work"
fi

echo
echo "4. Setting up translation with existing llama.cpp..."
# Create a fixed translation script
cat > fixed_translate.py << 'PYEOF'
#!/usr/bin/env python3
import sys
import os
import subprocess

def translate_text(text):
    """Translate using llama.cpp"""
    # Use the existing llama.cpp binary
    llama_binary = "/home/milosvasic/llama.cpp"
    
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
    
    print(f"Running: {llama_binary} -m {model_path}")
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        if result.returncode != 0:
            print(f"llama.cpp failed with code {result.returncode}")
            print(f"Stderr: {result.stderr}")
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
        
        return '\n'.join(result_lines).strip()
        
    except subprocess.TimeoutExpired:
        print("Translation timed out")
        return ""
    except Exception as e:
        print(f"Translation failed: {e}")
        return ""

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 fixed_translate.py <input.md> <output.md>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    # Read input
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Translate just the first paragraph for testing
    paragraphs = [p for p in content.split('\n\n') if p.strip()]
    if paragraphs:
        # Skip title headers and translate first actual paragraph
        first_paragraph = None
        for p in paragraphs:
            if not p.startswith('#'):
                first_paragraph = p
                break
        
        if first_paragraph:
            print(f"Translating: {first_paragraph}")
            translated = translate_text(first_paragraph)
            
            # Write result
            with open(output_file, 'w', encoding='utf-8') as f:
                f.write(f"# Преведена содржина\n\n{translated}\n\n[Тест превода - само први пасус]")
            
            print(f"Translation written to {output_file}")
        else:
            print("No suitable paragraph found for translation")
    else:
        print("No content to translate")
PYEOF

chmod +x fixed_translate.py

echo
echo "5. Running translation test with fixed script..."
python3 fixed_translate.py "materials/books/test_book_small_original.md" "materials/books/test_book_small_translated.md"

echo
echo "6. Checking result..."
if [ -f "materials/books/test_book_small_translated.md" ]; then
    echo "Translation successful!"
    cat "materials/books/test_book_small_translated.md"
else
    echo "Translation still failed"
fi

EOF