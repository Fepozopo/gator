package main

import (
	"context"
	"fmt"
)

// handlerReset handles the "reset" command, which deletes all users in the database.
// It executes a query to remove all entries from the users table. If the operation
// is successful, a confirmation message is printed. If the operation fails, an error
// message is returned.
func handlerReset(s *state, cmd command) error {
	// Call the query to delete all users
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset the users table: %w", err)
	}

	fmt.Print("All users have been successfully deleted.\n")
	return nil
}
