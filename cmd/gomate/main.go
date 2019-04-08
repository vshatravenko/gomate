package main

import (
	"fmt"
	"os"
)

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
