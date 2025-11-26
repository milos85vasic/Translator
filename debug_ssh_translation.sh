#!/bin/bash

echo "=== Debugging SSH Translation Process ==="
echo

# First upload our test FB2 file
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null internal/materials/books/test_book_small.fb2 milosvasic@thinker.local:/tmp/test_book_small.fb2

echo
echo "1. Testing SSH connection and file upload..."
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
echo "Connected to remote host successfully"
ls -la /tmp/test_book_small.fb2

# Create working directory
mkdir -p /tmp/translate-ssh/materials/books
cp /tmp/test_book_small.fb2 /tmp/translate-ssh/materials/books/
cd /tmp/translate-ssh

echo
echo "2. Testing FB2 to Markdown conversion..."
# Check if Python3 is available
python3 --version

# Create a simple FB2 to markdown converter
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

# Run conversion
echo "Running conversion..."
python3 fb2_to_markdown.py "materials/books/test_book_small.fb2" "materials/books/test_book_small_original.md"

echo
echo "3. Checking conversion results..."
if [ -f "materials/books/test_book_small_original.md" ]; then
    echo "Markdown file created successfully:"
    cat "materials/books/test_book_small_original.md"
else
    echo "Failed to create markdown file"
fi

EOF

echo
echo "4. Checking results on local machine..."
echo "If this completes successfully, we'll know the issue is in the translation step only."