#!/bin/bash

echo "Testing SSH translation with simple markdown file..."

# Create SSH connection and test translation directly
ssh milosvasic@thinker.local << 'EOF'
cd /tmp/translate-ssh

# Check if files exist
ls -la

# Check if llama.cpp works
if [ -f "./llama.cpp" ]; then
    chmod +x ./llama.cpp
    echo "Testing llama.cpp binary..."
    ./llama.cpp --help | head -5
fi

# Check if translation script exists
if [ -f "./translate_llm_only.py" ]; then
    echo "Testing translation script..."
    python3 translate_llm_only.py "test_simple.md" "test_simple_translated.md"
fi

EOF