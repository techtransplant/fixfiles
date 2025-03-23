package main

import (
	"regexp"
)

// Store files that have already been processed to avoid duplicates
var processedFiles = make(map[string]struct{})

// File extensions to consider for import analysis
var supportedExtensions = map[string]struct{}{
	".go":    {},
	".js":    {},
	".jsx":   {},
	".ts":    {},
	".tsx":   {},
	".py":    {},
	".html":  {},
	".css":   {},
	".json":  {},
	".vue":   {},
	".scss":  {},
	".sass":  {},
	".less":  {},
	".mjs":   {},
	".cjs":   {},
	".rs":    {},
	".rb":    {},
	".php":   {},
	".java":  {},
	".swift": {},
	".kt":    {},
}

// ImportPatterns maps file extensions to regular expressions that match import statements
var importPatterns = map[string][]*regexp.Regexp{
	".js": {
		regexp.MustCompile(`import\s+.*?\s+from\s+['"](.+?)['"]`),
		regexp.MustCompile(`import\s+['"](.+?)['"]`),
		regexp.MustCompile(`require\(['"](.+?)['"]\)`),
		regexp.MustCompile(`import\(.*?['"](.+?)['"].*?\)`),
	},
	".jsx": {
		regexp.MustCompile(`import\s+.*?\s+from\s+['"](.+?)['"]`),
		regexp.MustCompile(`import\s+['"](.+?)['"]`),
		regexp.MustCompile(`require\(['"](.+?)['"]\)`),
	},
	".ts": {
		regexp.MustCompile(`import\s+.*?\s+from\s+['"](.+?)['"]`),
		regexp.MustCompile(`import\s+['"](.+?)['"]`),
		regexp.MustCompile(`import\s+type\s+.*?\s+from\s+['"](.+?)['"]`),
		regexp.MustCompile(`require\(['"](.+?)['"]\)`),
	},
	".tsx": {
		regexp.MustCompile(`import\s+.*?\s+from\s+['"](.+?)['"]`),
		regexp.MustCompile(`import\s+['"](.+?)['"]`),
		regexp.MustCompile(`import\s+type\s+.*?\s+from\s+['"](.+?)['"]`),
		regexp.MustCompile(`require\(['"](.+?)['"]\)`),
	},
	".py": {
		regexp.MustCompile(`from\s+(\S+)\s+import\s+`),
		regexp.MustCompile(`import\s+(\S+)`),
		regexp.MustCompile(`import\s+(\S+)\s+as\s+`),
	},
	".go": {
		regexp.MustCompile(`import\s+[(\s]+"(.+?)"`),
		regexp.MustCompile(`import\s+(\S+)\s+".+?"`),
	},
	".html": {
		regexp.MustCompile(`<script\s+src=['"](.+?)['"]\s*>`),
		regexp.MustCompile(`<link\s+.*?href=['"](.+?)['"]\s*>`),
		regexp.MustCompile(`<img\s+.*?src=['"](.+?)['"]\s*>`),
	},
	".css": {
		regexp.MustCompile(`@import\s+['"](.+?)['"]`),
		regexp.MustCompile(`@import\s+url\(['"](.+?)['"]\)`),
	},
	".scss": {
		regexp.MustCompile(`@import\s+['"](.+?)['"]`),
		regexp.MustCompile(`@import\s+url\(['"](.+?)['"]\)`),
		regexp.MustCompile(`@use\s+['"](.+?)['"]`),
	},
	".vue": {
		regexp.MustCompile(`import\s+.*?\s+from\s+['"](.+?)['"]`),
		regexp.MustCompile(`import\s+['"](.+?)['"]`),
		regexp.MustCompile(`require\(['"](.+?)['"]\)`),
		regexp.MustCompile(`<script\s+src=['"](.+?)['"]\s*>`),
	},
	".php": {
		regexp.MustCompile(`require[_once]*\s*\(['"](.+?)['"]\)`),
		regexp.MustCompile(`include[_once]*\s*\(['"](.+?)['"]\)`),
		regexp.MustCompile(`use\s+([^;]+)`),
	},
	".rb": {
		regexp.MustCompile(`require\s+['"](.+?)['"]`),
		regexp.MustCompile(`require_relative\s+['"](.+?)['"]`),
		regexp.MustCompile(`load\s+['"](.+?)['"]`),
	},
	".java": {
		regexp.MustCompile(`import\s+([^;]+)`),
	},
	".kt": {
		regexp.MustCompile(`import\s+([^;]+)`),
	},
}

// Initialize patterns for other file types that use the same patterns as JS
func init() {
	// JavaScript-like imports
	jsLike := []string{".mjs", ".cjs"}
	for _, ext := range jsLike {
		importPatterns[ext] = importPatterns[".js"]
	}

	// CSS-like imports
	cssLike := []string{".less", ".sass"}
	for _, ext := range cssLike {
		importPatterns[ext] = importPatterns[".css"]
	}

	// File types that typically don't have imports but may be referenced
	noImports := []string{".json", ".svg", ".png", ".jpg", ".jpeg", ".gif"}
	for _, ext := range noImports {
		importPatterns[ext] = []*regexp.Regexp{}
	}
}
