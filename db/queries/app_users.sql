-- name: GetUsers :one
SELECT * FROM app_users
WHERE user_id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM app_users
WHERE username = ? LIMIT 1;

-- name: GetUserByUserID :one
SELECT * FROM app_users
WHERE user_id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM app_users;

-- name: CreateUser :one
INSERT INTO app_users (
	username, email,
	password_hash
) VALUES (
	?, ?,
	?
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app_users
WHERE user_id = ?;

-- name: UpdateUser :one
UPDATE app_users
set username = ?,
email = ?,
password_hash = ?,
role = ?
WHERE user_id = ?
RETURNING *;
