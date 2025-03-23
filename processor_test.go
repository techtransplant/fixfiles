package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestProcessFile tests the file processing functionality
func TestProcessFile(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "process-file-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create directory structure
	componentsDir := filepath.Join(tempDir, "src", "components")
	contextDir := filepath.Join(tempDir, "src", "context")
	apiDir := filepath.Join(tempDir, "src", "lib", "api")

	dirs := []string{componentsDir, contextDir, apiDir}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files with imports
	files := map[string]string{
		filepath.Join(componentsDir, "Component.tsx"): `
import React from 'react';
import { useAuth } from '../context/Auth';
import { getData } from '../lib/api/data';
`,
		filepath.Join(contextDir, "Auth.tsx"): `
import React, { createContext } from 'react';
export const useAuth = () => {};
`,
		filepath.Join(apiDir, "data.ts"): `
export const getData = () => {};
`,
	}

	for path, content := range files {
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Reset processed files for the test
	processedFiles = make(map[string]struct{})

	// Process the main component file
	results := make(map[string]string)
	mainFilePath := filepath.Join(componentsDir, "Component.tsx")
	err = ProcessFile(mainFilePath, tempDir, results)
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	// Check if all three files were processed
	expectedFiles := []string{
		mainFilePath,
		filepath.Join(contextDir, "Auth.tsx"),
		filepath.Join(apiDir, "data.ts"),
	}

	if len(results) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(results))
	}

	for _, expectedFile := range expectedFiles {
		if _, exists := results[expectedFile]; !exists {
			t.Errorf("Expected file %s not found in results", expectedFile)
		}
	}
}

// TestExtractImports tests the import extraction functionality
func TestExtractImports(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "extract-imports-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file with different import styles
	jsFilePath := filepath.Join(tempDir, "test.js")
	jsContent := `
import React from 'react';
import { useState } from 'react';
import './styles.css';
import AuthContext from '../context/AuthContext';
import { getData } from '@/lib/api/data';
const utils = require('./utils');
`
	err = os.WriteFile(jsFilePath, []byte(jsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Extract imports
	imports, err := ExtractImports(jsFilePath, tempDir)
	if err != nil {
		t.Fatalf("ExtractImports failed: %v", err)
	}

	// We expect only local imports to be included (not 'react')
	expectedImportCount := 4 // './styles.css', '../context/AuthContext', '@/lib/api/data', './utils'
	if len(imports) != expectedImportCount {
		t.Errorf("Expected %d imports, got %d", expectedImportCount, len(imports))
	}

	// Check for specific import paths
	containsPath := func(paths []string, substr string) bool {
		for _, path := range paths {
			if filepath.Base(path) == substr || substr == filepath.Base(filepath.Dir(path)) {
				return true
			}
		}
		return false
	}

	// Check for expected import types
	if !containsPath(imports, "styles.css") {
		t.Errorf("Expected to find styles.css in imports")
	}

	if !containsPath(imports, "AuthContext") {
		t.Errorf("Expected to find AuthContext in imports")
	}

	if !containsPath(imports, "data") {
		t.Errorf("Expected to find data.ts in imports")
	}

	if !containsPath(imports, "utils") {
		t.Errorf("Expected to find utils in imports")
	}
}
