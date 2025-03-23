package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFormatResults tests the results formatting functionality
func TestFormatResults(t *testing.T) {
	// Create test data
	results := map[string]string{
		"/path/to/file1.js": "console.log('Hello');\n",
		"/path/to/file2.js": "function test() { return true; }",
	}

	// Format the results
	output := FormatResults(results)

	// Check that each file is included with the correct header and footer
	for filePath, content := range results {
		expectedHeader := "{{ BEGIN CONTENTS OF " + filePath + " }}"
		expectedFooter := "{{ END CONTENTS OF " + filePath + " }}"

		if !strings.Contains(output, expectedHeader) {
			t.Errorf("Expected output to contain header: %s", expectedHeader)
		}

		if !strings.Contains(output, expectedFooter) {
			t.Errorf("Expected output to contain footer: %s", expectedFooter)
		}

		if !strings.Contains(output, content) {
			t.Errorf("Expected output to contain file content: %s", content)
		}
	}

	// Check that the summary line includes the right information
	if !strings.Contains(output, "Generated") || !strings.Contains(output, "lines of code from") || !strings.Contains(output, "files") {
		t.Errorf("Expected output to contain summary with line count and file count")
	}
}

// TestWriteResultsToFile tests the file writing functionality
func TestWriteResultsToFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "write-results-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory for the test
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalWd)

	os.Chdir(tempDir)

	// Test content
	testContent := "Test formatted content\nWith multiple lines\n"

	// Write to file
	outputFile, err := WriteResultsToFile(testContent)
	if err != nil {
		t.Fatalf("WriteResultsToFile failed: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s does not exist", outputFile)
	}

	// Read the file content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the content was written correctly
	if !strings.Contains(string(content), testContent) {
		t.Errorf("Expected output file to contain the test content")
	}

	// Check that the file name follows the expected pattern
	if !strings.HasPrefix(filepath.Base(outputFile), "error-context-") {
		t.Errorf("Expected output file name to start with 'error-context-'")
	}

	if !strings.HasSuffix(filepath.Base(outputFile), ".txt") {
		t.Errorf("Expected output file to have .txt extension")
	}
}
