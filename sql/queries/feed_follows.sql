-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowsByUserId :many
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
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE id = $1 AND user_id = $2;
