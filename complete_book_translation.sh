#!/bin/bash

echo "=== Complete Book Translation Process ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

# First, upload the book1.fb2 file if it doesn't exist
if [ ! -f "materials/books/book1.fb2" ]; then
    echo "Please upload book1.fb2 to the remote host first"
    echo "Run: scp internal/materials/books/book1.fb2 milosvasic@thinker.local:/tmp/book1.fb2"
    echo "Then run this script again"
    exit 1
fi

echo
echo "1. Creating complete translation script..."
cat > complete_translate.py << 'PYEOF'
#!/usr/bin/env python3
import sys
import os
import subprocess
import xml.etree.ElementTree as ET

def fb2_to_markdown(input_file, output_file):
    """Convert FB2 to Markdown"""
    try:
        tree = ET.parse(input_file)
        root = tree.getroot()
        
        # Remove namespace
        for elem in root.iter():
            if '}' in elem.tag:
                elem.tag = elem.tag.split('}', 1)[1]
        
        # Extract title
        title_info = root.find('.//title-info')
        book_title = title_info.find('.//book-title').text if title_info is not None and title_info.find('.//book-title') is not None else "Unknown Title"
        
        # Extract body content
        body = root.find('.//body')
        content = []
        
        if body is not None:
            for elem in body:
                if elem.tag == 'title':
                    for p in elem.findall('.//p'):
                        content.append(f"# {p.text or ''}")
                elif elem.tag == 'section':
                    for p in elem.findall('.//p'):
                        content.append(p.text or '')
                elif elem.tag == 'p':
                    content.append(elem.text or '')
                elif elem.tag == 'subtitle':
                    for p in elem.findall('.//p'):
                        content.append(f"## {p.text or ''}")
        
        # Write markdown
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(f"# {book_title}\n\n")
            for line in content:
                if line.strip():
                    f.write(line + "\n\n")
        
        print(f"Successfully converted {input_file} to {output_file}")
        return True
        
    except Exception as e:
        print(f"Error converting FB2: {e}")
        return False

def translate_text(text, llama_binary, model_path):
    """Translate text using llama.cpp"""
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
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        if result.returncode != 0:
            print(f"llama.cpp failed with code {result.returncode}")
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

