#!/usr/bin/env python3
"""
Test script to verify translation works with llama.cpp on remote worker
"""

import sys
import os

# Add the scripts path to import translation functions
sys.path.append('internal/scripts')

try:
    from translate_llm_only import get_translation_provider, translate_text
    
    print("Testing translation provider availability...")
    provider, config = get_translation_provider()
    print(f"Available provider: {provider}")
    
    if provider:
        test_text = "Это тестовый текст на русском языке."
        print(f"Translating: {test_text}")
        
        translated = translate_text(test_text, "ru", "sr")
        print(f"Translation result: {translated}")
        
        if translated and translated != test_text:
            print("✓ Translation successful!")
        else:
            print("✗ Translation failed - no change detected")
    else:
        print("✗ No translation provider available")
        
except Exception as e:
    print(f"Error during translation test: {e}")
    import traceback
    traceback.print_exc()