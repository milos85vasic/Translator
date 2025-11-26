#!/bin/bash

echo "=== Testing with a Smaller Sample First ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Killing the current translation process..."
pkill -f "complete_translate.py"

echo
echo "2. Creating a smaller sample of the book..."
# Take first 50 lines of the original markdown
head -50 "materials/books/book1_original.md" > "materials/books/book1_sample.md"

echo "Sample created with $(wc -l < materials/books/book1_sample.md) lines"
echo "Sample content:"
cat "materials/books/book1_sample.md"

echo
echo "3. Translating the sample..."
python3 complete_translate.py "materials/books/book1.fb2" "materials/books/book1_sample_translated.md" "materials/books/book1_sample.epub"

echo
echo "4. Checking sample translation results..."
if [ -f "materials/books/book1_sample_translated.md" ]; then
    echo "✓ Sample translation successful!"
    echo "File size: $(du -h materials/books/book1_sample_translated.md | cut -f1)"
    echo "Content:"
    cat "materials/books/book1_sample_translated.md"
    
    # Check for Cyrillic characters
    if grep -q '[\u0400-\u04FF]' "materials/books/book1_sample_translated.md"; then
        echo "✓ Translation contains Cyrillic characters"
    else
        echo "✗ Translation does not contain Cyrillic characters"
    fi
else
    echo "✗ Sample translation failed"
fi

if [ -f "materials/books/book1_sample.epub" ]; then
    echo "✓ Sample EPUB created!"
    file "materials/books/book1_sample.epub"
else
    echo "✗ Sample EPUB not created"
fi

EOF