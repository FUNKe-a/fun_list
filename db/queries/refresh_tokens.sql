-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE user_id = ? LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
	token_hash, user_id,
	issued_at, expires_at
) VALUES (
	?, ?,
	?, ?
)
RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token_hash = ?;

-- name: RotateToken :one
UPDATE refresh_tokens
SET token_hash = sqlc.arg(new_hash)
WHERE token_hash = sqlc.arg(old_hash)
AND expires_at > sqlc.arg(now)
RETURNING *;
