package fb2

import (
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Fatal("NewParser() returned nil")
	}
}

func TestParse(t *testing.T) {
	parser := NewParser()

	// Test parsing a valid FB2 file
	tempDir := t.TempDir()
	validFB2 := `<?xml version="1.0" encoding="utf-8"?>
<FictionBook xmlns="http://www.gribuser.ru/xml/fictionbook/2.0">
  <description>
    <title-info>
      <genre>detective</genre>
      <author>
        <first-name>John</first-name>
        <last-name>Doe</last-name>
      </author>
      <book-title>Test Book</book-title>
      <lang>ru</lang>
    </title-info>
  </description>
  <body>
    <section>
      <p>Test paragraph</p>
    </section>
  </body>
</FictionBook>`

	filename := filepath.Join(tempDir, "test.fb2")
	err := os.WriteFile(filename, []byte(validFB2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fb, err := parser.Parse(filename)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	if fb == nil {
		t.Fatal("Parse() returned nil FictionBook")
	}

	if fb.GetTitle() != "Test Book" {
		t.Errorf("Expected title 'Test Book', got '%s'", fb.GetTitle())
	}

	if fb.GetLanguage() != "ru" {
		t.Errorf("Expected language 'ru', got '%s'", fb.GetLanguage())
	}
}

func TestParseReader(t *testing.T) {
	parser := NewParser()

	validFB2 := `<?xml version="1.0" encoding="utf-8"?>
<FictionBook xmlns="http://www.gribuser.ru/xml/fictionbook/2.0">
  <description>
    <title-info>
      <book-title>Test Book</book-title>
      <lang>ru</lang>
    </title-info>
  </description>
  <body>
    <section>
      <p>Test content</p>
    </section>
  </body>
</FictionBook>`

	reader := strings.NewReader(validFB2)
	fb, err := parser.ParseReader(reader)
	if err != nil {
		t.Fatalf("ParseReader() failed: %v", err)
	}

	if fb.GetTitle() != "Test Book" {
		t.Errorf("Expected title 'Test Book', got '%s'", fb.GetTitle())
	}
}

func TestParseInvalidXML(t *testing.T) {
	parser := NewParser()

	invalidXML := `<FictionBook xmlns="http://www.gribuser.ru/xml/fictionbook/2.0">
  <description>
    <title-info>
      <book-title>Unclosed tag`

	reader := strings.NewReader(invalidXML)
	_, err := parser.ParseReader(reader)
	if err == nil {
		t.Error("ParseReader() should have failed with invalid XML")
	}
}

func TestParseNonexistentFile(t *testing.T) {
	parser := NewParser()

	_, err := parser.Parse("nonexistent.fb2")
	if err == nil {
		t.Error("Parse() should have failed with nonexistent file")
	}
}

func TestWrite(t *testing.T) {
	parser := NewParser()

	fb := &FictionBook{
		Description: Description{
			TitleInfo: TitleInfo{
				BookTitle: "Test Book",
				Lang:      "ru",
				Author: []Author{{
					FirstName: "John",
					LastName:  "Doe",
				}},
			},
		},
		Body: []Body{{
			Section: []Section{{
				Paragraph: []Paragraph{{
					Text: "Test paragraph",
				}},
			}},
		}},
	}

	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "output.fb2")

	err := parser.Write(filename, fb)
	if err != nil {
		t.Fatalf("Write() failed: %v", err)
	}

	// Verify file was created and can be read back
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Write() did not create output file")
	}

	// Parse it back to verify content
	fb2, err := parser.Parse(filename)
	if err != nil {
		t.Fatalf("Failed to parse written file: %v", err)
	}

	if fb2.GetTitle() != "Test Book" {
		t.Errorf("Round-trip title mismatch: expected 'Test Book', got '%s'", fb2.GetTitle())
	}
}

