package main

import (
	"context"
	"errors"
	"fmt"
)

// handlerLogin handles the "login command".
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required")
	}
	username := cmd.args[0]

	// Check if the user exists in the database
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		if err.Error() == "sql: now rows in result set" {
			return fmt.Errorf("user '%s' does not exist\n", username)
		}
		return fmt.Errorf("failed to fetch user: %w\n", err)
	}

	// Set the current user in the config
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("failed to set user: %w\n", err)
	}

	fmt.Printf("Logged in as '%s'\n", username)
	return nil
}
