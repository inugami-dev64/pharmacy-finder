-- +goose Up
-- +goose StatementBegin
CREATE TABLE moderator_users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    email VARCHAR(64) NOT NULL,
    "password" VARCHAR(60) NOT NULL,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    registration_timestamp TIMESTAMP DEFAULT now() NOT NULL,
    last_login_timestamp TIMESTAMP DEFAULT now() NOT NULL,
    administrator BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE UNIQUE INDEX idx_moderator_users_username ON moderator_users (username);
CREATE UNIQUE INDEX idx_moderator_users_email ON moderator_users (email);
CREATE INDEX idx_moderator_users_first_name_last_name ON moderator_users (first_name, last_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_moderator_users_first_name_last_name;
DROP INDEX idx_moderator_users_email;
DROP INDEX idx_moderator_users_username;
DROP TABLE moderator_users;
-- +goose StatementEnd
