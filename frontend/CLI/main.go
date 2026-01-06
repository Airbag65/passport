package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	p, err := InitParser()
	if err != nil {
		log.Fatal(err)
	}
	command, err := p.Parse(os.Args)
	if err != nil {
		fmt.Println("Invalid usage.\nUsage: passport <command>\nRun 'passport help' for a list of commands and their usage")
		panic("error")
	}
	c := CreateCommand(command)
	c.Execute()
}

