package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveImportPath attempts to resolve an import path to a real file path
func ResolveImportPath(importPath string, currentDir, projectRoot string) (string, error) {
	// First, try standard resolution based on import type
	resolvedPath, err := resolveImportByType(importPath, currentDir, projectRoot)
	if err != nil {
		return "", err
	}

	// If the resolved path already has an extension, check if it exists directly
	if filepath.Ext(resolvedPath) != "" {
		if _, err := os.Stat(resolvedPath); err == nil {
			return resolvedPath, nil
		}
	} else {
		// Try with each supported extension
		for ext := range supportedExtensions {
			testPath := resolvedPath + ext
			if _, err := os.Stat(testPath); err == nil {
				return testPath, nil
			}
		}

		// For JS/TS, try index files
		jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
		for _, ext := range jsExtensions {
			if _, ok := supportedExtensions[ext]; ok {
				indexFile := filepath.Join(resolvedPath, "index"+ext)
				if _, err := os.Stat(indexFile); err == nil {
					return indexFile, nil
				}
			}
		}
	}

	// If we still can't find the file, return the best guess
	return resolvedPath, nil
}

// resolveImportByType handles different import styles
func resolveImportByType(importPath string, currentDir, projectRoot string) (string, error) {
	// Handle different import styles
	if strings.HasPrefix(importPath, ".") {
		// Relative import (e.g., "./utils" or "../components")
		absPath, err := filepath.Abs(filepath.Join(currentDir, importPath))
		if err != nil {
			return "", err
		}
		return absPath, nil
	} else if strings.HasPrefix(importPath, "/") {
		// Absolute import within the project
		return filepath.Join(projectRoot, importPath), nil
	} else if strings.HasPrefix(importPath, "~") {
		// Some projects use ~ to refer to the project root
		return filepath.Join(projectRoot, strings.TrimPrefix(importPath, "~")), nil
	} else if strings.HasPrefix(importPath, "@") {
		// Handle aliased imports like @/components
		// For Next.js, @ often refers to src directory
		srcPath := filepath.Join(projectRoot, "src", strings.TrimPrefix(importPath, "@/"))

		// Try with src directory first
		if _, err := os.Stat(srcPath); err == nil {
			return srcPath, nil
		}

		// If file not found in src, try adding extension
		for ext := range supportedExtensions {
			testPath := srcPath + ext
			if _, err := os.Stat(testPath); err == nil {
				return testPath, nil
			}
		}

		// If still not found, fall back to project root
		return filepath.Join(projectRoot, strings.TrimPrefix(importPath, "@/")), nil
	}

	// For non-relative imports, try to find them in node_modules or similar
	possiblePaths := []string{
		filepath.Join(projectRoot, "node_modules", importPath),
		filepath.Join(projectRoot, "vendor", importPath),
		filepath.Join(projectRoot, "lib", importPath),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// If we can't resolve it, just return as is - we'll try to resolve extensions later
	return importPath, nil
}

// FindProjectRoot attempts to find the root directory of the project
func FindProjectRoot(startPath string) (string, error) {
	// Common project root indicators
	indicators := []string{
		"go.mod",
		"package.json",
		".git",
		"requirements.txt",
		"setup.py",
		"docker-compose.yml",
		"Makefile",
	}

	dir := startPath
	for {
		for _, indicator := range indicators {
			if _, err := os.Stat(filepath.Join(dir, indicator)); err == nil {
				return dir, nil
			}
		}

		// Move up one directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// We've reached the filesystem root
			break
		}
		dir = parentDir
	}

	// If we couldn't find a project root, use the directory of the initial file
	return filepath.Dir(startPath), nil
}
