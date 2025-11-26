package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"digital.vasic.translator/pkg/fb2"
	"digital.vasic.translator/pkg/logger"
	"digital.vasic.translator/pkg/markdown"
	"digital.vasic.translator/pkg/sshworker"
	"digital.vasic.translator/pkg/version"
)

// TestEBookTranslationWorkflow is a comprehensive test suite for the ebook translation workflow
func TestEBookTranslationWorkflow(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "ebook_translation_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test logger
	testLogger := logger.NewLogger(logger.Config{
		Level:  "debug",
		Format: "text",
	})

	// Test 1: Create test FB2 file
	fb2Path := filepath.Join(tempDir, "test_book.fb2")
	if err := createTestFB2File(fb2Path); err != nil {
		t.Fatalf("Failed to create test FB2 file: %v", err)
	}

	// Test 2: FB2 to Markdown conversion
	t.Run("FB2ToMarkdown", func(t *testing.T) {
		testFB2ToMarkdownConversion(t, testLogger, fb2Path, tempDir)
	})

	// Test 3: Markdown to EPUB conversion
	t.Run("MarkdownToEPUB", func(t *testing.T) {
		mdPath := filepath.Join(tempDir, "test_book_original.md")
		testMarkdownToEPUBConversion(t, testLogger, mdPath, tempDir)
	})

	// Test 4: Codebase version hashing
	t.Run("CodebaseHashing", func(t *testing.T) {
		testCodebaseHashing(t, tempDir)
	})

	// Test 5: File verification
	t.Run("FileVerification", func(t *testing.T) {
		testFileVerification(t, tempDir)
	})
}

// TestSSHWorkerConnection tests SSH worker connectivity and command execution
func TestSSHWorkerConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping SSH test in short mode")
	}

	// These would be your actual SSH connection parameters
	// For testing, we'll use environment variables to avoid hardcoding credentials
	host := os.Getenv("TEST_SSH_HOST")
	user := os.Getenv("TEST_SSH_USER")
	pass := os.Getenv("TEST_SSH_PASS")

	if host == "" || user == "" || pass == "" {
		t.Skip("SSH credentials not provided in environment variables")
	}

	testLogger := logger.NewLogger(logger.Config{
		Level:  "debug",
		Format: "text",
	})

	sshConfig := sshworker.SSHWorkerConfig{
		Host:              host,
		Username:          user,
		Password:          pass,
		Port:              22,
		RemoteDir:         "/tmp",
		ConnectionTimeout: 30 * time.Second,
		CommandTimeout:    30 * time.Second,
	}

	worker, err := sshworker.NewSSHWorker(sshConfig, testLogger)
	if err != nil {
		t.Fatalf("Failed to create SSH worker: %v", err)
	}

	ctx := context.Background()

	// Test connection
	if err := worker.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect to SSH worker: %v", err)
	}
	defer worker.Disconnect()

	// Test command execution
	output, err := worker.ExecuteCommand(ctx, "echo 'SSH test successful'", 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to execute SSH command: %v", err)
	}

	expected := "SSH test successful"
	if !strings.Contains(output, expected) {
		t.Errorf("SSH command output mismatch. Expected: %s, Got: %s", expected, output)
	}

	// Test file transfer
	testFile := filepath.Join(os.TempDir(), "ssh_test.txt")
	testContent := "This is a test file for SSH transfer"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	remotePath := "/tmp/ssh_test.txt"
	if err := worker.TransferFile(ctx, testFile, remotePath); err != nil {
		t.Fatalf("Failed to transfer file: %v", err)
	}

	// Verify file exists on remote
	output, err = worker.ExecuteCommand(ctx, fmt.Sprintf("cat %s", remotePath), 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to read transferred file: %v", err)
	}

	if !strings.Contains(output, testContent) {
		t.Errorf("Transferred file content mismatch. Expected: %s, Got: %s", testContent, output)
	}

	// Clean up remote file
	_, _ = worker.ExecuteCommand(ctx, fmt.Sprintf("rm -f %s", remotePath), 5*time.Second)
}

