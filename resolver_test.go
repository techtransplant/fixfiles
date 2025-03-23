package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFindProjectRoot tests the project root finding functionality
func TestFindProjectRoot(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "project-root-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a nested directory structure
	subDir := filepath.Join(tempDir, "src", "components")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	// Create a package.json file to indicate project root
	rootIndicator := filepath.Join(tempDir, "package.json")
	err = os.WriteFile(rootIndicator, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create root indicator file: %v", err)
	}

	// Test finding project root from the subdirectory
	foundRoot, err := FindProjectRoot(subDir)
	if err != nil {
		t.Fatalf("FindProjectRoot failed: %v", err)
	}

	// Normalize paths for comparison
	expectedRoot := filepath.Clean(tempDir)
	foundRoot = filepath.Clean(foundRoot)

	if foundRoot != expectedRoot {
		t.Errorf("FindProjectRoot returned wrong path. Got: %s, Want: %s", foundRoot, expectedRoot)
	}
}

// TestResolveImportPath tests the import path resolution
func TestResolveImportPath(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "import-path-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create directory structure
	srcDir := filepath.Join(tempDir, "src")
	componentsDir := filepath.Join(srcDir, "components")
	contextDir := filepath.Join(srcDir, "context")
	libDir := filepath.Join(srcDir, "lib", "api")

	dirs := []string{componentsDir, contextDir, libDir}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files
	files := map[string]string{
		filepath.Join(componentsDir, "TestComponent.tsx"): "// Test component",
		filepath.Join(contextDir, "AuthContext.tsx"):      "// Auth context",
		filepath.Join(libDir, "climate.ts"):               "// Climate API",
	}

	for path, content := range files {
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	// Test cases
	testCases := []struct {
		name         string
		importPath   string
		currentDir   string
		projectRoot  string
		expectedPath string
	}{
		{
			name:         "Relative import from component",
			importPath:   "../context/AuthContext",
			currentDir:   componentsDir,
			projectRoot:  tempDir,
			expectedPath: filepath.Join(contextDir, "AuthContext.tsx"),
		},
		{
			name:         "Import with @ prefix",
			importPath:   "@/lib/api/climate",
			currentDir:   componentsDir,
			projectRoot:  tempDir,
			expectedPath: filepath.Join(libDir, "climate.ts"),
		},
		{
			name:         "Already has extension",
			importPath:   "../context/AuthContext.tsx",
			currentDir:   componentsDir,
			projectRoot:  tempDir,
			expectedPath: filepath.Join(contextDir, "AuthContext.tsx"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resolved, err := ResolveImportPath(tc.importPath, tc.currentDir, tc.projectRoot)
			if err != nil {
				t.Fatalf("ResolveImportPath failed: %v", err)
			}

			// Normalize paths for comparison
			expected := filepath.Clean(tc.expectedPath)
			resolved = filepath.Clean(resolved)

			if resolved != expected {
				t.Errorf("ResolveImportPath returned wrong path. Got: %s, Want: %s", resolved, expected)
			}
		})
	}
}
