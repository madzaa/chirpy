-- +goose Up
CREATE TABLE chirps
(
    id         uuid primary key,
    created_at timestamp                                    not null,
    updated_at timestamp                                    not null,
    body       varchar                                      not null,
    user_id    uuid references users (id) on delete cascade not null
);

-- +goose down
DROP TABLE chirps;
