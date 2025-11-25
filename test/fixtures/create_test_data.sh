#!/bin/bash

# Create test data directory
mkdir -p test/fixtures/ebooks
mkdir -p test/fixtures/translations
mkdir -p test/fixtures/configs

# Create test ebooks
cat > test/fixtures/ebooks/sample.txt << 'EOF'
This is a sample English text for testing translation functionality.
It contains multiple sentences and should be suitable for translation testing.
The text includes various punctuation and sentence structures to test comprehensive translation.
EOF

cat > test/fixtures/ebooks/sample.html << 'EOF'
<!DOCTYPE html>
<html>
<head><title>Sample Document</title></head>
<body>
<h1>Test Document</h1>
<p>This is a test HTML document for translation.</p>
<p>It contains headings and paragraphs for comprehensive testing.</p>
<ul>
<li>First item in list</li>
<li>Second item in list</li>
<li>Third item in list</li>
</ul>
</body>
</html>
EOF

cat > test/fixtures/ebooks/sample.fb2 << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<FictionBook xmlns="http://www.gribuser.ru/xml/fictionbook/2.0">
<description>
<title-info>
<book-title>Test Book</book-title>
<author>
<first-name>Test</first-name>
<last-name>Author</last-name>
</author>
<lang>en</lang>
</title-info>
</description>
<body>
<section>
<title><p>Test Chapter</p></title>
<p>This is a test FB2 document for translation testing.</p>
<p>It contains section titles and paragraphs to validate conversion accuracy.</p>
<subtitle><p>Section 1.1</p></subtitle>
<p>This is a subsection with additional content for comprehensive testing.</p>
</section>
</body>
</FictionBook>
EOF

# Create multi-language test files
cat > test/fixtures/ebooks/spanish_sample.txt << 'EOF'
Este es un texto de muestra en español para probar la funcionalidad de traducción.
Contiene múltiples oraciones y debe ser adecuado para las pruebas de traducción.
El texto incluye varios signos de puntuación y estructuras de oraciones.
EOF

cat > test/fixtures/ebooks/german_sample.txt << 'EOF'
Dies ist ein deutscher Beispieltext zum Testen der Übersetzungsfunktionalität.
Er enthält mehrere Sätze und eignet sich für Übersetzungstests.
Der Text enthält verschiedene Satzzeichen und Satzstrukturen.
EOF

cat > test/fixtures/ebooks/russian_sample.txt << 'EOF'
Это образец русского текста для проверки функции перевода.
Он содержит несколько предложений и подходит для тестирования перевода.
Текст включает различные знаки препинания и структуры предложений.
EOF

# Create test configurations
cat > test/fixtures/configs/test_config.json << 'EOF'
{
  "provider": "mock",
  "model": "test-model",
  "temperature": 0.3,
  "max_tokens": 1000,
  "target_language": "es",
  "output_format": "epub"
}
EOF

cat > test/fixtures/configs/performance_test.json << 'EOF'
{
  "provider": "mock",
  "model": "performance-test",
  "temperature": 0.1,
  "max_tokens": 2000,
  "target_language": "fr",
  "output_format": "txt"
}
EOF

cat > test/fixtures/configs/security_test.json << 'EOF'
{
  "provider": "openai",
  "model": "gpt-4",
  "temperature": 0.2,
  "max_tokens": 1500,
  "target_language": "de",
  "output_format": "fb2",
  "api_key": "test-key-for-security-testing"
}
EOF

# Create expected translation results for testing
cat > test/fixtures/translations/expected_translations.json << 'EOF'
{
  "translations": {
    "en_to_es": {
      "source": "Hello world",
      "expected": "Hola mundo"
    },
    "en_to_fr": {
      "source": "Hello world",
      "expected": "Bonjour le monde"
    },
    "en_to_de": {
      "source": "Hello world",
      "expected": "Hallo Welt"
    },
    "en_to_ru": {
      "source": "Hello world",
      "expected": "Привет мир"
    },
    "es_to_en": {
      "source": "Hola mundo",
      "expected": "Hello world"
    }
  }
}
EOF

# Create test data for batch processing
cat > test/fixtures/ebooks/batch_test_1.txt << 'EOF'
First document for batch processing test.
This document contains multiple sentences.
The purpose is to test batch translation functionality.
EOF

cat > test/fixtures/ebooks/batch_test_2.txt << 'EOF'
Second document for batch processing test.
This is another file in the batch.
Batch processing should handle multiple files efficiently.
EOF

cat > test/fixtures/ebooks/batch_test_3.txt << 'EOF'
Third document for batch processing test.
This completes our batch test set.
All files should be processed correctly.
EOF

echo "Test data created successfully"
echo "Files created in test/fixtures/:"
find test/fixtures -type f | sort