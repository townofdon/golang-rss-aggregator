// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2
`

type DeleteFeedFollowParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollow, arg.ID, arg.UserID)
	return err
}

const getFeedFollowsByUserId = `-- name: GetFeedFollowsByUserId :many
SELECT
  f.id as feed_id,
  f.created_at,
  f.updated_at,
  f.name,
  f.url,
  f.user_id,
  ff.id AS feed_follow_id
FROM feed_follows ff
JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1
`

type GetFeedFollowsByUserIdRow struct {
	FeedID       uuid.UUID `json:"feed_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Url          string    `json:"url"`
	UserID       uuid.UUID `json:"user_id"`
	FeedFollowID uuid.UUID `json:"feed_follow_id"`
}

func (q *Queries) GetFeedFollowsByUserId(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsByUserIdRow
	for rows.Next() {
		var i GetFeedFollowsByUserIdRow
		if err := rows.Scan(
			&i.FeedID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.FeedFollowID,
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
