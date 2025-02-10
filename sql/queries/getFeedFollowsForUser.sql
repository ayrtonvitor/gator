-- name: GetFeedFollowsForUser :many
SELECT f.name as feed_name, u.name as user_name
FROM feed_follows ff
INNER JOIN users u
  ON u.id = ff.user_id
INNER JOIN feeds f
  ON f.id = ff.feed_id
WHERE u.Name = $1;
