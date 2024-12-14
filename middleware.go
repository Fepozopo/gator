package main

import (
	"context"
	"fmt"

	"github.com/Fepozopo/gator/internal/database"
	"github.com/google/uuid"
)

// middlewareLoggedIn wraps a command handler with middleware that checks if the user is logged in before executing the handler.
// If the user is not logged in, the middleware returns an error.
// If the user is logged in, the middleware passes the state, command, and current user to the wrapped handler.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Fetch the current user's name from the config
		if s.cfg.CurrentUserName == "" {
			return fmt.Errorf("no user is currently logged in")
		}

		// Retrieve the current user from the database
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to fetch current user: %w", err)
		}

		// Ensure the user exists
		if currentUser.ID == uuid.Nil {
			return fmt.Errorf("no logged-in user found")
		}

		// Pass the state, command, and current user to the wrapped handler
		return handler(s, cmd, *&currentUser)
	}
}
