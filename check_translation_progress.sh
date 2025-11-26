#!/bin/bash

echo "=== Checking Translation Progress ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Checking if translation script is running..."
ps aux | grep -E "(python3|complete_translate)" | grep -v grep

echo
echo "2. Checking files created so far..."
ls -la materials/books/book1* 2>/dev/null || echo "No files found yet"

echo
echo "3. Checking if any partial translation exists..."
if [ -f "materials/books/book1_translated.md" ]; then
    echo "Partial translation exists!"
    echo "File size: $(wc -l < materials/books/book1_translated.md) lines"
    echo "Last few lines:"
    tail -5 "materials/books/book1_translated.md"
else
    echo "No translation file yet"
fi

echo
echo "4. Checking system resources..."
echo "Memory usage:"
free -h
echo
echo "CPU load:"
uptime

EOF