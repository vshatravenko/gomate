package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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

	var state string
	err = db.Get("state", &state)
	if err != nil {
		return err
	}

	if state == "started" {
		fmt.Println("The timer is already started!")
		return nil
	}

	startTime := time.Now()

	err = db.Put("state", "started")
	if err != nil {
		return err
	}

	err = db.Put("startTime", startTime)
	if err != nil {
		return err
	}

	fmt.Println("Started the timer!")
	return nil
}

func handleStop() error {
	db, err := storage.Open(storageFile)
	if err != nil {
		return err
	}
	defer db.Close()

	var state string
	err = db.Get("state", &state)
	if err != nil {
		return err
	}

	if state == "stopped" {
		fmt.Println("The timer is already stopped!")
		return nil
	}

	defer fmt.Println("Stopped the timer!")
	return db.Put("state", "stopped")
}

func handleStatus() error {
	db, err := storage.Open(storageFile)
	if err != nil {
		return err
	}
	defer db.Close()

	var state string
	err = db.Get("state", &state)
	if err != nil {
		return err
	}

	var startTime time.Time
	err = db.Get("startTime", &startTime)
	if err != nil {
		return err
	}

	fmt.Printf("The timer is currently %s!\n", state)
	if state == "started" {
		fmt.Printf("Started at %s\n", startTime.Format("3:04PM"))
	}
	return nil
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
		return handleStatus()
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
