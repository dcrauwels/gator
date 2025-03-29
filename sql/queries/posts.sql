-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
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
ON CONFLICT (url)
DO UPDATE SET updated_at = $3, title = $4, description = $6
RETURNING *;

-- name: GetPostsByUserName :many
SELECT *
FROM posts
WHERE feed_id IN (
    SELECT feed_follows.feed_id
    FROM feed_follows
    INNER JOIN users ON feed_follows.user_id = users.id
    WHERE users.name = $1
)
ORDER BY published_at DESC
LIMIT $2;

-- name: ResetPosts :exec
DELETE FROM posts;

