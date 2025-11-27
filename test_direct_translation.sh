#!/bin/bash

echo "Testing direct translation workflow execution..."

# Upload a simple test file
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && echo -e '# Тестовая глава\\n\\nЭто тестовый абзац на русском языке. Мы хотим проверить перевод.\\n\\nВторой абзац для проверки.' > test_direct.md"

echo "1. Test file created"
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && cat test_direct.md"

echo "2. Running translation with timeout and error output..."
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && timeout 120 python3 internal/scripts/translate_llm_only.py test_direct.md test_direct_translated.md 2>&1; echo 'Exit code: \$?'"

echo "3. Checking if translation completed..."
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && if [ -f test_direct_translated.md ]; then echo 'Translation file exists:'; head -10 test_direct_translated.md; else echo 'Translation file NOT created'; fi"

echo "4. Cleaning up..."
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && rm -f test_direct.md test_direct_translated.md"