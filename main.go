package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

// Page holds the text content to pass into the template
type Page struct {
	Content string
}

func main() {
	// CLI flag for choosing which .txt file to use
	fileFlag := flag.String("file", "first-post.txt", "The .txt file to convert into HTML")
	flag.Parse()

	// Read the input file
	content, err := os.ReadFile(*fileFlag)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", *fileFlag, err)
	}

	// Parse the template
	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Prepare the data
	page := Page{Content: string(content)}

	// Render to stdout
	fmt.Println("===== Rendering HTML to stdout =====")
	err = tmpl.Execute(os.Stdout, page)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	// Determine output file name (.txt → .html)
	outputFile := strings.TrimSuffix(*fileFlag, ".txt") + ".html"

	// Write to output file
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating output file %s: %v", outputFile, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, page)
	if err != nil {
		log.Fatalf("Error writing HTML to file: %v", err)
	}

	fmt.Printf("\nHTML file successfully generated: %s\n", outputFile)
}
