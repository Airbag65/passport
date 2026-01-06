package main

import (
	"fmt"
	"os"
	"passport-cli/net"
	"syscall"
	"text/tabwriter"

	"github.com/Airbag65/argparse"
	"golang.org/x/crypto/ssh/terminal"
)

func CreateCommand(pc *argparse.ParsedCommand) Command {
	switch pc.Command {
	case "status":
		return &StatusCommand{}
	case "login":
		return &LoginCommand{}
	case "signout":
		return &SignOutCommand{}
	case "signup":
		return &SignUpCommand{}
	case "add":
		return &AddCommand{}
	case "help":
		return &HelpCommand{}
	case "get":
		return &GetCommand{
			FlagExists: pc.Option != "",
			FlagValue: pc.Parameter,
		}
	case "list", "ls":
		return &ListCommand{}
	case "remove", "rm":
		return &RemoveCommand{
			FlagExists: pc.Option != "",
			FlagValue: pc.Parameter,
		}
	}
	return nil
}

func (c *StatusCommand) Execute() error {
	if net.ValidTokenExists() {
		green.Println("You are signed in to PASSPORT\nPASSPORT is ready to use")
		return nil
	}
	red.Println("You are not signed in to PASSPORT\nRun 'passport login' to sign in")
	return nil
}

func (c *LoginCommand) Execute() error {
	var email string
	fmt.Print("Email: ")
	fmt.Scan(&email)
	fmt.Print("Password: ")
	passBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return err
	}
	res, err := net.Login(email, string(passBytes))
	if err != nil {
		red.Printf("Something went wrong! err: %v\n", err)
		return err
	}
	switch res.ResponseCode {
	case 200:
		green.Printf("You are now logged in as '%s %s'\n", res.Name, res.Surname)
	case 404:
		yellow.Printf("Account with email '%s' does not exist\n", email)
	case 418:
		yellow.Printf("Already logged in with email '%s'\n", email)
	case 401:
		red.Println("Incorrect password")
	}
	return nil
}

func (c *SignOutCommand) Execute() error {
	status, err := net.SignOut()
	if err != nil {
		red.Printf("Something went wrong! err: %v\n", err)
		return err
	}

	switch status {
	case 200:
		green.Println("You are now signed out")
	case 304:
		yellow.Println("You were already signed out")
	}
	return nil
}

func (c *SignUpCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *AddCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *ListCommand) Execute() error {
	EnsureLoggedIn()
	for _, host := range net.GetHostNames() {
		fmt.Printf("%s\n", host)
	}
	return nil
}

func (c *GetCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *RemoveCommand) Execute() error {
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *HelpCommand) Execute() error {
	blue.Println(LoadTitle())	
	fmt.Println("Usage: passport <command> [flag] [<value>]")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Println("COMMANDS:")
	fmt.Fprintln(w, "\tstatus\tCheck login status")
	fmt.Fprintln(w, "\tlogin\tLogin to passport")
	fmt.Fprintln(w, "\tsignout\tSign out from passport")
	fmt.Fprintln(w, "\tsignup\tRegister a new passport account")
	fmt.Fprintln(w, "\tadd\tAdd a new password to your passport account")
	fmt.Fprintln(w, "\tget [-h --host] <hostname>\tRetrieve the password of the specified hostname")
	fmt.Fprintln(w, "\tlist\tList all the hosts you have registered passwords for")
	fmt.Fprintln(w, "\tls\tList all the hosts you have registered passwords for")
	fmt.Fprintln(w, "\tremove [-h --host] <hostname>\tRemove the password of the specified hostname. Also removes the host from passport")
	fmt.Fprintln(w, "\trm [-h --host] <hostname>\tRemove the password of the specified hostname. Also removes the host from passport")
	fmt.Fprintln(w, "\thelp\tLists all possible commands and their usage")
	w.Flush()
	return nil
}
