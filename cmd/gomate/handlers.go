package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vshatravenko/gomate/pkg/storage"
)

var (
	commands    = []string{"start", "stop", "status", "daemon"}
	storageDir  = os.Getenv("HOME") + "/.gomate"
	storageFile = storageDir + "/main.db"
)

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

	err = db.Put("startTime", startTime)
	if err != nil {
		return err
	}

	err = db.Put("remainingTime", 25*time.Minute)
	if err != nil {
		return err
	}

	err = db.Put("state", "started")
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

	err = db.Put("remainingTime", time.Duration(0))
	if err != nil {
		return err
	}

	defer fmt.Println("Stopped the timer!")
	return db.Put("state", "stopped")
}

// handleDaemon starts a process which would check the state of GoMate and act accordingly each second
func handleDaemon() error {
	sigs := make(chan os.Signal, 1)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Register sigs to receive shutdown signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("GoMate daemon started successfully!\nPress Ctrl+C to stop")
	for {
		select {
		case <-ticker.C:
			db, err := storage.Open(storageFile)

			if err != nil {
				return err
			}

			var state string
			err = db.Get("state", &state)
			if err != nil {
				return err
			}

			if state == "started" {
				var remainingTime time.Duration
				err = db.Get("remainingTime", &remainingTime)
				if err != nil {
					return err
				}

				remainingTime -= 1 * time.Second
				err = db.Put("remainingTime", remainingTime)
				if err != nil {
					return err
				}

				if remainingTime == 0*time.Second {
					err = db.Put("state", "stopped")
					if err != nil {
						return err
					}
				}

				if err = db.Close(); err != nil {
					return err
				}
			}

		case <-sigs:
			fmt.Println("\nReceived a shutdown signal, exiting")
			return nil
		}
	}
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

	var remainingTime time.Duration
	err = db.Get("remainingTime", &remainingTime)
	if err != nil {
		return err
	}

	fmt.Printf("The timer is currently %s!\n", state)
	if state == "started" {
		fmt.Printf("Started at %s\n", startTime.Format("3:04PM"))
		fmt.Printf("Remaining time: %s\n", remainingTime)
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
	case "daemon":
		return handleDaemon()
	default:
		fmt.Println("Available commands:", strings.Join(commands, " "))
		return errors.New("unrecognized command")
	}
}
