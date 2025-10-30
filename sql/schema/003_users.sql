-- +goose Up
ALTER table users
    ADD hash_password varchar;

-- +goose Down
ALTER table users
    DROP hash_password;