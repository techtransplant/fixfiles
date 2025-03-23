package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProcessFile analyzes a file and its dependencies recursively
func ProcessFile(filePath string, projectRoot string, results map[string]string) error {
	// Normalize the path
	filePath = filepath.Clean(filePath)

	// Skip if we've already processed this file
	if _, exists := processedFiles[filePath]; exists {
		return nil
	}

	// Check if the file exists
	info, err := os.Stat(filePath)
	if err != nil {
		// Try to find the file with extensions if it doesn't have one
		foundFile := false
		if filepath.Ext(filePath) == "" {
			for ext := range supportedExtensions {
				testPath := filePath + ext
				if _, err := os.Stat(testPath); err == nil {
					filePath = testPath
					foundFile = true
					break
				}
			}
		}

		if !foundFile {
			// Try index.* files for directories
			if dirInfo, dirErr := os.Stat(filePath); dirErr == nil && dirInfo.IsDir() {
				for ext := range supportedExtensions {
					indexPath := filepath.Join(filePath, "index"+ext)
					if _, err := os.Stat(indexPath); err == nil {
						filePath = indexPath
						foundFile = true
						break
					}
				}
			}
		}

		if !foundFile {
			return fmt.Errorf("file not found: %s", filePath)
		}

		info, err = os.Stat(filePath)
		if err != nil {
			return err
		}
	}

	// Skip directories
	if info.IsDir() {
		return nil
	}

	// Skip unsupported file types
	fileExt := filepath.Ext(filePath)
	if _, supported := supportedExtensions[fileExt]; !supported {
		return nil
	}

	// Mark as processed
	processedFiles[filePath] = struct{}{}

	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Add to results
	results[filePath] = string(content)

	// Extract imports
	imports, err := ExtractImports(filePath, projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to extract imports from %s: %v\n", filePath, err)
		// Continue even if we can't extract imports
	}

	// Process each imported file recursively
	for _, importPath := range imports {
		err = ProcessFile(importPath, projectRoot, results)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not process import %s: %v\n", importPath, err)
		}
	}

	return nil
}

// ExtractImports finds all import statements in a file
func ExtractImports(filePath string, projectRoot string) ([]string, error) {
	fileExt := filepath.Ext(filePath)
	patterns, ok := importPatterns[fileExt]
	if !ok {
		return nil, fmt.Errorf("unsupported file extension: %s", fileExt)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var imports []string
	fileDir := filepath.Dir(filePath)

	for _, pattern := range patterns {
		matches := pattern.FindAllSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				importPath := string(match[1])

				// Try to resolve the import path to an actual file
				resolvedPath, err := ResolveImportPath(importPath, fileDir, projectRoot)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not resolve import path %s: %v\n", importPath, err)
					continue
				}

				// Check if this is a relative path that might be part of the project
				if strings.HasPrefix(importPath, ".") ||
					strings.HasPrefix(importPath, "/") ||
					strings.HasPrefix(importPath, "@") ||
					strings.HasPrefix(importPath, "~") ||
					// For local imports without special prefixes
					(!strings.Contains(importPath, "/") && !isBuiltinModule(importPath, fileExt)) {
					imports = append(imports, resolvedPath)
				}
				// Skip external imports (node_modules, npm packages, etc.)
			}
		}
	}

	return imports, nil
}

// isBuiltinModule checks if an import refers to a built-in module
func isBuiltinModule(importPath string, fileExt string) bool {
	// JavaScript/TypeScript built-in modules
	if fileExt == ".js" || fileExt == ".jsx" || fileExt == ".ts" || fileExt == ".tsx" {
		builtins := []string{"react", "react-dom", "next", "axios", "lodash"}
		for _, builtin := range builtins {
			if importPath == builtin || strings.HasPrefix(importPath, builtin+"/") {
				return true
			}
		}
	}

	// Python built-in modules
	if fileExt == ".py" {
		builtins := []string{"os", "sys", "re", "math", "json", "datetime", "collections"}
		for _, builtin := range builtins {
			if importPath == builtin || strings.HasPrefix(importPath, builtin+".") {
				return true
			}
		}
	}

	// Go built-in packages
	if fileExt == ".go" {
		builtins := []string{"fmt", "os", "io", "net", "http", "strings", "time"}
		for _, builtin := range builtins {
			if importPath == builtin || strings.HasPrefix(importPath, builtin+"/") {
				return true
			}
		}
	}

	return false
}

// FindImportsFromError tries to extract file paths from an error message
func FindImportsFromError(errorMsg string, projectRoot string) []string {
	var paths []string

	// Common patterns for file paths in error messages
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`in file '(.*?)'`),
		regexp.MustCompile(`in ([\w\/\.-]+\.\w+) on line \d+`),
		regexp.MustCompile(`cannot find module '(.*?)'`),
		regexp.MustCompile(`([\w\/\.-]+\.\w+):\d+:\d+`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(errorMsg, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				path := match[1]
				// Convert to absolute path if it's relative
				if !filepath.IsAbs(path) {
					path = filepath.Join(projectRoot, path)
				}
				paths = append(paths, path)
			}
		}
	}

	return paths
}
