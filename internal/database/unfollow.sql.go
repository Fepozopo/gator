// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: unfollow.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const deleteFeedFollowByUserAndURL = `-- name: DeleteFeedFollowByUserAndURL :exec
DELETE FROM feed_follows
USING feeds
WHERE feed_follows.feed_id = feeds.id
    AND feed_follows.user_id = $1
    AND feeds.url = $2
`

type DeleteFeedFollowByUserAndURLParams struct {
	UserID uuid.UUID
	Url    string
}

func (q *Queries) DeleteFeedFollowByUserAndURL(ctx context.Context, arg DeleteFeedFollowByUserAndURLParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowByUserAndURL, arg.UserID, arg.Url)
	return err
}