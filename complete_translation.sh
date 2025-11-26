#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"
OUTPUT_FILE="/Users/milosvasic/Projects/Translate/internal/materials/books/book1_full_translated.md"

echo "=== Final Full Translation ==="
ssh $USER@$HOST "
cd $REMOTE_DIR

# Extract text from FB2
python3 << 'PYEOF'
import re

# Read FB2 file
with open('book1.fb2', 'r', encoding='utf-8') as f:
    content = f.read()

# Find content after title page
start = content.find('</cite></section><section>')
if start != -1:
    content = content[start:]

# Extract paragraphs
paragraphs = re.findall(r'<p>([^<]*)</p>', content)

# Filter and clean paragraphs
cleaned = []
for p in paragraphs[:10]:  # First 10 meaningful paragraphs
    p = p.strip()
    if (len(p) > 50 and 
        not any(skip in p.lower() for skip in ['copyright', '©', 'издание', 'isbn', 'jo nesbø', 'gribuser']) and
        not re.match(r'^[A-Z\s\-\.\',:\d]+$', p) and
        'глава' not in p.lower()):
        cleaned.append(p)

# Save to file
with open('to_translate.txt', 'w', encoding='utf-8') as f:
    for p in cleaned:
        f.write(p + '\n')

print(f'Prepared {len(cleaned)} paragraphs for translation')
PYEOF

# Create final translation script
cat > translate_final.py << 'PYEOF'
import subprocess
import re
import sys

def translate_with_chatgpt_style(text):
    '''Translate using a chat-like prompt'''
    
    # Create prompt that forces response
    prompt_text = f'''You are a translator. Translate the following Russian text to Serbian Cyrillic. Return ONLY the Serbian translation.

Russian: {text}
Serbian:'''
    
    # Run with --color to get cleaner output
    cmd = [
        '/home/milosvasic/llama.cpp/build/bin/llama-cli',
        '-m', '/home/milosvasic/models/tiny-llama-working.gguf',
        '--prompt', prompt_text,
        '-n', '200',
        '--temp', '0.2',  # Very low temp for consistency
        '--repeat-penalty', '1.0',
        '--no-display-prompt',
        '--log-disable'
    ]
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=20)
        output = result.stdout
        
        # Get everything after "Serbian:" in the output
        if 'Serbian:' in output:
            translation = output.split('Serbian:')[-1].strip()
        else:
            # Try to extract Cyrillic lines
            lines = output.split('\n')
            translation = ''
            collecting = False
            
            for line in lines:
                # Skip model stats
                if any(skip in line.lower() for skip in ['main:', 'llama_', 'ggml', 'build', 'sampler', 'perf']):
                    continue
                
                # Start collecting after we see our prompt pattern
                if 'Serbian:' in line or collecting:
                    collecting = True
                    # Keep the part after Serbian:
                    if 'Serbian:' in line:
                        translation = line.split('Serbian:')[-1].strip()
                    elif line.strip():
                        translation += ' ' + line.strip()
        
        # Clean up - keep only Cyrillic and punctuation
        if translation:
            translation = re.sub(r'[^А-Яа-яЉљЊњЋћЂђЏџ\s\.,!?;:\-\"\'()]', '', translation)
            translation = re.sub(r'\s+', ' ', translation).strip()
            return translation
        
    except Exception as e:
        print(f'Error: {e}')
    
    return None

# Read paragraphs
with open('to_translate.txt', 'r', encoding='utf-8') as f:
    paragraphs = [p.strip() for p in f if p.strip()]

print(f'Translating {len(paragraphs)} paragraphs...')

# Translate each
translations = []
for i, paragraph in enumerate(paragraphs, 1):
    print(f'\\nParagraph {i}: {paragraph[:50]}...')
    
    # Try translation
    translation = translate_with_chatgpt_style(paragraph)
    
    if translation and len(translation) > 10:
        print(f'→ {translation[:50]}...')
        translations.append(translation)
    else:
        print('→ Translation failed, keeping original')
        translations.append(f'[Failed] {paragraph}')

# Save final result
with open('book1_translated_final.md', 'w', encoding='utf-8') as f:
    f.write('# Крв на снегу\\n')
    f.write('# Превод с руског на српски ћирилицу\\n\\n')
    
    for i, (orig, trans) in enumerate(zip(paragraphs, translations), 1):
        f.write(f'## Пасус {i}\\n\\n')
        f.write(f'{trans}\\n\\n')

print('\\nTranslation complete!')
PYEOF

# Run translation
python3 translate_final.py
"

# Download the result
echo -e "\n=== Downloading translation ==="
scp $USER@$HOST:$REMOTE_DIR/book1_translated_final.md $OUTPUT_FILE

# Display part of the result
echo -e "\n=== First part of translation ==="
head -50 $OUTPUT_FILE