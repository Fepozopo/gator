package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	// Call the query to delete all users
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset the users table: %w\n", err)
	}

	fmt.Print("All users have been successfully deleted.\n")
	return nil
}
