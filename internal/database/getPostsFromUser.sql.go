// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: getPostsFromUser.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getPostsFromUser = `-- name: GetPostsFromUser :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id
FROM posts p
INNER JOIN feed_follows ff
  ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
LIMIT $2
`

type GetPostsFromUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsFromUser(ctx context.Context, arg GetPostsFromUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsFromUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
