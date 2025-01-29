-- name: CreatePost :one
INSERT INTO posts (
  id,
  created_at,
  updated_at,
  title,
  description,
  published_at,
  url,
  feed_id
)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
	p.id AS post_id,
	p.published_at,
	p.title,
	p.description,
	p.url,
	f.name AS feed_name,
	f.url AS feed_url
FROM posts p
JOIN feed_follows ff ON ff.feed_id = p.feed_id
JOIN feeds f ON f.id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2 OFFSET $3;
