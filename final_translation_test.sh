#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"
LLAMA_BINARY="/home/milosvasic/llama.cpp/build/bin/llama-cli"
MODEL_PATH="/home/milosvasic/models/tiny-llama-working.gguf"

# Create remote directory and run translation
echo "=== Running translation ==="
ssh $USER@$HOST "
mkdir -p $REMOTE_DIR && cd $REMOTE_DIR

# Extract text from FB2
python3 << 'PYEOF'
import re

# Read FB2 file
with open('book1.fb2', 'r', encoding='utf-8') as f:
    content = f.read()

# Extract paragraphs after Chapter 1
chapter_match = re.search(r'Глава 1.*?</title>(.*?)(?=Глава|</section>)', content, re.DOTALL)
if chapter_match:
    chapter_content = chapter_match.group(1)
else:
    start_match = re.search(r'</cite></section><section>', content)
    if start_match:
        chapter_content = content[start_match.end():]
        chapter_content = chapter_content[:10000]
    else:
        chapter_content = content[10000:15000]

# Extract paragraphs
paragraphs = re.findall(r'<p>([^<]*)</p>', chapter_content)

# Filter meaningful paragraphs
meaningful = []
for p in paragraphs:
    p = p.strip()
    if (len(p) > 30 and 
        not any(skip in p.lower() for skip in ['copyright', '©', 'издание', 'isbn', 'jo nesbø']) and
        not re.match(r'^[A-Z\s\-\.\',:]+$', p) and
        'глава' not in p.lower()):
        meaningful.append(p)

# Take first 3 paragraphs
selected = meaningful[:3]
print(f'Extracted {len(selected)} paragraphs')

# Save to file
with open('russian_paragraphs.txt', 'w', encoding='utf-8') as f:
    for p in selected:
        f.write(p + '\n\n')
PYEOF

# Create translation script
cat > translate_single.py << 'PYEOF'
import subprocess
import re
import sys

def translate_text(text):
    '''Translate using llama.cpp'''
    prompt = f'Преведи на српски ћирилицу:\n\n{text}'
    
    cmd = [
        '$LLAMA_BINARY',
        '-m', '$MODEL_PATH',
        '--prompt', prompt,
        '-n', '300',
        '--temp', '0.3',
        '--no-display-prompt',
        '--log-disable'
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
    output = result.stdout
    
    # Find Cyrillic content
    lines = output.split('\n')
    translation = ''
    
    for line in lines:
        # Skip prompt and system messages
        if any(skip in line.lower() for skip in ['translate', 'prompt', 'system', 'main:', 'llama_', 'ggml', 'build']):
            continue
        
        # Check for Cyrillic
        if re.search(r'[А-ЯЉЊЋЂЏ][а-яљњћђџ]', line):
            translation += line.strip() + ' '
    
    # Clean up
    if translation:
        translation = re.sub(r'[^А-Яа-яЉљЊњЋћЂђЏџ\s\.,!?;:\-\"\'()]', '', translation)
        translation = re.sub(r'\s+', ' ', translation).strip()
    
    return translation if translation else '[Translation failed]'

# Read and translate
with open('russian_paragraphs.txt', 'r', encoding='utf-8') as f:
    paragraphs_text = f.read()

paragraphs = [p.strip() for p in paragraphs_text.split('\n\n') if p.strip()]

# Create translation
with open('book1_translated_sample.md', 'w', encoding='utf-8') as f:
    f.write('# Превод: Крв на снегу\n\n')
    f.write('# Први пасуси из књige\n\n')
    
    for i, paragraph in enumerate(paragraphs, 1):
        print(f'Translating paragraph {i}...')
        translation = translate_text(paragraph)
        
        f.write(f'## Пасус {i}\n\n')
        f.write(f'{translation}\n\n')
        
        print(f'  Done: {translation[:50]}...')

print('Translation complete!')
PYEOF

# Run translation
python3 translate_single.py
"

# Download the result
echo -e "\n=== Downloading translation ==="
scp $USER@$HOST:$REMOTE_DIR/book1_translated_sample.md /Users/milosvasic/Projects/Translate/internal/materials/books/

# Display the result
echo -e "\n=== Translation Result ==="
cat /Users/milosvasic/Projects/Translate/internal/materials/books/book1_translated_sample.md