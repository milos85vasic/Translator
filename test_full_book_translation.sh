#!/bin/bash

echo "=== Testing Full Book Translation with Fixed Script ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Uploading book1.fb2 file..."
# First convert book1.fb2 to markdown on the remote
cat > fb2_to_markdown.py << 'PYEOF'
#!/usr/bin/env python3
import sys
import xml.etree.ElementTree as ET

def fb2_to_markdown(input_file, output_file):
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

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 fb2_to_markdown.py <input.fb2> <output.md>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    success = fb2_to_markdown(input_file, output_file)
    sys.exit(0 if success else 1)
PYEOF

chmod +x fb2_to_markdown.py

# Upload the FB2 file from local to remote
echo "Uploading book1.fb2..."
echo "Note: You should upload the file manually before running this script"
echo "Assuming book1.fb2 is already uploaded..."

# Convert to markdown
echo "Converting to markdown..."
python3 fb2_to_markdown.py "materials/books/book1.fb2" "materials/books/book1_original.md"

echo
echo "2. Checking conversion results..."
if [ -f "materials/books/book1_original.md" ]; then
    echo "Markdown created successfully"
    echo "File size: $(wc -c < materials/books/book1_original.md) bytes"
    echo "Number of lines: $(wc -l < materials/books/book1_original.md) lines"
    echo "First 200 characters:"
    head -c 200 "materials/books/book1_original.md"
    echo ""
else
    echo "Failed to create markdown"
    exit 1
fi

echo
echo "3. Running translation with fixed script..."
# Use the final translate script we created earlier
python3 final_translate.py "materials/books/book1_original.md" "materials/books/book1_translated.md"

echo
echo "4. Checking translation results..."
if [ -f "materials/books/book1_translated.md" ]; then
    echo "Translation successful!"
    echo "File size: $(wc -c < materials/books/book1_translated.md) bytes"
    echo "Number of lines: $(wc -l < materials/books/book1_translated.md) lines"
    echo "First 200 characters:"
    head -c 200 "materials/books/book1_translated.md"
    echo ""
    
    # Check for Cyrillic characters
    if grep -q '[\u0400-\u04FF]' "materials/books/book1_translated.md"; then
        echo "✓ Translation contains Cyrillic characters"
    else
        echo "✗ Translation does not contain Cyrillic characters"
    fi
else
    echo "Translation failed"
fi

EOF