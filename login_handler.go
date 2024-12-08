package main

import (
	"errors"
	"fmt"
)

// handlerLogin handles the "login command".
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required")
	}
	username := cmd.args[0]
	err := s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User set to '%s'\n", username)
	return nil
}
