package main

import (
	"context"
	"fmt"
)

// handlerFeeds prints all feeds with their associated user names to the console.
// It takes no arguments, and returns an error if any arguments are passed.
func handlerFeeds(s *state, cmd command) error {
	// Ensure no arguments are passed
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: feeds (no arguments allowed)")
	}

	// Fetch all feeds with their associated user names
	feeds, err := s.db.GetAllFeedsWithUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch feeds: %w", err)
	}

	// Print the feeds to the console
	fmt.Print("Feeds:\n")
	for _, feed := range feeds {
		fmt.Printf("Feed Name: %s\nFeed URL: %s\nUser Name: %s\n\n", feed.FeedName, feed.FeedUrl, feed.UserName)
	}

	return nil
}
