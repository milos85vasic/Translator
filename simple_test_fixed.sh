#!/bin/bash

HOST="thinker.local"
USER="milosvasic"
REMOTE_DIR="/tmp/translate-ssh"

echo "=== Creating remote directory ==="
ssh $USER@$HOST "mkdir -p $REMOTE_DIR"

echo -e "\n=== Uploading sample file ==="
scp /Users/milosvasic/Projects/Translate/internal/materials/books/book1_sample_original.md $USER@$HOST:$REMOTE_DIR/

echo -e "\n=== Testing translation command directly on remote ==="
ssh $USER@$HOST "cd $REMOTE_DIR && /home/milosvasic/llama.cpp/build/bin/llama-cli -m /home/milosvasic/models/tiny-llama-working.gguf --prompt 'Translate this Russian text to Serbian Cyrillic: Снег танцевал в свете фонаря, подобно хлопковому пуху.' -n 100 --temp 0.7 --ctx-size 2048"