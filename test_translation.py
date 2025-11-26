#!/usr/bin/env python3
"""
Simple test to check if translation infrastructure works
"""

import sys
import os

def main():
    if len(sys.argv) != 3:
        print("Usage: python3 test_translation.py <input> <output>")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    print(f"Testing translation from {input_file} to {output_file}")
    
    # Read input file
    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            content = f.read()
        print(f"Read {len(content)} characters from input file")
    except Exception as e:
        print(f"Error reading input file: {e}")
        sys.exit(1)
    
    # Simple mock translation - just add Serbian Cyrillic chars
    mock_translation = content.replace('Hello', 'Здраво').replace('The', 'У').replace('and', 'и')
    
    # Write output file
    try:
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(mock_translation)
        print(f"Wrote {len(mock_translation)} characters to output file")
    except Exception as e:
        print(f"Error writing output file: {e}")
        sys.exit(1)
    
    print("Test translation completed successfully!")

if __name__ == "__main__":
    main()