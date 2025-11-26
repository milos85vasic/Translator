#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"

echo "=== Testing with conversation format ==="
ssh $USER@$HOST "
cd $REMOTE_DIR

# Create a simple test with system prompt
cat > test_prompt.txt << 'PROMPT'
<|im_start|>system
You are a translator. Translate the given text to Serbian Cyrillic. Return only the translation.<|im_end|>
<|im_start|>user
Translate this Russian text to Serbian Cyrillic: Снег танцевал в свете фонаря<|im_end|>
<|im_start|>assistant
PROMPT

/home/milosvasic/llama.cpp/build/bin/llama-cli \
  -m /home/milosvasic/models/tiny-llama-working.gguf \
  --file test_prompt.txt \
  -n 50 \
  --temp 0.3 \
  --no-display-prompt \
  2>&1 | tail -10
"