func TestWriteToWriter(t *testing.T) {
	parser := NewParser()

	fb := &FictionBook{
		Description: Description{
			TitleInfo: TitleInfo{
				BookTitle: "Writer Test",
				Lang:      "sr",
			},
		},
	}

	var buf bytes.Buffer
	err := parser.WriteToWriter(&buf, fb)
	if err != nil {
		t.Fatalf("WriteToWriter() failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Writer Test") {
		t.Error("WriteToWriter() output does not contain expected title")
	}

	if !strings.Contains(output, xml.Header) {
		t.Error("WriteToWriter() output does not contain XML header")
	}
}

func TestGetLanguage(t *testing.T) {
	fb := &FictionBook{
		Description: Description{
			TitleInfo: TitleInfo{
				Lang: "ru",
			},
		},
	}

	if lang := fb.GetLanguage(); lang != "ru" {
		t.Errorf("GetLanguage() = %s, want ru", lang)
	}
}

func TestSetLanguage(t *testing.T) {
	fb := &FictionBook{}

	fb.SetLanguage("sr")
	if lang := fb.GetLanguage(); lang != "sr" {
		t.Errorf("SetLanguage() failed: got %s, want sr", lang)
	}
}

func TestGetTitle(t *testing.T) {
	fb := &FictionBook{
		Description: Description{
			TitleInfo: TitleInfo{
				BookTitle: "Test Title",
			},
		},
	}

	if title := fb.GetTitle(); title != "Test Title" {
		t.Errorf("GetTitle() = %s, want Test Title", title)
	}
}

func TestSetTitle(t *testing.T) {
	fb := &FictionBook{}

	fb.SetTitle("New Title")
	if title := fb.GetTitle(); title != "New Title" {
		t.Errorf("SetTitle() failed: got %s, want New Title", title)
	}
}

func TestComplexFB2Structure(t *testing.T) {
	parser := NewParser()

	complexFB2 := `<?xml version="1.0" encoding="utf-8"?>
<FictionBook xmlns="http://www.gribuser.ru/xml/fictionbook/2.0" xmlns:l="http://www.w3.org/1999/xlink">
  <description>
    <title-info>
      <genre>detective</genre>
      <genre>thriller</genre>
      <author>
        <first-name>Иван</first-name>
        <middle-name>Иванович</middle-name>
        <last-name>Иванов</last-name>
      </author>
      <book-title>Тестовая книга</book-title>
      <annotation>
        <p>Аннотация книги</p>
      </annotation>
      <lang>ru</lang>
      <src-lang>ru</src-lang>
    </title-info>
    <document-info>
      <author>
        <first-name>Иван</first-name>
        <last-name>Иванов</last-name>
      </author>
      <date value="2023-01-01">2023</date>
      <id>test-id</id>
      <version>1.0</version>
    </document-info>
  </description>
  <body>
    <title>
      <p>Глава 1</p>
    </title>
    <section id="section1">
      <title>
        <p>Введение</p>
      </title>
      <p>Первый абзац текста.</p>
      <p>Второй абзац с <emphasis>выделением</emphasis> и <strong>жирным</strong> текстом.</p>
      <empty-line/>
      <p>Третий абзац.</p>
    </section>
    <section>
      <title>
        <p>Поэзия</p>
      </title>
      <poem>
        <stanza>
          <v>Строка 1</v>
          <v>Строка 2</v>
        </stanza>
      </poem>
    </section>
  </body>
  <binary id="cover.jpg" content-type="image/jpeg">
    /9j/4AAQSkZJRgABAQEAYABgAAD//2Q==
  </binary>
</FictionBook>`

	reader := strings.NewReader(complexFB2)
	fb, err := parser.ParseReader(reader)
	if err != nil {
		t.Fatalf("Failed to parse complex FB2: %v", err)
	}

	// Test various fields
	if fb.GetTitle() != "Тестовая книга" {
		t.Errorf("Title mismatch: got '%s'", fb.GetTitle())
	}

	if fb.GetLanguage() != "ru" {
		t.Errorf("Language mismatch: got '%s'", fb.GetLanguage())
	}

	if len(fb.Description.TitleInfo.Author) != 1 {
		t.Error("Expected 1 author")
	}

	if len(fb.Body) != 1 {
		t.Error("Expected 1 body")
	}

	if len(fb.Body[0].Section) != 2 {
		t.Error("Expected 2 sections")
	}

	if len(fb.Binary) != 1 {
		t.Error("Expected 1 binary element")
	}
}

func TestXMLNamespaces(t *testing.T) {
	// Test that XML namespaces are handled correctly
	fb := &FictionBook{
		Description: Description{
			TitleInfo: TitleInfo{
				BookTitle: "Namespace Test",
				Lang:      "ru",
			},
		},
	}

	parser := NewParser()
	var buf bytes.Buffer
	err := parser.WriteToWriter(&buf, fb)
	if err != nil {
		t.Fatalf("WriteToWriter() failed: %v", err)
	}

	output := buf.String()

	// Check that the FictionBook element has the correct namespace
	if !strings.Contains(output, `<FictionBook xmlns="http://www.gribuser.ru/xml/fictionbook/2.0"`) {
		t.Error("Output does not contain correct FictionBook namespace")
	}
}
