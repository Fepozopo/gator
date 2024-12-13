package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Fepozopo/gator/internal/database"
	"github.com/google/uuid"
)

// handlerAddFeed creates a new feed in the database and automatically follows it
// for the current user. It takes two arguments: the name of the feed, and the
// URL of the feed. The function returns an error if the feed cannot be created.
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	// Create a new feed
	now := time.Now()

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
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
		UserID:    user.ID,
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
		fmt.Printf("Feed created and followed by %s:\n", user.Name)
		fmt.Printf("ID: %s\n", newFeed.ID)
		fmt.Printf("Name: %s\n", newFeed.Name)
		fmt.Printf("URL: %s\n", newFeed.Url)
		fmt.Printf("User ID: %s\n", newFeed.UserID)
	}

	return nil
}
