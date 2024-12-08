package main

import (
	"fmt"
	"os"

	config "github.com/Fepozopo/gator/internal/config"
)

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize application state
	appState := &state{config: &cfg}

	// Initialize commands and register handlers
	cmds := &commands{}
	cmds.register("login", handlerLogin)

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Error: not enough arguments provided\n")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	cmd := command{name: cmdName, args: cmdArgs}

	// Run the command
	if err := cmds.run(appState, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
