package main

import (
	"io"
	"os"
	"passport-cli/net"
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

func EnsureLoggedIn() {
	if !net.ValidTokenExists() {
		red.Println("You are signed out and are thus unable to use PASSPORT\nRun 'passport login' to login")
		panic("Not Logged In")
	}
}
