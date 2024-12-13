package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Fepozopo/gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing feed URL")
	}
	feedURL := cmd.args[0]

	// Delete the feed follow record
	err := s.db.DeleteFeedFollowByUserAndURL(context.Background(), database.DeleteFeedFollowByUserAndURLParams{
		UserID: user.ID,
		Url:    feedURL,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return fmt.Errorf("no feed follow record found for URL: %s", feedURL)
		}
		return fmt.Errorf("failed to unfollow feed: %v", err)
	}

	fmt.Printf("Successfully unfollowed feed: %s\n", feedURL)
	return nil
}
