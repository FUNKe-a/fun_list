-- name: GetUser :one
SELECT * FROM app_users
WHERE user_id = ? LIMIT 1;
