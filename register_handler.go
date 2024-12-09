package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	database "github.com/Fepozopo/gator/internal/database"
	"github.com/google/uuid"
)

// handlerRegister handles the "register" command, which creates a new user
// with the given username. The username is taken from the command arguments.
// If the username already exists, an error is returned. Upon successful
// creation, the user is set as the current user in the configuration, and
// a success message is printed with the user details.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required")
	}
	name := cmd.args[0]

	// Create a new user
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	// Use the generated query to insert the user
	user, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_name_key\"" {
			return fmt.Errorf("a user with the name '%s' already exists\n", name)
		}
		return fmt.Errorf("failed to create user: %w\n", err)
	}

	// Set the current user in the config
	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to set current user: %w\n", err)
	}

	// Print a success message and log the user details
	fmt.Printf("User '%s' created successfully\n", user.Name)
	fmt.Printf("User Details: %+v\n", user)
	return nil
}
