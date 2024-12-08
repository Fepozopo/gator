package main

import (
	"fmt"
	"log"

	config "github.com/Fepozopo/gator/internal/config"
)

func main() {
	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	fmt.Printf("Initial Config: %+v\n", cfg)

	// Set the current user to "shane" and update the config file
	err = cfg.SetUser("shane")
	if err != nil {
		log.Fatalf("Failed to update user in config: %v", err)
	}
	fmt.Print("Updated current_user_name to 'shane'.\n")

	// Read the config file again and print the contents
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Failed to read config after update: %v", err)
	}
	fmt.Printf("Updated Config: %+v\n", cfg)
}
