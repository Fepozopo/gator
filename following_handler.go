package main

import (
	"context"
	"fmt"
)

// handlerFollowing handles the "following" command, which prints all feeds that the current user is following.
func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: following (no arguments allowed)")
	}

	// Get the current user
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("current user not found: %w", err)
	}

	// Fetch feed follows for the user
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch following: %w", err)
	}

	// Print the feed names
	fmt.Printf("Feeds followed by %s:\n", currentUser.Name)
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}
