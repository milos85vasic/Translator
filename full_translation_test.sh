#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"
LLAMA_BINARY="/home/milosvasic/llama.cpp/build/bin/llama-cli"
MODEL_PATH="/home/milosvasic/models/tiny-llama-working.gguf"

# Create remote directory
echo "=== Creating remote directory ==="
ssh $USER@$HOST "mkdir -p $REMOTE_DIR"

# Upload the full sample file
echo -e "\n=== Uploading sample file ==="
scp /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sample_original.md $USER@$HOST:$REMOTE_DIR/

# Create translation script on remote
echo -e "\n=== Creating translation script ==="
ssh $USER@$HOST "cat > $REMOTE_DIR/translate.py << 'EOF'
import re
import subprocess

# Read the Russian text
with open('book1_sample_original.md', 'r') as f:
    text = f.read().strip()

# Split into paragraphs
paragraphs = [p.strip() for p in text.split('\\n\\n') if p.strip()]
print(f'Found {len(paragraphs)} paragraphs to translate')

translated_paragraphs = []

for i, paragraph in enumerate(paragraphs, 1):
    print(f'Translating paragraph {i}/{len(paragraphs)}...')
    
    # Prepare prompt
    prompt = f'Translate this Russian text to Serbian Cyrillic. Return ONLY the Serbian translation, no explanations: {paragraph}'
    
    # Run llama.cpp
    cmd = [
        '$LLAMA_BINARY',
        '-m', '$MODEL_PATH',
        '--prompt', prompt,
        '-n', '200',
        '--temp', '0.7',
        '--ctx-size', '2048',
        '--no-display-prompt',
        '--color'
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True, cwd='$REMOTE_DIR')
    
    # Extract the actual translation (remove prompt echo)
    output = result.stdout
    
    # Find where the actual translation starts
    lines = output.split('\\n')
    translation = ''
    started = False
    
    for line in lines:
        # Skip prompt-related lines
        if 'Translate this Russian' in line or line.strip() == prompt[:50]:
            started = True
            continue
        elif started and line.strip():
            translation += line.strip() + ' '
        elif started and not line.strip() and translation:
            # Empty line after we started collecting translation
            break
    
    if not translation:
        # Fallback: extract last non-empty line that's not the prompt
        for line in reversed(lines):
            if line.strip() and 'Translate this Russian' not in line and '>' not in line and 'system' not in line.lower():
                translation = line.strip()
                break
    
    if translation:
        # Clean up the translation
        translation = re.sub(r'^[>\\s]*', '', translation)  # Remove leading > and spaces
        translation = re.sub(r'[^А-Яа-яЁёЉљЊњЋћЂђЏџЄєІіҐґ.,!?;:\\s\\-]', '', translation)  # Keep only Cyrillic and punctuation
        translated_paragraphs.append(translation)
        print(f'  → {translation[:100]}...')
    else:
        print(f'  → FAILED to extract translation')
        # Use original as fallback
        translated_paragraphs.append(paragraph)

# Write the translation
with open('book1_sample_translated.md', 'w') as f:
    for i, (orig, trans) in enumerate(zip(paragraphs, translated_paragraphs), 1):
        f.write(f'Параграф {i}:\\n')
        f.write(f'{trans}\\n\\n')

print('\\nTranslation complete!')
EOF"

# Run the translation
echo -e "\n=== Running translation ==="
ssh $USER@$HOST "cd $REMOTE_DIR && python3 translate.py"

# Download the result
echo -e "\n=== Downloading translation ==="
scp $USER@$HOST:$REMOTE_DIR/book1_sample_translated.md /Users/milosvasic/Projects/Translate/internal/materials/books/

# Display the result
echo -e "\n=== Translation Result ==="
cat /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sample_translated.md