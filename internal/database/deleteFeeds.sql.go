// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: deleteFeeds.sql

package database

import (
	"context"
)

const deleteFeeds = `-- name: DeleteFeeds :exec
DELETE FROM feeds
`

func (q *Queries) DeleteFeeds(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteFeeds)
	return err
}
