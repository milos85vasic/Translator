package main

import (
	"fmt"
	"digital.vasic.translator/pkg/markdown"
	"digital.vasic.translator/pkg/ebook"
)

func main() {
	c := markdown.NewMarkdownToEPUBConverter()
	c.ConvertMarkdownToEPUB("/Users/milosvasic/Projects/Translate/test_sample.md", "/tmp/test_sample.epub")
	p := ebook.NewUniversalParser()
	b, _ := p.Parse("/tmp/test_sample.epub")
	fmt.Printf("Chapters: %d\n", len(b.Chapters))
}