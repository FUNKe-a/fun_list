-- name: GetDirector :one
SELECT * FROM directors
WHERE director_id = ? LIMIT 1;

-- name: ListDirectors :many
SELECT * FROM directors;

-- name: CreateDirector :one
INSERT INTO directors (
	name, date_of_birth,
	date_of_death, biography
) VALUES (
	?, ?,
	?, ?
)
RETURNING *;

-- name: DeleteDirector :exec
DELETE FROM directors
WHERE director_id = ?;

-- name: UpdateDirector :one
UPDATE directors
set name = ?,
date_of_birth = ?,
date_of_death = ?,
biography = ?
WHERE director_id = ?
RETURNING *;