// TestCodebaseVersionConsistency tests codebase version verification between local and remote
func TestCodebaseVersionConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping codebase version test in short mode")
	}

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "codebase_version_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test codebase files
	testFiles := map[string]string{
		"main.go":           "package main\n\nfunc main() { println(\"test\") }",
		"utils.go":          "package main\n\nfunc test() { println(\"utils\") }",
		"config.json":       `{"test": "value"}`,
		"README.md":         "# Test Project\n\nThis is a test project",
		"scripts/build.sh":  "#!/bin/bash\necho Building...",
	}

	for file, content := range testFiles {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	// Test codebase hashing
	hasher := version.NewCodebaseHasher()
	hash1, err := hasher.CalculateHash()
	if err != nil {
		t.Fatalf("Failed to calculate codebase hash: %v", err)
	}

	// Wait a moment to ensure different timestamp
	time.Sleep(100 * time.Millisecond)

	hash2, err := hasher.CalculateHash()
	if err != nil {
		t.Fatalf("Failed to calculate second codebase hash: %v", err)
	}

	// Hashes should be the same for unchanged content
	if hash1 != hash2 {
		t.Errorf("Codebase hashes differ for unchanged content: %s vs %s", hash1, hash2)
	}

	// Modify a file
	testFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc main() { println(\"modified\") }"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	hash3, err := hasher.CalculateHash()
	if err != nil {
		t.Fatalf("Failed to calculate third codebase hash: %v", err)
	}

	// Hash should change after file modification
	if hash2 == hash3 {
		t.Errorf("Codebase hash should have changed after file modification: %s", hash3)
	}

	t.Logf("Codebase hash test passed: %s -> %s -> %s", hash1, hash2, hash3)
}

// TestCompleteTranslationWorkflow tests the complete translation workflow
func TestCompleteTranslationWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping complete workflow test in short mode")
	}

	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "complete_workflow_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test FB2 file
	fb2Path := filepath.Join(tempDir, "test_book.fb2")
	if err := createTestFB2File(fb2Path); err != nil {
		t.Fatalf("Failed to create test FB2 file: %v", err)
	}

	// Test logger
	testLogger := logger.NewLogger(logger.Config{
		Level:  "info",
		Format: "text",
	})

	// Create workflow with dummy SSH parameters (will fail connection tests but allows testing other parts)
	workflow, err := NewEBookTranslationWorkflow(fb2Path, "sr-cyrl", "dummy.host", "dummy_user", "dummy_pass")
	if err != nil {
		t.Fatalf("Failed to create workflow: %v", err)
	}

	// Test individual components that don't require SSH connection

	// Test FB2 to Markdown conversion
	originalMdPath, err := workflow.convertFB2ToMarkdown()
	if err != nil {
		t.Errorf("FB2 to Markdown conversion failed: %v", err)
	}

	// Verify original markdown file
	if err := workflow.verifyMarkdownFile(originalMdPath, "original"); err != nil {
		t.Errorf("Original markdown verification failed: %v", err)
	}

	// Create a dummy translated markdown file for testing
	translatedMdPath := strings.TrimSuffix(originalMdPath, filepath.Ext(originalMdPath)) + "_translated.md"
	dummyTranslatedContent := "# Тест книга\n\nОво је преведени садржај на ћирилицу.\n\nОво је тестни пасус који треба да садржи ћириличка слова као што су ћ, ћ, ч, ћ, ш, ж, љ, њ, и други."
	if err := os.WriteFile(translatedMdPath, []byte(dummyTranslatedContent), 0644); err != nil {
		t.Fatalf("Failed to create dummy translated markdown: %v", err)
	}

	// Verify translated markdown file
	if err := workflow.verifyMarkdownFile(translatedMdPath, "translated"); err != nil {
		t.Errorf("Translated markdown verification failed: %v", err)
	}

	// Verify Cyrillic content
	if err := workflow.verifyLanguage(translatedMdPath); err != nil {
		t.Errorf("Cyrillic verification failed: %v", err)
	}

	// Test Markdown to EPUB conversion
	finalEpubPath, err := workflow.convertMarkdownToEPUB(translatedMdPath)
	if err != nil {
		t.Errorf("Markdown to EPUB conversion failed: %v", err)
	}

	// Verify final EPUB file
	if err := workflow.verifyEPUBFile(finalEpubPath); err != nil {
		t.Errorf("EPUB verification failed: %v", err)
	}

	// Test final verification of all outputs
	if err := workflow.verifyFinalOutput(finalEpubPath, translatedMdPath, originalMdPath); err != nil {
		t.Errorf("Final output verification failed: %v", err)
	}

	t.Logf("Complete workflow test passed with files: %s, %s, %s", originalMdPath, translatedMdPath, finalEpubPath)
}

