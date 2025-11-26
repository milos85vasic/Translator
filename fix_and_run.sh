#!/bin/bash

echo "=== Fix and Run Translation ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Fixing missing import..."
sed -i '1i import sys' final_book_translate.py

echo
echo "2. Running fixed translation..."
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