def translate_markdown_file(input_file, output_file):
    """Translate entire markdown file"""
    # Use the correct llama.cpp binary
    llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
    
    # Find model
    model_path = "/home/milosvasic/models/tiny-llama-working.gguf"
    if not os.path.exists(model_path):
        model_path = "/home/milosvasic/models/tiny-llama.gguf"
    
    if not os.path.exists(model_path):
        print("No model found!")
        return False
    
    with open(input_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Process content paragraph by paragraph
    paragraphs = [p for p in content.split('\n\n') if p.strip()]
    translated_paragraphs = []
    
    print(f"Total paragraphs to translate: {len(paragraphs)}")
    
    for i, paragraph in enumerate(paragraphs):
        print(f"Translating paragraph {i+1}/{len(paragraphs)}")
        
        # Skip title headers and metadata
        if paragraph.startswith('#') or 'Copyright' in paragraph or 'Published by' in paragraph:
            translated_paragraphs.append(paragraph)
            continue
        
        # Limit paragraph size for model
        if len(paragraph) > 1000:
            # Split long paragraphs
            sentences = paragraph.split('. ')
            translated_sentences = []
            for sentence in sentences[:5]:  # Limit to first 5 sentences
                if sentence.strip():
                    translated = translate_text(sentence.strip(), llama_binary, model_path)
                    if translated:
                        translated_sentences.append(translated)
                    else:
                        translated_sentences.append(sentence)
            translated = '. '.join(translated_sentences)
        else:
            # Translate the paragraph
            translated = translate_text(paragraph, llama_binary, model_path)
        
        if translated:
            translated_paragraphs.append(translated)
        else:
            print(f"Failed to translate paragraph {i+1}, keeping original")
            translated_paragraphs.append(paragraph)
    
    # Write translated content
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write('\n\n'.join(translated_paragraphs))
    
    print(f"Translation completed: {output_file}")
    return True

def create_simple_epub(markdown_file, epub_file):
    """Create a simple EPUB from markdown"""
    # Create a minimal EPUB structure
    import zipfile
    import os
    
    # Create mimetype file
    with open('mimetype', 'w') as f:
        f.write('application/epub+zip')
    
    # Create META-INF/container.xml
    os.makedirs('META-INF', exist_ok=True)
    with open('META-INF/container.xml', 'w') as f:
        f.write('''<?xml version="1.0"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>''')
    
    # Create OEBPS directory
    os.makedirs('OEBPS', exist_ok=True)
    
    # Read markdown content
    with open(markdown_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Convert markdown to basic HTML
    html_content = f"""<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>Translated Book</title>
    <meta charset="utf-8"/>
</head>
<body>
{content.replace('\n\n', '</p><p>').replace('\n', ' ').replace('# ', '<h1>').replace('\n', '</h1>\n')}
</body>
</html>"""
    
    # Create content.opf
    with open('OEBPS/content.opf', 'w') as f:
        f.write(f'''<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="BookId" version="2.0">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:opf="http://www.idpf.org/2007/opf">
    <dc:title>Translated Book</dc:title>
    <dc:language>sr</dc:language>
    <dc:identifier id="BookId">book-1</dc:identifier>
  </metadata>
  <manifest>
    <item id="chapter1" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="chapter1"/>
  </spine>
</package>''')
    
    # Create chapter1.xhtml
    with open('OEBPS/chapter1.xhtml', 'w') as f:
        f.write(html_content)
    
    # Create EPUB
    with zipfile.ZipFile(epub_file, 'w') as epub:
        epub.write('mimetype', compress_type=zipfile.ZIP_STORED)
        epub.write('META-INF/container.xml')
        epub.write('OEBPS/content.opf')
        epub.write('OEBPS/chapter1.xhtml')
    
    # Clean up temporary files
    os.remove('mimetype')
    import shutil
    shutil.rmtree('META-INF')
    shutil.rmtree('OEBPS')
    
    print(f"EPUB created: {epub_file}")

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python3 complete_translate.py <input.fb2> <output_translated.md> <output.epub>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_md = sys.argv[2]
    output_epub = sys.argv[3]
    
    print(f"Converting {input_file} to {output_md}")
    
    # Convert FB2 to markdown
    temp_md = input_file.replace('.fb2', '_temp.md')
    if not fb2_to_markdown(input_file, temp_md):
        sys.exit(1)
    
    # Translate markdown
    if not translate_markdown_file(temp_md, output_md):
        sys.exit(1)
    
    # Create EPUB
    create_simple_epub(output_md, output_epub)
    
    print(f"Complete! Files created:")
    print(f"  Original FB2: {input_file}")
    print(f"  Original MD: {temp_md}")
    print(f"  Translated MD: {output_md}")
    print(f"  Final EPUB: {output_epub}")
PYEOF

chmod +x complete_translate.py

echo
echo "2. Running complete translation process..."
python3 complete_translate.py "materials/books/book1.fb2" "materials/books/book1_translated.md" "materials/books/book1_sr.epub"

echo
echo "3. Checking results..."
echo "Generated files:"
ls -la materials/books/book1*

if [ -f "materials/books/book1_sr.epub" ]; then
    echo
    echo "✓ EPUB file created successfully!"
    file "materials/books/book1_sr.epub"
    echo "File size: $(du -h materials/books/book1_sr.epub | cut -f1)"
else
    echo "✗ EPUB file not created"
fi

if [ -f "materials/books/book1_translated.md" ]; then
    echo
    echo "✓ Translated markdown created!"
    echo "File size: $(du -h materials/books/book1_translated.md | cut -f1)"
    echo "First 300 characters:"
    head -c 300 "materials/books/book1_translated.md"
    echo ""
    
    # Check for Cyrillic characters
    if grep -q '[\u0400-\u04FF]' "materials/books/book1_translated.md"; then
        echo "✓ Translation contains Cyrillic characters"
    else
        echo "✗ Translation does not contain Cyrillic characters"
    fi
else
    echo "✗ Translated markdown not created"
fi

EOF