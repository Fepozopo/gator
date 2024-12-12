package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Fepozopo/gator/internal/database"
	"github.com/google/uuid"
)

// handlerFollow handles the "follow" command, which creates a new feed follow
// record between the current user and the given feed. The feed is looked up by
// URL, and an error is returned if the feed does not exist. The current user is
// looked up by name, and an error is returned if the current user does not
// exist. The feed follow record is then created, and a success message is
// printed with the user and feed details.
func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: follow <feed_url>")
	}
	feedURL := cmd.args[0]

	// Look up the feed by URL
	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	// Get the current user
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("current user not found: %w", err)
	}

	// Create a feed follow record
	now := time.Now()

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}

	// Print success message
	fmt.Printf("User %s is now following feed: %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}
