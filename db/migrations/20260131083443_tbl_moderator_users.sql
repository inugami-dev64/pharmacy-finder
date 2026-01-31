-- +goose Up
-- +goose StatementBegin
CREATE TABLE moderator_users (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    "password" VARCHAR(60) NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    administrator BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_moderator_users_username ON moderator_users (username);
CREATE INDEX idx_moderator_users_first_name_last_name ON moderator_users (first_name, last_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_moderator_users_username;
DROP INDEX idx_moderator_users_first_name_last_name;
-- +goose StatementEnd
