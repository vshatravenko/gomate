package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/vshatravenko/gomate/pkg/storage"
)

var (
	commands    = []string{"start", "stop", "status", "daemon"}
	storageDir  = os.Getenv("HOME") + "/.gomate"
	storageFile = storageDir + "/main.db"
)

func handleStart() error {
	db, err := storage.Open(storageFile)
	if err != nil {
		return err
	}
	defer db.Close()

	defer fmt.Println("Started the timer!")
	return db.Put("state", "started")
}

func handleStop() error {
	db, err := storage.Open(storageFile)
	if err != nil {
		return err
	}
	defer db.Close()

	defer fmt.Println("Stopped the timer!")
	return db.Put("state", "stopped")
}

func handleCmd() error {
	if len(os.Args[1:]) == 0 {
		return errors.New("no commands passed")
	}

	switch os.Args[1] {
	case "start":
		return handleStart()
	case "stop":
		return handleStop()
	case "status":
		fmt.Println("status")
		return nil
	default:
		fmt.Println("Available commands:", strings.Join(commands, " "))
		return errors.New("unrecognized command")
	}
}

func initStorage() error {
	if _, ok := os.Stat(storageDir); os.IsNotExist(ok) {
		fmt.Printf("Initializing the storage directory at %s\n", storageDir)
		err := os.Mkdir(storageDir, 0755)

		if err != nil {
			return err
		}

		fmt.Println("Initializing the storage")
		db, err := storage.Open(storageFile)
		if err != nil {
			return err
		}
		defer db.Close()

		return db.Put("state", "stopped")
	}

	return nil
}

func main() {
	err := initStorage()
	if err != nil {
		fmt.Printf("Error while initializing the storage dir: %s\n", err)
		os.Exit(1)
	}

	err = handleCmd()
	if err != nil {
		fmt.Printf("Error while handling a command: %s\n", err)
		os.Exit(-1)
	}
}
