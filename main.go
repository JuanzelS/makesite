package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

type Page struct {
	Content string
}

func main() {
	// Flags
	fileFlag := flag.String("file", "", "The .txt file to convert into HTML")
	dirFlag := flag.String("dir", "", "Directory containing .txt files to convert into HTML")
	flag.Parse()

	// Handle --file
	if *fileFlag != "" {
		generateHTML(*fileFlag)
		return
	}

	// Handle --dir
	if *dirFlag != "" {
		err := filepath.Walk(*dirFlag, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
				fmt.Println("Found:", info.Name())
				generateHTML(path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	fmt.Println("Please provide either --file=<filename> or --dir=<directory>")
}

func generateHTML(inputFile string) {
	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", inputFile, err)
	}

	// Default to raw text
	htmlContent := string(content)

	// If Markdown, convert with goldmark
	if strings.HasSuffix(inputFile, ".md") {
		var buf bytes.Buffer
		if err := goldmark.Convert(content, &buf); err != nil {
			log.Fatalf("Error converting Markdown: %v", err)
		}
		htmlContent = buf.String()
	}

	// Parse template
	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	page := Page{Content: htmlContent}

	// Determine output file name
	ext := filepath.Ext(inputFile) // handles both .txt and .md
	outputFile := strings.TrimSuffix(inputFile, ext) + ".html"

	// Write to output file
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", outputFile, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, page)
	if err != nil {
		log.Fatalf("Error writing HTML to file: %v", err)
	}

	fmt.Printf("Generated: %s\n", outputFile)
}