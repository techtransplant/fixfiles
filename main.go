package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// formatPath ensures consistent path formatting
func formatPath(path string) string {
	return filepath.Clean(path)
}

func main() {
	// Parse command line arguments
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: fixfiles PATH")
		fmt.Println("  PATH: Path to the file with the error")
		os.Exit(1)
	}

	filePath := args[0]

	// Get the absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf("Error: could not get absolute path: %v\n", err)
		os.Exit(1)
	}

	// Find project root
	projectRoot, err := FindProjectRoot(absPath)
	if err != nil {
		fmt.Printf("Error: could not find project root: %v\n", err)
		os.Exit(1)
	}

	// Process the file and its dependencies
	results := make(map[string]string)
	err = ProcessFile(absPath, projectRoot, results)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	// Format and write results to file
	PrintResults(results)
}
