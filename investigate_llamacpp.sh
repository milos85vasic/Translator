#!/bin/bash

echo "=== Investigating llama.cpp installation ==="
echo

ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null milosvasic@thinker.local << 'EOF'
echo
echo "1. Checking llama.cpp directory structure..."
ls -la /home/milosvasic/llama.cpp/

echo
echo "2. Looking for llama.cpp binaries in common locations..."
find /home/milosvasic/llama.cpp -name "main" -o -name "llama" -o -name "llama-cli" -o -name "llama.cpp" -type f 2>/dev/null | head -10

echo
echo "3. Checking if there's a build directory..."
ls -la /home/milosvasic/llama.cpp/build/ 2>/dev/null || echo "No build directory"

echo
echo "4. Looking for compiled binaries in build..."
find /home/milosvasic/llama.cpp/build -name "main" -o -name "llama" -o -name "llama-cli" -type f 2>/dev/null

echo
echo "5. If no binary found, trying to build llama.cpp..."
if ! find /home/milosvasic/llama.cpp -name "main" -type f 2>/dev/null | grep -q .; then
    echo "No llama.cpp binary found, attempting to build..."
    cd /home/milosvasic/llama.cpp
    
    # Check if CMakeLists.txt exists
    if [ -f "CMakeLists.txt" ]; then
        echo "Found CMakeLists.txt, attempting to build with cmake..."
        mkdir -p build
        cd build
        cmake .. -DLLAMA_BLAS=OFF -DLLAMA_CUBLAS=OFF
        make -j2 llama-cli
        
        echo "Build completed. Checking for binary..."
        ls -la ./bin/llama-cli 2>/dev/null || ls -la ./main 2>/dev/null || echo "Build failed or binary not found"
    else
        echo "No CMakeLists.txt found, checking for Makefile..."
        if [ -f "Makefile" ]; then
            make -j2
        else
            echo "No build system found"
        fi
    fi
fi

echo
echo "6. Final check for usable llama.cpp binary..."
find /home/milosvasic/llama.cpp -name "main" -o -name "llama-cli" -type f -executable 2>/dev/null

EOF