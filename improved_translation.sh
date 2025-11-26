#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"
LLAMA_BINARY="/home/milosvasic/llama.cpp/build/bin/llama-cli"
MODEL_PATH="/home/milosvasic/models/tiny-llama-working.gguf"

# Create remote directory
echo "=== Creating remote directory ==="
ssh $USER@$HOST "mkdir -p $REMOTE_DIR"

# Upload the original FB2
echo -e "\n=== Uploading FB2 file ==="
scp /Users/milosvasic/Projects/Translate/internal/materials/books/book1.fb2 $USER@$HOST:$REMOTE_DIR/

# Create a better extraction and translation script
echo -e "\n=== Creating translation script ==="
ssh $USER@$HOST "cat > $REMOTE_DIR/extract_and_translate.sh << 'EOF'
#!/bin/bash

# Extract text from FB2 using Python
python3 << 'PYTHON_EOF'
import re
import sys

# Read FB2 file
with open('book1.fb2', 'r', encoding='utf-8') as f:
    content = f.read()

# Extract paragraphs after Chapter 1
chapter_match = re.search(r'Глава 1.*?</title>(.*?)(?=Глава|$)', content, re.DOTALL)
if chapter_match:
    chapter_content = chapter_match.group(1)
else:
    # Fallback: extract paragraphs after the title page
    start_match = re.search(r'</cite></section><section>', content)
    if start_match:
        chapter_content = content[start_match.end():]
        # Limit to first 10000 characters
        chapter_content = chapter_content[:10000]
    else:
        chapter_content = content[10000:15000]

# Extract paragraphs using regex
paragraphs = re.findall(r'<p>([^<]*)</p>', chapter_content)

# Filter meaningful paragraphs (longer than 30 chars, not copyright/info)
meaningful = []
for p in paragraphs:
    p = p.strip()
    if (len(p) > 30 and 
        not any(skip in p.lower() for skip in ['copyright', '©', 'издание', 'isbn', 'jo nesbø']) and
        not re.match(r'^[A-Z\s\-\.]+$', p) and
        not 'глава' in p.lower()):
        meaningful.append(p)

# Take first 5 paragraphs
selected = meaningful[:5]
print(f'Extracted {len(selected)} paragraphs for translation')

# Save to temporary file
with open('russian_paragraphs.txt', 'w', encoding='utf-8') as f:
    for i, p in enumerate(selected, 1):
        f.write(f'{p}\\n\\n')

print('Saved paragraphs to russian_paragraphs.txt')
PYTHON_EOF

# Now translate each paragraph
echo "=== Translating paragraphs ==="

mkdir -p translations

# Read paragraphs and translate
python3 << 'TRANSLATE_EOF'
import subprocess
import re
import os

def translate_paragraph(text, max_attempts=3):
    """Translate a paragraph using llama.cpp"""
    
    for attempt in range(max_attempts):
        # Try different prompt strategies
        if attempt == 0:
            prompt = f"Translate to Serbian Cyrillic: {text}"
        elif attempt == 1:
            prompt = f"Преведи на српски ћирилицу: {text}"
        else:
            prompt = f"Convert this Russian to Serbian Cyrillic writing: {text}"
        
        # Run llama.cpp with minimal output
        cmd = [
            '$LLAMA_BINARY',
            '-m', '$MODEL_PATH',
            '--prompt', prompt,
            '-n', '300',
            '--temp', '0.3',  # Lower temp for more consistent output
            '--repeat-penalty', '1.1',
            '--no-display-prompt',
            '--log-disable'
        ]
        
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
            output = result.stdout
            
            # Extract translation - look for Cyrillic content
            lines = output.split('\\n')
            translation = ''
            
            for line in lines:
                # Skip prompt and system messages
                if any(skip in line.lower() for skip in ['translate', 'prompt', 'system', 'main:', 'llama_', 'ggml', 'build', 'sampler', 'common']):
                    continue
                
                # Check if line contains Cyrillic
                if re.search(r'[А-ЯЉЊЋЂЏ][а-яљњћђџ]', line):
                    translation += line.strip() + ' '
            
            # Clean the translation
            if translation:
                # Remove any remaining artifacts
                translation = re.sub(r'[^А-Яа-яЉљЊњЋћЂђЏџ\\s.,!?;:\\-\\'\"()]', '', translation)
                translation = re.sub(r'\\s+', ' ', translation).strip()
                
                # Basic validation - should have at least some Cyrillic
                if len(translation) > 10 and re.search(r'[А-ЯЉЊЋЂЏ][а-яљњћђџ]', translation):
                    return translation
        except Exception as e:
            print(f'Translation attempt {attempt + 1} failed: {e}')
            continue
    
    # If all attempts fail, return original with note
    return f"[Translation failed] {text}"

# Read the paragraphs
with open('russian_paragraphs.txt', 'r', encoding='utf-8') as f:
    paragraphs_text = f.read()

paragraphs = [p.strip() for p in paragraphs_text.split('\\n\\n') if p.strip()]

# Translate each paragraph
translations = []
for i, paragraph in enumerate(paragraphs, 1):
    print(f'Translating paragraph {i}/{len(paragraphs)}...')
    print(f'Original: {paragraph[:100]}...')
    
    translation = translate_paragraph(paragraph)
    translations.append(translation)
    
    print(f'Translation: {translation[:100]}...')
    print()

# Save the translation
with open('book1_translated_sample.md', 'w', encoding='utf-8') as f:
    f.write('# Превод: Крв на снегу\\n\\n')
    f.write('# Преузето из руског оригинала\\n\\n')
    
    for i, (orig, trans) in enumerate(zip(paragraphs, translations), 1):
        f.write(f'## Параграф {i}\\n\\n')
        f.write(f'{trans}\\n\\n')

print('Translation complete! Saved to book1_translated_sample.md')
TRANSLATE_EOF
EOF"

# Make script executable and run it
echo -e "\n=== Running extraction and translation ==="
ssh $USER@$HOST "chmod +x $REMOTE_DIR/extract_and_translate.sh && cd $REMOTE_DIR && ./extract_and_translate.sh"

# Download the result
echo -e "\n=== Downloading translation ==="
scp $USER@$HOST:$REMOTE_DIR/book1_translated_sample.md /Users/milosvasic/Projects/Translate/internal/materials/books/

# Display the result
echo -e "\n=== Translation Result ==="
cat /Users/milosvasic/Projects/Translate/internal/materials/books/book1_translated_sample.md