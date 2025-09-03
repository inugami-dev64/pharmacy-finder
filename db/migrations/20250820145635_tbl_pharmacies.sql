-- +goose Up
-- +goose StatementBegin
CREATE TYPE chain_t AS ENUM ('Apotheka', 'Südameapteek', 'Benu', 'Euroapteek');
CREATE TABLE pharmacies (
    id BIGSERIAL PRIMARY KEY,
    pharmacy_id BIGINT NOT NULL, -- ID of the pharmacy as scraped
    chain chain_t NOT NULL,
    "name" VARCHAR(256) NOT NULL,
    "address" VARCHAR(64) NOT NULL,
    city VARCHAR(32) NOT NULL,
    county VARCHAR(32) NOT NULL,
    postal_code VARCHAR(6) NOT NULL,
    email VARCHAR(32) NOT NULL,
    phone_number VARCHAR(12) NOT NULL, -- We only support Estonian phone numbers for now
    mod_time TIMESTAMP NOT NULL DEFAULT now(),
    latitude REAL NOT NULL,
    longitude REAL NOT NULL
);

CREATE INDEX idx_pharmacies_latitude_longitude ON pharmacies (latitude, longitude);
CREATE UNIQUE INDEX idx_pharmacies_pharmacy_id_chain ON pharmacies (pharmacy_id, chain);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_pharmacies_latitude_longitude;
DROP INDEX idx_pharmacies_pharmacy_id_chain;
DROP TABLE pharmacies;
DROP TYPE chain_t;
-- +goose StatementEnd
