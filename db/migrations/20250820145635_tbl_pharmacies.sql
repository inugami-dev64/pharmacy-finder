-- +goose Up
-- +goose StatementBegin
CREATE TYPE chain_t AS ENUM ('Apotheka', 'Südameapteek', 'Benu', 'Euroapteek');
CREATE TABLE pharmacies (
    id BIGSERIAL PRIMARY KEY,
    chain chain_t NOT NULL,
    "name" VARCHAR(256) NOT NULL,
    "address" VARCHAR(64) NOT NULL,
    city VARCHAR(32) NOT NULL,
    county VARCHAR(32) NOT NULL,
    postal_code VARCHAR(6) NOT NULL,
    phone_number VARCHAR(12) NOT NULL, -- We only support Estonian phone numbers for now
    latitude REAL NOT NULL,
    longitude REAL NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pharmacies;
DROP TYPE chain_t;
-- +goose StatementEnd