// createTestFB2File creates a minimal test FB2 file for testing
func createTestFB2File(path string) error {
	fb2Content := `<?xml version="1.0" encoding="UTF-8"?>
<FictionBook xmlns="http://www.gribuser.ru/xml/fictionbook/2.0" xmlns:l="http://www.w3.org/1999/xlink">
  <description>
    <title-info>
      <genre>fiction</genre>
      <author>
        <first-name>Test</first-name>
        <last-name>Author</last-name>
      </author>
      <book-title>Test Book</book-title>
      <annotation>
        <p>This is a test book for translation workflow testing.</p>
      </annotation>
      <lang>en</lang>
    </title-info>
  </description>
  <body>
    <title>
      <p>Chapter 1</p>
    </title>
    <section>
      <title>
        <p>Introduction</p>
      </title>
      <p>This is the first paragraph of our test book. It contains simple English text that should be translated to Serbian Cyrillic.</p>
      <p>This is the second paragraph. It continues with more test content to ensure proper translation and formatting.</p>
      <subtitle>A Subtitle</subtitle>
      <p>This paragraph follows a subtitle to test proper markdown formatting.</p>
      <epigraph>
        <p>This is a test epigraph that should be formatted as a blockquote.</p>
        <text-author>Test Quote Author</text-author>
      </epigraph>
      <p>Regular text continues after the epigraph.</p>
    </section>
    <section>
      <title>
        <p>Chapter 2</p>
      </title>
      <p>This is the second chapter with more content to translate.</p>
      <poem>
        <title>
          <p>Test Poem</p>
        </title>
        <stanza>
          <v>This is the first verse of a test poem.</v>
          <v>This is the second verse of the test poem.</v>
        </stanza>
      </poem>
      <p>Text continues after the poem.</p>
      <cite>
        <p>This is a test citation that should be formatted as a blockquote.</p>
        <text-author>Test Citation Author</text-author>
      </cite>
      <p>Final paragraph of the test book.</p>
    </section>
  </body>
</FictionBook>`

	return os.WriteFile(path, []byte(fb2Content), 0644)
}

// testFB2ToMarkdownConversion tests the FB2 to Markdown conversion
func testFB2ToMarkdownConversion(t *testing.T, testLogger logger.Logger, fb2Path, tempDir string) {
	converter := fb2.NewMarkdownConverter(testLogger)
	mdPath := filepath.Join(tempDir, "test_book.md")

	if err := converter.ConvertToMarkdown(fb2Path, mdPath); err != nil {
		t.Fatalf("FB2 to Markdown conversion failed: %v", err)
	}

	// Verify markdown file was created
	if _, err := os.Stat(mdPath); err != nil {
		t.Fatalf("Markdown file was not created: %v", err)
	}

	// Verify markdown content
	content, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("Failed to read markdown file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "# Test Book") {
		t.Errorf("Markdown content missing title: %s", contentStr)
	}

	if !strings.Contains(contentStr, "## Authors") {
		t.Errorf("Markdown content missing authors: %s", contentStr)
	}

	if !strings.Contains(contentStr, "This is the first paragraph") {
		t.Errorf("Markdown content missing first paragraph: %s", contentStr)
	}

	t.Logf("FB2 to Markdown conversion test passed: %s", mdPath)
}

