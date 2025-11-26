#!/bin/bash

# Check status and upload our sample file
echo "=== Checking SSH Worker Status ==="
./translate-ssh status --input /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sample_original.md

echo -e "\n=== Uploading sample file ==="
./translate-ssh upload \
  --source /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sample_original.md \
  --target book1_sample_original.md

echo -e "\n=== Translating sample ==="
./translate-ssh translate \
  --input book1_sample_original.md \
  --output book1_sample_translated.md \
  --engine "llama.cpp" \
  --model "tiny-llama-working.gguf" \
  --prompt "Translate this Russian text to Serbian Cyrillic. Return only the Serbian translation without any explanations." \
  --language "sr-Cyrl" \
  --monitor

echo -e "\n=== Downloading translation ==="
./translate-ssh download \
  --source book1_sample_translated.md \
  --target /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sample_translated.md

echo -e "\n=== Displaying translation ==="
cat /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sample_translated.md