package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Fepozopo/gator/internal/database"
)

// handlerBrowse prints a list of the posts for the currently logged-in user.
// If a single argument is provided, it is interpreted as an integer and used
// as a limit for the number of posts to retrieve. If no argument is provided, a
// default limit of 2 is used. The retrieved posts are printed with their title,
// URL, description (if any), and publication date.
func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	// Try to convert the first argument to an integer
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil || parsedLimit <= 0 {
			return fmt.Errorf("invalid limit: %v", cmd.args[0])
		}
		limit = parsedLimit
	}

	// Retrieve a list of posts for a specific user from the database
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}

	// Print a list of the posts
	for _, post := range posts {
		fmt.Printf("\n\n\n========================================\nTitle: %s\n\n* URL: %s\n", post.Title, post.Url)
		if post.Description.Valid {
			fmt.Printf("\n* Description: %s\n", post.Description.String)
		}
		fmt.Printf("* Published at: %s\n========================================", post.PublishedAt.Time)
	}

	return nil
}
