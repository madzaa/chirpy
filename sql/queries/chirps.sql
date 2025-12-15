-- name: CreateChirps :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(),
        now(),
        now(),
        $1,
        $2)
RETURNING *;

-- name: GetChirps :many
SELECT *
FROM chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteUserChirpById :one
DELETE
FROM chirps
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: DeleteChirp :exec
delete
from chirps;

-- name: DeleteChirpById :exec
delete
from chirps
where id = $1;

-- name: GetChirpsByUser :many
SELECT *
FROM chirps
WHERE user_id = $1;