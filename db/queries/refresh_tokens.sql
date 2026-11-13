-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token_hash,
    expires_at
) VALUES (
    ?, ?, ?
)
RETURNING *;


-- name: GetActiveRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token_hash = ?
  AND expires_at > unixepoch()
  AND revoked_at IS NULL
LIMIT 1;


-- name: RevokeRefreshToken :execrows
UPDATE refresh_tokens
SET revoked_at = unixepoch()
WHERE token_hash = ?
  AND revoked_at IS NULL;


-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked_at = unixepoch()
WHERE user_id = ?
  AND revoked_at IS NULL;


-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at <= unixepoch()
   OR revoked_at IS NOT NULL;

-- name: RotateRefreshToken :execrows
UPDATE refresh_tokens
SET revoked_at = unixepoch()
WHERE token_hash = ?
  AND expires_at > unixepoch()
  AND revoked_at IS NULL;
