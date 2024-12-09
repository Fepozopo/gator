package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	config "github.com/Fepozopo/gator/internal/config"
	database "github.com/Fepozopo/gator/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		os.Exit(1)
	}

	// Connect to the database
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to the database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize the database queries
	dbQueries := database.New(db)

	// Initialize application state
	appState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Initialize commands and register handlers
	cmds := &commands{}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

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
