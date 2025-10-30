-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hash_password)
VALUES (gen_random_uuid(),
        now(),
        now(),
        $1, $2)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
where email = $1;

-- name: DeleteUsers :exec
DELETE
FROM users;