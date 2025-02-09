-- name: GetFeeds :many
SELECT f.*, u.name AS user_name
FROM feeds f
LEFT JOIN users u
  ON f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;
