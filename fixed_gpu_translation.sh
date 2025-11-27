#!/bin/bash
# Fixed GPU Translation with Proper Prompt Engineering

echo "=== Fixed GPU Translation ==="

# Kill any existing processes
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && pkill -f 'gpu_translate.py' 2>/dev/null; pkill -f 'simple_gpu_translate.py' 2>/dev/null; true"

# Create a properly working translation script
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && cat > fixed_gpu_translate.py << 'EOF'
#!/usr/bin/env python3
import sys
import os
import subprocess
import time

def translate_text(text, model_path, llama_binary):
    \"\"\"Translate a single text using llama.cpp with proper prompt\"\"\"
    
    # Better prompt that forces translation
    prompt = f\"\"\"<|im_start|>system
You are a professional Russian to Serbian translator. Translate only the given text.
<|im_end|>
<|im_start|>user
Translate to Serbian Cyrillic:
{text}
<|im_end|>
<|im_start|>assistant
\"\"\"
    
    cmd = [
        llama_binary,
        '-m', model_path,
        '--n-gpu-layers', '99',
        '-p', prompt,
        '--ctx-size', '2048',
        '--temp', '0.2',
        '-n', '1000',
        '--repeat-penalty', '1.1',
        '--in-prefix', ' ',
        '--in-suffix', 'assistant\n'
    ]
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        if result.returncode == 0:
            output = result.stdout.strip()
            # Extract the assistant's response
            if '<|im_start|>assistant' in output:
                translation = output.split('<|im_start|>assistant')[-1].strip()
            else:
                translation = output.strip()
            
            # Clean up any remaining markers
            translation = translation.replace('<|im_end|>', '').strip()
            return translation
        else:
            print(f'llama.cpp error: {result.stderr}')
            return text
    except Exception as e:
        print(f'Translation error: {e}')
        return text

def main():
    if len(sys.argv) != 3:
        print('Usage: python3 fixed_gpu_translate.py input.md output.md')
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    print(f'Starting GPU-accelerated translation: {input_file} -> {output_file}')
    
    model_path = '/home/milosvasic/models/tiny-llama-working.gguf'
    llama_binary = '/home/milosvasic/llama.cpp/build/bin/llama-cli'
    
    # Verify model exists
    if not os.path.exists(model_path):
        print(f'Model not found: {model_path}')
        sys.exit(1)
    
    # Read input
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split into paragraphs
    paragraphs = [p.strip() for p in content.split('\\n\\n') if p.strip()]
    print(f'Found {len(paragraphs)} paragraphs')
    
    # Test translation with first paragraph
    print('\\nTesting translation with first paragraph...')
    test_translation = translate_text(paragraphs[0], model_path, llama_binary)
    print(f'Original: {paragraphs[0][:100]}...')
    print(f'Translated: {test_translation[:100]}...')
    
    if test_translation == paragraphs[0]:
        print('Translation failed - output same as input')
        sys.exit(1)
    
    # Translate all paragraphs
    translated = []
    start_time = time.time()
    
    for i, paragraph in enumerate(paragraphs[:10]):  # Limit to 10 for testing
        print(f'Translating paragraph {i+1}/10...')
        
        translation = translate_text(paragraph, model_path, llama_binary)
        translated.append(translation)
        
        print(f'  ✓ {len(paragraph)} -> {len(translation)} chars')
    
    # Write output
    output_content = '\\n\\n'.join(translated)
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(output_content)
    
    elapsed = time.time() - start_time
    print(f'\\nTranslation completed in {elapsed:.1f} seconds!')
    print(f'Output: {output_file} ({len(output_content)} characters)')

if __name__ == '__main__':
    main()
EOF"

ssh milosvasic@thinker.local "cd /tmp/translate-ssh && chmod +x fixed_gpu_translate.py"

echo "Running fixed GPU translation..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && python3 fixed_gpu_translate.py book1_original.md book1_fixed_gpu.md"

echo "Checking results..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && echo '=== Output Stats ===' && wc -c book1_fixed_gpu.md && echo '=== Sample Output ===' && head -20 book1_fixed_gpu.md"