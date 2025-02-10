-- name: MarkFeedFetched :exec
UPDATE feeds SET updated_at = NOW(), last_fetched_at = NOW()
WHERE id = $1;
