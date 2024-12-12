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

	// Get the current user
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("current user not found: %w", err)
	}

	// Create a new feed
	now := time.Now()

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUser.ID,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	// Automatically follow the created feed
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		// Print the new feed details
		fmt.Print("Feed created:\n")
		fmt.Printf("ID: %s\n", newFeed.ID)
		fmt.Printf("Name: %s\n", newFeed.Name)
		fmt.Printf("URL: %s\n", newFeed.Url)
		fmt.Printf("User ID: %s\n\n", newFeed.UserID)
		fmt.Printf("Feed created, but failed to auto-follow: %v\n", err)
	} else {
		// Print the new feed details
		fmt.Printf("Feed created and followed by %s:\n", currentUser.Name)
		fmt.Printf("ID: %s\n", newFeed.ID)
		fmt.Printf("Name: %s\n", newFeed.Name)
		fmt.Printf("URL: %s\n", newFeed.Url)
		fmt.Printf("User ID: %s\n", newFeed.UserID)
	}

	return nil
}
