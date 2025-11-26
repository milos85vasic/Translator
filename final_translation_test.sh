#!/bin/bash

echo "=== Final Book Translation with Working Method ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Creating final translation script..."

cat > final_book_translate.py << 'PYEOF'
#!/usr/bin/env python3
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

def translate_text(text):
    """Translate text using llama.cpp"""
    llama_binary = "/home/milosvasic/llama.cpp/build/bin/llama-cli"
    model_path = "/home/milosvasic/models/tiny-llama-working.gguf"

    if not os.path.exists(model_path):
        model_path = "/home/milosvasic/models/tiny-llama.gguf"
    
    # Simple prompt that worked in tests
    prompt = f"Russian to Serbian Cyrillic: {text}"

    cmd = [
        llama_binary,
        '-m', model_path,
        '--n-gpu-layers', '0',
        '-p', prompt,
        '--ctx-size', '2048',
        '--temp', '0.1',
        '-n', '200'
    ]
    
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
        if result.returncode != 0:
            print(f"llama.cpp failed with code {result.returncode}")
            return ""
        
        output = result.stdout
        
        # Extract translation - look for content after "Russian to Serbian Cyrilic:" marker
        if "Russian to Serbian Cyrillic:" in output:
            translation = output.split("Russian to Serbian Cyrillic:", 1)[1].strip()
        elif "Russian to Serbian Cyrilic:" in output:
            translation = output.split("Russian to Serbian Cyrilic:", 1)[1].strip()
        else:
            # Take everything after first newline
            lines = output.split('\n')
            if len(lines) > 1:
                translation = '\n'.join(lines[1:]).strip()
            else:
                translation = output.strip()
        
        # Clean up any remaining artifacts
        translation = translation.replace(prompt, "").strip()
        translation = translation.replace("Russian to Serbian Cyrillic:", "").strip()
        
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
    
    print(f"Total paragraphs to translate: {len(paragraphs)}")
    
    for i, paragraph in enumerate(paragraphs):
        print(f"Translating paragraph {i+1}/{len(paragraphs)}")
        
        # Skip title headers and metadata
        if paragraph.startswith('#') or 'Copyright' in paragraph or 'Published by' in paragraph or '©' in paragraph:
            translated_paragraphs.append(paragraph)
            continue
        
        # Translate paragraph
        translated = translate_text(paragraph)
        
        if translated:
            # Verify it contains Cyrillic and is different
            has_cyrillic = any('\u0400' <= char <= '\u04FF' for char in translated)
            
            if has_cyrillic and translated != paragraph:
                translated_paragraphs.append(translated)
            else:
                # Try once more with different prompt
                alt_prompt = f"Translate to Serbian Cyrillic: {paragraph}"
                cmd = [
                    "/home/milosvasic/llama.cpp/build/bin/llama-cli",
                    '-m', model_path,
                    '--n-gpu-layers', '0',
                    '-p', alt_prompt,
                    '--ctx-size', '2048',
                    '--temp', '0.1',
                    '-n', '200'
                ]
                
                try:
                    result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
                    if result.returncode == 0:
                        output = result.stdout
                        if "Translate to Serbian Cyrillic:" in output:
                            alt_translation = output.split("Translate to Serbian Cyrillic:", 1)[1].strip()
                        else:
                            lines = output.split('\n')
                            if len(lines) > 1:
                                alt_translation = '\n'.join(lines[1:]).strip()
                            else:
                                alt_translation = output.strip()
                        
                        alt_translation = alt_translation.replace(alt_prompt, "").strip()
                        
                        if any('\u0400' <= char <= '\u04FF' for char in alt_translation):
                            translated_paragraphs.append(alt_translation)
                            print("  ✓ Second attempt successful")
                        else:
                            translated_paragraphs.append(paragraph)
                            print("  ✗ Second attempt failed, keeping original")
                    else:
                        translated_paragraphs.append(paragraph)
                        print("  ✗ Second attempt failed, keeping original")
                except:
                    translated_paragraphs.append(paragraph)
                    print("  ✗ Second attempt failed, keeping original")
        else:
            translated_paragraphs.append(paragraph)
            print("  ✗ Translation failed, keeping original")
    
    # Write translated content
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write('\n\n'.join(translated_paragraphs))
    
    print(f"Translation completed: {output_file}")
    return True

def create_simple_epub(markdown_file, epub_file):
    """Create a simple EPUB from markdown"""
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
    
    # Simple markdown to HTML conversion
    html_content = content
    # Convert headers
    html_content = html_content.replace('# ', '<h1>').replace('\n', '</h1>\n')
    html_content = html_content.replace('## ', '<h2>').replace('\n', '</h2>\n')
    # Convert paragraphs
    paragraphs = html_content.split('\n\n')
    html_paragraphs = []
    for p in paragraphs:
        if p.strip() and not p.startswith('<h'):
            html_paragraphs.append(f'<p>{p.strip()}</p>')
        else:
            html_paragraphs.append(p.strip())
    
    html_content = '\n'.join(html_paragraphs)
    
    # Create final HTML
    final_html = f'''<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title>Translated Book</title>
    <meta charset="utf-8"/>
</head>
<body>
{html_content}
</body>
</html>'''
    
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
        f.write(final_html)
    
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
        print("Usage: python3 final_book_translate.py <input.fb2> <output_translated.md> <output.epub>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_md = sys.argv[2]
    output_epub = sys.argv[3]
    
    print(f"Starting translation of {input_file}")
    
    # Convert FB2 to markdown
    temp_md = input_file.replace('.fb2', '_temp.md')
    if not fb2_to_markdown(input_file, temp_md):
        sys.exit(1)
    
    # Translate first 100 paragraphs for testing
    with open(temp_md, 'r', encoding='utf-8') as f:
        content = f.read()
    
    paragraphs = [p for p in content.split('\n\n') if p.strip()]
    limited_content = '\n\n'.join(paragraphs[:100])  # First 100 paragraphs only
    
    test_md = temp_md.replace('_temp.md', '_limited.md')
    with open(test_md, 'w', encoding='utf-8') as f:
        f.write(limited_content)
    
    print(f"Translating limited version with {len(paragraphs[:100])} paragraphs")
    if translate_markdown_file(test_md, output_md):
        # Create EPUB
        create_simple_epub(output_md, output_epub)
        
        print(f"\nTranslation completed!")
        print(f"Files created:")
        print(f"  Original FB2: {input_file}")
        print(f"  Original MD: {temp_md}")
        print(f"  Limited MD: {test_md}")
        print(f"  Translated MD: {output_md}")
        print(f"  Final EPUB: {output_epub}")
    else:
        print("Translation failed")
        sys.exit(1)

PYEOF

chmod +x final_book_translate.py

echo
echo "2. Running final book translation (limited to first 100 paragraphs for testing)..."
python3 final_book_translate.py "materials/books/book1.fb2" "materials/books/book1_translated.md" "materials/books/book1_sr.epub"

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
    echo "First 500 characters:"
    head -c 500 "materials/books/book1_translated.md"
    echo ""
    
    # Check for Cyrillic characters
    if grep -q '[\u0400-\u04FF]' "materials/books/book1_translated.md"; then
        echo "✓ Translation contains Cyrillic characters"
    else
        echo "✗ Translation does not contain Cyrillic characters"
    fi
    
    # Line count
    echo "Number of lines: $(wc -l < materials/books/book1_translated.md)"
else
    echo "✗ Translated markdown not created"
fi

EOF