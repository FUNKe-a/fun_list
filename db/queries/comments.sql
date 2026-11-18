-- name: GetComment :one
SELECT * FROM comments
WHERE comment_id = ? LIMIT 1;

-- name: ListComments :many
SELECT * FROM comments;

-- name: ListMediaComments :many
SELECT *
FROM comments
WHERE media_id = sqlc.arg('media_id')
  AND comment_id >= sqlc.arg('start_key')
ORDER BY comment_id DESC
LIMIT sqlc.arg('page_size');

-- name: CreateComment :one
INSERT INTO comments (
	user_id, media_id,
	content, rating
) VALUES (
	?, ?,
	?, ?
)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE comment_id = ?;

-- name: UpdateComment :one
UPDATE comments
set user_id = ?,
media_id = ?,
content = ?,
rating = ?
WHERE comment_id = ?
RETURNING *;
