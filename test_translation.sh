#!/bin/bash

# Test translation with llama.cpp on remote worker
echo "Testing translation on remote worker..."

# Upload test script to remote
scp test_translation_remote.py milosvasic@thinker.local:/tmp/translate-workspace/

# Execute test script remotely
ssh milosvasic@thinker.local "cd /tmp/translate-workspace && python3 test_translation_remote.py"

# Clean up
ssh milosvasic@thinker.local "rm -f /tmp/translate-workspace/test_translation_remote.py"