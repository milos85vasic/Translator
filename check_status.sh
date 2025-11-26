#!/bin/bash

echo "=== Checking Current Translation Status ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

echo
echo "1. Checking running processes..."
ps aux | grep -E "(python|translate)" | grep -v grep

echo
echo "2. Checking sample files..."
ls -la materials/books/book1_sample* 2>/dev/null || echo "No sample files found"

echo
echo "3. Checking if any translations completed..."
if [ -f "materials/books/book1_sample_translated.md" ]; then
    echo "Sample translation found!"
    echo "Content:"
    head -50 "materials/books/book1_sample_translated.md"
else
    echo "No sample translation yet"
fi

echo
echo "4. System load..."
uptime

EOF