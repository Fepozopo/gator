package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Fepozopo/gator/internal/database"
)

// handlerAgg handles the "agg" command, which starts an RSS feed aggregation
// process that runs indefinitely. It requires a time duration argument
// specifying the interval between each feed collection cycle.
//
// The function initializes a ticker with the given interval, and repeatedly
// calls the scrapeFeeds function to fetch and process feeds. If an error
// occurs during scraping, it logs the error to the console.
//
// Args:
//
//	s: The application state, containing database queries and configuration.
//	cmd: The command input, which should include the time interval between
//	     requests as an argument.
//
// Returns:
//
//	An error if the time_between_reqs argument is missing or invalid.
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("time_between_reqs argument is required")
	}

	// Parse the argument into a time.Duration value
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid time duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	// Run the scraper in a loop
	for {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("\nError scraping feeds: %v\n", err)
		}
		<-ticker.C
	}
}

// scrapeFeeds retrieves the next feed to be fetched from the database, marks it as fetched,
// fetches the RSS feed from the feed's URL, and prints the feed's items to the console.
// If there are no feeds to fetch, it returns without error. If any step fails, it returns an error.
func scrapeFeeds(s *state) error {
	// Get the next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No feeds to fetch.\n")
			return nil
		}
		return fmt.Errorf("failed to get next feed: %w", err)
	}

	// Mark the feed as fetched
	now := time.Now()
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  now,
			Valid: true, // Set to true because we're providing a value
		},
		UpdatedAt: now,
		ID:        feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	// Fetch the feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed from %s: %w", feed.Url, err)
	}

	// Print the feed's items
	fmt.Printf("\nFetched feed: %s (%s)\n", feed.Name, feed.Url)
	for _, item := range rssFeed.Items {
		fmt.Printf("* %s\n", item.Title)
	}

	return nil
}
