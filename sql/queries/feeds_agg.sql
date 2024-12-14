-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
WHERE last_fetched_at IS NULL
   OR last_fetched_at = (
       SELECT MIN(last_fetched_at)
       FROM feeds
   )
LIMIT 1;
