package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func LoadTitle() string {
	var builder strings.Builder
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	fullPath := filepath.Join(homeDir, ".passport/title.txt")

	file, err := os.Open(fullPath)
	if err != nil {
		return ""
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return ""
	}

	builder.WriteString(string(content))

	return builder.String()
}
