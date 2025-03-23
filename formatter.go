package main

import (
	"fmt"
	"strings"
)

// FormatResults formats the collected file contents for output
func FormatResults(results map[string]string) string {
	var builder strings.Builder

	for filePath, content := range results {
		fmt.Fprintf(&builder, "{{ BEGIN CONTENTS OF %s }}\n", filePath)
		fmt.Fprint(&builder, content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Fprint(&builder, "\n")
		}
		fmt.Fprintf(&builder, "{{ END CONTENTS OF %s }}\n\n", filePath)
	}

	totalLines := 0
	for _, content := range results {
		totalLines += strings.Count(content, "\n") + 1
	}

	fmt.Fprintf(&builder, "------------------------------\n")
	fmt.Fprintf(&builder, "Generated %d lines of code from %d files\n", totalLines, len(results))

	return builder.String()
}

// PrintResults outputs the formatted results to stdout and file
func PrintResults(results map[string]string) {
	formattedContent := FormatResults(results)
	fmt.Print(formattedContent)
}
