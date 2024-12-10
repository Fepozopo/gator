package main

import (
	"context"
	"fmt"
)

// handlerAgg fetches the RSS feed from the Wagslane blog and prints it to the
// console.
func handlerAgg(s *state, cmd command) error {
	const feedURL = "https://www.wagslane.dev/index.xml"

	// Fetch the RSS feed
	ctx := context.Background()
	feed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	// Print the RSS feed to the console
	fmt.Printf("Title: %s\n", feed.Title)
	fmt.Printf("Description: %s\n", feed.Description)
	fmt.Printf("Link: %s\n", feed.Link)
	fmt.Printf("\nItems:\n")
	for _, item := range feed.Items {
		fmt.Printf("----- %s -----\n\n* %s\n\n* %s\n\n\n\n", item.Title, item.Description, item.Link)
	}

	return nil
}
