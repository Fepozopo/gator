package main

import (
	"context"
	"fmt"

	"github.com/Fepozopo/gator/internal/database"
)

// handlerFollowing handles the "following" command, which prints all feeds that the current user is following.
func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: following (no arguments allowed)")
	}

	// Fetch feed follows for the user
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch following: %w", err)
	}

	// Print the feed names
	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}
