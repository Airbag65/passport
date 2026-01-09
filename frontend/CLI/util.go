package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"passport-cli/net"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Airbag65/argparse"
	"golang.org/x/term"
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

func GetPassword(prompt string) string {
	fmt.Print(prompt)
	passBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	fmt.Println("")
	return string(passBytes)
}

func InitParser() (*argparse.Parser, error) {
	commands := []string{"status", "login", "signout", "signup", "add", "get", "list", "ls", "remove", "rm", "help"}

	p := argparse.New()
	hostDesc := "Specify which host to direct the command at"
	hostFlag := argparse.NewFlag("--host", hostDesc, true)
	hFlag := argparse.NewFlag("-h", hostDesc, true)
	for _, comm := range commands {
		switch comm {
		case "get", "remove", "rm":
			err := p.AddCommand(comm, argparse.AddFlag(hostFlag), argparse.AddFlag(hFlag))
			if err != nil {
				return nil, err
			}
		default:
			err := p.AddCommand(comm)
			if err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}

func YesNoConfirmation(prompt string, defaultYes bool) bool {
	fmt.Print(prompt)
	if defaultYes {
		fmt.Print(" [Y/n] ")
	} else {
		fmt.Print(" [y/N] ")
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Text() == "Y" || scanner.Text() == "y" {
		return true
	} else if scanner.Text() == "" && defaultYes {
		return true
	}
	return false
}

func PrintFramedWord(word string) {
	fmt.Printf("+%s+\n|", strings.Repeat("-", len(word)))
	fmt.Print(word)
	fmt.Printf("|\n+%s+\n", strings.Repeat("-", len(word)))
}
