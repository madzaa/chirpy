-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1,
        now(),
        now(),
        $2,
        $3,
        $4)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.*
FROM users u
         INNER JOIN refresh_tokens rt
                    ON u.id = rt.user_id
WHERE rt.token = $1
  AND rt.expires_at > now()
  AND rt.revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $1,
    updated_at = $1
WHERE token = $2;


