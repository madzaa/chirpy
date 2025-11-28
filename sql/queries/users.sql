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

-- name: GetUserById :one
SELECT *
FROM users
where id = $1;

-- name: DeleteUsers :exec
DELETE
FROM users;

-- name: UpdateUsers :exec
UPDATE users
SET email         = $1,
    hash_password = $2,
    updated_at    = now()
WHERE id = $3;