// testMarkdownToEPUBConversion tests the Markdown to EPUB conversion
func testMarkdownToEPUBConversion(t *testing.T, testLogger logger.Logger, mdPath, tempDir string) {
	// First ensure we have a markdown file
	if _, err := os.Stat(mdPath); err != nil {
		// Create a simple markdown file for testing
		simpleMD := `---
title: Test Book
author: Test Author
language: sr
---

# Test Book

This is a test chapter.

## Chapter 1

This is the first paragraph of the test book.

This is the second paragraph.

### Subtitle

This paragraph follows a subtitle.

> This is a test epigraph that should be formatted as a blockquote.
> 
> — Test Quote Author

Regular text continues after the epigraph.

## Chapter 2

This is the second chapter with more content.

    This is the first verse of a test poem.
    This is the second verse of the test poem.

Text continues after the poem.

> This is a test citation that should be formatted as a blockquote.
> 
> — Test Citation Author

Final paragraph of the test book.
`
		if err := os.WriteFile(mdPath, []byte(simpleMD), 0644); err != nil {
			t.Fatalf("Failed to create test markdown file: %v", err)
		}
	}

	converter := markdown.NewMarkdownToEPUBConverter()
	epubPath := filepath.Join(tempDir, "test_book.epub")

	if err := converter.ConvertMarkdownToEPUB(mdPath, epubPath); err != nil {
		t.Fatalf("Markdown to EPUB conversion failed: %v", err)
	}

	// Verify EPUB file was created
	if _, err := os.Stat(epubPath); err != nil {
		t.Fatalf("EPUB file was not created: %v", err)
	}

	// Verify EPUB file size (should be reasonable)
	info, err := os.Stat(epubPath)
	if err != nil {
		t.Fatalf("Failed to stat EPUB file: %v", err)
	}

	if info.Size() < 1000 {
		t.Errorf("EPUB file seems too small: %d bytes", info.Size())
	}

	t.Logf("Markdown to EPUB conversion test passed: %s", epubPath)
}

// testCodebaseHashing tests the codebase hashing functionality
func testCodebaseHashing(t *testing.T, tempDir string) {
	hasher := version.NewCodebaseHasher()
	
	// Create some test files
	testFiles := map[string]string{
		"test1.go":   "package main\n\nfunc main() {}",
		"test2.txt":  "This is a test file",
		"config.json": `{"test": "value"}`,
	}

	for file, content := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Calculate hash
	hash, err := hasher.CalculateHash()
	if err != nil {
		t.Fatalf("Failed to calculate codebase hash: %v", err)
	}

	// Hash should not be empty
	if hash == "" {
		t.Error("Codebase hash is empty")
	}

	// Hash should be consistent
	hash2, err := hasher.CalculateHash()
	if err != nil {
		t.Fatalf("Failed to calculate second codebase hash: %v", err)
	}

	if hash != hash2 {
		t.Errorf("Codebase hash is not consistent: %s vs %s", hash, hash2)
	}

	// Modify a file and verify hash changes
	testFile := filepath.Join(tempDir, "test1.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc main() { println(\"modified\") }"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	hash3, err := hasher.CalculateHash()
	if err != nil {
		t.Fatalf("Failed to calculate third codebase hash: %v", err)
	}

	if hash2 == hash3 {
		t.Errorf("Codebase hash should have changed after file modification")
	}

	t.Logf("Codebase hashing test passed: %s -> %s -> %s", hash, hash2, hash3)
}

// testFileVerification tests file verification functions
func testFileVerification(t *testing.T, tempDir string) {
	// Create test files
	emptyFile := filepath.Join(tempDir, "empty.txt")
	smallFile := filepath.Join(tempDir, "small.txt")
	validFile := filepath.Join(tempDir, "valid.txt")

	if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	if err := os.WriteFile(smallFile, []byte("small"), 0644); err != nil {
		t.Fatalf("Failed to create small file: %v", err)
	}

	validContent := "This is a valid test file with sufficient content for verification.\n"
	validContent += "It contains multiple lines and enough characters to pass validation.\n"
	validContent += "The content should be meaningful and properly formatted.\n"
	if err := os.WriteFile(validFile, []byte(validContent), 0644); err != nil {
		t.Fatalf("Failed to create valid file: %v", err)
	}

	// Create workflow instance for testing
	workflow := &EBookTranslationWorkflow{
		logger: logger.NewLogger(logger.Config{
			Level:  "debug",
			Format: "text",
		}),
	}

	// Test empty file verification
	if err := workflow.verifyMarkdownFile(emptyFile, "empty"); err == nil {
		t.Error("Empty file verification should have failed")
	}

	// Test small file verification
	if err := workflow.verifyMarkdownFile(smallFile, "small"); err == nil {
		t.Error("Small file verification should have failed")
	}

	// Test valid file verification
	if err := workflow.verifyMarkdownFile(validFile, "valid"); err != nil {
		t.Errorf("Valid file verification failed: %v", err)
	}

	// Test Cyrillic verification
	cyrillicFile := filepath.Join(tempDir, "cyrillic.txt")
	cyrillicContent := "Ово је текст на ћирилицу. Садржи слова као што су ћ, ћ, ч, ћ, ш, ж, љ, њ."
	if err := os.WriteFile(cyrillicFile, []byte(cyrillicContent), 0644); err != nil {
		t.Fatalf("Failed to create Cyrillic file: %v", err)
	}

	workflow.targetLanguage = "sr-cyrl"
	if err := workflow.verifyLanguage(cyrillicFile); err != nil {
		t.Errorf("Cyrillic verification failed: %v", err)
	}

	// Test non-Cyrillic file with Cyrillic target
	latinFile := filepath.Join(tempDir, "latin.txt")
	latinContent := "Ovo je tekst na latinici. Ne sadrzi cirilicna slova."
	if err := os.WriteFile(latinFile, []byte(latinContent), 0644); err != nil {
		t.Fatalf("Failed to create Latin file: %v", err)
	}

	if err := workflow.verifyLanguage(latinFile); err == nil {
		t.Error("Latin file should have failed Cyrillic verification")
	}

	t.Logf("File verification test passed")
}

// BenchmarkCodebaseHashing benchmarks the codebase hashing performance
func BenchmarkCodebaseHashing(b *testing.B) {
	hasher := version.NewCodebaseHasher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hasher.CalculateHash()
		if err != nil {
			b.Fatalf("Failed to calculate codebase hash: %v", err)
		}
	}
}

