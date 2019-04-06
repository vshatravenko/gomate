package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var commands = []string{"start", "stop", "status", "daemon"}

func handleCmd() error {
	if len(os.Args[1:]) == 0 {
		return errors.New("no commands passed")
	}

	switch os.Args[1] {
	case "start":
		fmt.Println("start")
		return nil
	case "stop":
		fmt.Println("start")
		return nil
	case "status":
		fmt.Println("start")
		return nil
	default:
		fmt.Println("Available commands:", strings.Join(commands, " "))
		return errors.New("unrecognized command")
	}
}

func main() {
	err := handleCmd()

	if err != nil {
		fmt.Printf("Error while handling a command: %s\n", err)
		os.Exit(-1)
	}
}
