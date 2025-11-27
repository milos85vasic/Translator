#!/bin/bash
# Direct GPU Translation with Progress

echo "=== Restarting GPU Translation with Better Monitoring ==="

# Kill existing process
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && pkill -f 'gpu_translate.py' && sleep 2"

# Create a simpler translation script that shows progress
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && cat > simple_gpu_translate.py << 'EOF'
#!/usr/bin/env python3
import sys
import os
import subprocess
import time

def main():
    if len(sys.argv) != 3:
        print('Usage: python3 simple_gpu_translate.py input.md output.md')
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    print(f'Starting translation: {input_file} -> {output_file}')
    
    # Read input
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split into manageable chunks (first 10 paragraphs for testing)
    paragraphs = [p.strip() for p in content.split('\\n\\n') if p.strip()]
    paragraphs = paragraphs[:10]  # Limit to first 10 for testing
    
    print(f'Translating {len(paragraphs)} paragraphs...')
    
    model_path = '/home/milosvasic/models/tiny-llama-working.gguf'
    llama_binary = '/home/milosvasic/llama.cpp/build/bin/llama-cli'
    
    translated = []
    
    for i, paragraph in enumerate(paragraphs):
        print(f'Paragraph {i+1}/{len(paragraphs)}: {paragraph[:50]}...')
        
        # Simple prompt
        prompt = f'Преведи на српски ћирилицу:\\n{paragraph}'
        
        cmd = [
            llama_binary,
            '-m', model_path,
            '--n-gpu-layers', '99',
            '-p', prompt,
            '--ctx-size', '1024',
            '--temp', '0.1',
            '-n', '500'
        ]
        
        try:
            start = time.time()
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
            elapsed = time.time() - start
            
            if result.returncode == 0:
                translation = result.stdout.strip()
                translated.append(translation)
                print(f'  ✓ Done in {elapsed:.1f}s: {translation[:50]}...')
            else:
                print(f'  ✗ Failed: {result.stderr[:100]}')
                translated.append(paragraph)
        except Exception as e:
            print(f'  ✗ Error: {e}')
            translated.append(paragraph)
    
    # Write output
    output_content = '\\n\\n'.join(translated)
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(output_content)
    
    print(f'Translation completed! Output: {output_file}')

if __name__ == '__main__':
    main()
EOF"

ssh milosvasic@thinker.local "cd /tmp/translate-ssh && chmod +x simple_gpu_translate.py"

echo "Starting simple GPU translation..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && python3 simple_gpu_translate.py book1_original.md book1_simple_gpu.md"

echo "Checking results..."
ssh milosvasic@thinker.local "cd /tmp/translate-ssh && echo '=== Output File ===' && ls -la book1_simple_gpu.md && echo '=== First Paragraph ===' && head -10 book1_simple_gpu.md"