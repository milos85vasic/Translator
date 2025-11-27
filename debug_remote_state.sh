#!/bin/bash

echo "Checking what's happening on remote worker during translation..."

# 1. Check if test file exists and its content
echo "1. Checking test file content:"
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && if [ -f book1_test.md ]; then echo 'File exists, content:'; head -10 book1_test.md; else echo 'Test file does not exist'; fi"

# 2. Try to run the translation script manually with verbose output
echo "2. Running translation script with verbose output:"
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && PYTHONPATH=/tmp/translate-workspace/internal/scripts python3 -v internal/scripts/translate_llm_only.py book1_test.md book1_test_translated.md 2>&1 | head -20"

# 3. Check if Python can import the script
echo "3. Testing Python import:"
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && python3 -c \"import sys; sys.path.append('internal/scripts'); import translate_llm_only; print('Script imported successfully')\""

# 4. Check process status
echo "4. Checking if any translation processes are running:"
ssh milosvasic@thinker.local "ps aux | grep -E '(python|llama)' | grep -v grep || echo 'No processes found'"