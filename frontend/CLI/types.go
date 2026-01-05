package main

import "github.com/fatih/color"

type Command interface {
	Execute() error
}

type StatusCommand struct{}

type LoginCommand struct{}

type SignOutCommand struct{}

type SignUpCommand struct{}

type AddCommand struct{}

type ListCommand struct{}

type HelpCommand struct{}

type GetCommand struct {
	FlagExists bool
	FlagValue  string
}

type RemoveCommand struct {
	FlagExists bool
	FlagValue  string
} 

var (
	red = color.New(color.FgRed)
	green = color.New(color.FgGreen)
)