// BenchmarkFB2ToMarkdown benchmarks the FB2 to Markdown conversion
func BenchmarkFB2ToMarkdown(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "fb2_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fb2Path := filepath.Join(tempDir, "benchmark.fb2")
	if err := createTestFB2File(fb2Path); err != nil {
		b.Fatalf("Failed to create test FB2 file: %v", err)
	}

	testLogger := logger.NewLogger(logger.Config{
		Level:  "error", // Minimal logging for benchmark
		Format: "text",
	})

	converter := fb2.NewMarkdownConverter(testLogger)
	mdPath := filepath.Join(tempDir, "benchmark.md")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := converter.ConvertToMarkdown(fb2Path, mdPath); err != nil {
			b.Fatalf("FB2 to Markdown conversion failed: %v", err)
		}
		os.Remove(mdPath) // Clean up for next iteration
	}
}

// BenchmarkMarkdownToEPUB benchmarks the Markdown to EPUB conversion
func BenchmarkMarkdownToEPUB(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "epub_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mdPath := filepath.Join(tempDir, "benchmark.md")
	simpleMD := `---
title: Benchmark Test Book
author: Test Author
language: sr
---

# Benchmark Test Book

This is a test chapter for benchmarking EPUB generation.

## Chapter 1

This is the first paragraph of the test book.

This is the second paragraph.

### Subtitle

This paragraph follows a subtitle.

> This is a test epigraph that should be formatted as a blockquote.
> 
> — Test Quote Author

Regular text continues after the epigraph.

## Chapter 2

This is the second chapter with more content.

    This is the first verse of a test poem.
    This is the second verse of the test poem.

Text continues after the poem.

> This is a test citation that should be formatted as a blockquote.
> 
> — Test Citation Author

Final paragraph of the test book.
`

	if err := os.WriteFile(mdPath, []byte(simpleMD), 0644); err != nil {
		b.Fatalf("Failed to create test markdown file: %v", err)
	}

	converter := markdown.NewMarkdownToEPUBConverter()
	epubPath := filepath.Join(tempDir, "benchmark.epub")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := converter.ConvertMarkdownToEPUB(mdPath, epubPath); err != nil {
			b.Fatalf("Markdown to EPUB conversion failed: %v", err)
		}
		os.Remove(epubPath) // Clean up for next iteration
	}
}