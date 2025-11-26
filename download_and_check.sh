#!/bin/bash

echo "=== Download and Check Results ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Checking what files we have..."
ls -la materials/books/book1*

echo
echo "2. Let's download the files to local machine to check..."

EOF

echo
echo "3. Downloading files from remote host..."
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local:/tmp/translate-ssh/materials/books/book1_original.md ./internal/materials/books/
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local:/tmp/translate-ssh/materials/books/book1_temp.md ./internal/materials/books/

echo
echo "4. Checking downloaded files..."
ls -la internal/materials/books/book1*

echo
echo "5. Preview of original markdown (first 50 lines):"
head -50 internal/materials/books/book1_original.md

echo
echo "6. Summary of progress so far:"
echo "✓ SSH connection to thinker.local works"
echo "✓ File upload to remote works"
echo "✓ FB2 to Markdown conversion works"
echo "✓ llama.cpp binary detection works"
echo "✓ llama.cpp model detection works"
echo "✓ Single paragraph translation to Cyrillic works"
echo "✗ Full book translation still needs debugging"

EOF