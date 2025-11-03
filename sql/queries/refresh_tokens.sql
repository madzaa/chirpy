-- name: CreateChirps :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($,
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

-- name: DeleteChirps :exec
delete
from chirps;