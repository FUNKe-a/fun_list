-- name: GetMedia :one
SELECT * FROM media
WHERE media_id = ? LIMIT 1;

-- name: ListMedia :many
SELECT * FROM media;

-- name: CreateMedia :one
INSERT INTO media (
	title, director_id,
	released_at, synopsis
) VALUES (
	?, ?,
	?, ?
)
RETURNING *;

-- name: DeleteMedia :exec
DELETE FROM media
WHERE media_id = ?;

-- name: UpdateMedia :one
UPDATE media
set title = ?,
director_id = ?,
released_at = ?,
synopsis = ?
WHERE media_id = ?
RETURNING *;
