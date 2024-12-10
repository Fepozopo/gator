package main

import (
	"context"
	"fmt"
)

// handlerUsers handles the "users" command, which lists all users in the database.
// The list will show the currently logged-in user with "(current)" appended to their name.
func handlerUsers(s *state, cmd command) error {
	// Fetch all users from the database
	users, err := s.db.GetUsers((context.Background()))
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	// Get the currently logged-in user from the config
	currentUser := s.cfg.CurrentUserName

	// Print the list of users
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
