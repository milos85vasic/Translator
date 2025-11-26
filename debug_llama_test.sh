#!/bin/bash

echo "Testing llama.cpp availability on remote worker..."

# Check if llama.cpp binary exists
echo "1. Checking llama.cpp binary location..."
ssh milosvasic@thinker.local "find /home/milosvasic/llama.cpp -name 'llama*' -type f 2>/dev/null | head -5"

# Check if models exist
echo "2. Checking for GGUF models..."
ssh milosvasic@thinker.local "find /home/milosvasic/models -name '*.gguf' 2>/dev/null | head -5"

# Test simple translation
echo "3. Testing simple translation command..."
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && timeout 30 python3 -c \"
import sys
sys.path.append('internal/scripts')
from translate_llm_only import find_llama_binary, has_llama_model
print('Llama binary:', find_llama_binary())
print('Has model:', has_llama_model())
\""

echo "4. Testing with very small text sample..."
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && echo 'Это тест.' > tiny_test.md && timeout 60 python3 internal/scripts/translate_llm_only.py tiny_test.md tiny_test_translated.md && echo 'Translation completed' && head -3 tiny_test_translated.md || echo 'Translation failed or timed out'"