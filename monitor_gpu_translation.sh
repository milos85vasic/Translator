#!/bin/bash
echo "=== GPU Translation Progress Monitor ==="
echo "Checking translation status..."

ssh milosvasic@thinker.local "cd /tmp/translate-ssh && \
echo '=== Process Status ===' && \
ps aux | grep 'gpu_translate.py' | grep -v grep || echo 'No translation process running' && \
echo && \
echo '=== Recent Logs ===' && \
tail -20 gpu_translation.log && \
echo && \
echo '=== File Status ===' && \
ls -la book1_gpu_translated.md 2>/dev/null || echo 'Output file not created yet' && \
if [ -f book1_gpu_translated.md ]; then \
  echo 'Output file size:' && \
  wc -c book1_gpu_translated.md && \
  echo 'Last few lines:' && \
  tail -5 book1_gpu_translated.md; \
fi"