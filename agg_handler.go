package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Fepozopo/gator/internal/database"
	"github.com/google/uuid"
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

// scrapeFeeds runs the RSS feed aggregation process, which fetches feeds from
// the database, marks them as fetched, fetches the feed content, and saves the
// feed items to the database as posts. If an error occurs during the process,
// it is propagated up the call stack.
//
// The function returns an error if any database query fails, or if the feed
// content cannot be fetched.
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
			Valid: true,
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

	// Iterate over each item in an RSS feed and parse the published date of each item
	for _, item := range rssFeed.Items {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("Failed to parse published date for %s: %v\n", item.Title, err)
			publishedAt = time.Time{} // Default to zero value
		}

		// Convert item.Description to an sql.NullString
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		// Create a new post in the database
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: sql.NullTime{Time: publishedAt, Valid: !publishedAt.IsZero()},
			FeedID:      feed.ID,
		})
		// Check if the creation is successful, a duplicate, or any other error
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				fmt.Printf("\nPost with URL %s already exists. Skipping.\n", item.Link)
				continue
			}
			return fmt.Errorf("failed to save post: %w", err)
		}
	}

	return nil
}
