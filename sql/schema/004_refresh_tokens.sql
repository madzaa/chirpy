-- +goose Up
CREATE TABLE refresh_tokens
(
    token      VARCHAR PRIMARY KEY,
    created_at TIMESTAMP                                    NOT NULL,
    updated_at TIMESTAMP                                    NOT NULL,
    user_id    UUID REFERENCES users (ID) ON DELETE CASCADE NOT NULL,
    expires_at TIMESTAMP,
    revoked_at TIMESTAMP DEFAULT NULL

);

-- +goose Down
DROP TABLE refresh_tokens;