#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"

echo "=== Testing llama.cpp directly ==="
ssh $USER@$HOST "
cd $REMOTE_DIR

# Test 1: Simple translation
echo 'Test 1: Basic translation'
/home/milosvasic/llama.cpp/build/bin/llama-cli \
  -m /home/milosvasic/models/tiny-llama-working.gguf \
  --prompt 'Translate to Serbian: Снег танцевал' \
  -n 30 \
  --temp 0.3 \
  2>&1 | grep -E '[А-ЯЉЊЋЂЏ]' | tail -3

echo ''
echo 'Test 2: Serbian prompt'
/home/milosvasic/llama.cpp/build/bin/llama-cli \
  -m /home/milosvasic/models/tiny-llama-working.gguf \
  --prompt 'Преведи на српски: Снег танцевал' \
  -n 30 \
  --temp 0.3 \
  2>&1 | grep -E '[А-ЯЉЊЋЂЏ]' | tail -3
"