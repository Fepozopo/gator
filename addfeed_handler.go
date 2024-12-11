package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Fepozopo/gator/internal/database"
	"github.com/google/uuid"
)

// handlerAddFeed adds a new feed for the currently logged-in user. It takes a name
// and a URL as arguments. If the user is not logged in, an error is returned.
// If the feed is successfully added, the function prints the new feed details.
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	// Get the current user from the config
	currentUser := s.cfg.CurrentUserName
	if currentUser == "" {
		return fmt.Errorf("no user logged in")
	}

	// Get the user from the database
	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Create a new feed
	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	// Print the new feed details
	fmt.Print("Feed created:\n")
	fmt.Printf("ID: %s\n", newFeed.ID)
	fmt.Printf("Name: %s\n", newFeed.Name)
	fmt.Printf("URL: %s\n", newFeed.Url)
	fmt.Printf("User ID: %s\n", newFeed.UserID)

	return nil
}
