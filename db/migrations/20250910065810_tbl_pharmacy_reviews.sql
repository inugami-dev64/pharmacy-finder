-- +goose Up
-- +goose StatementBegin
CREATE TYPE prescription_t AS ENUM ('Imago', 'GenderGP', 'National');
CREATE TABLE pharmacy_reviews (
    id BIGSERIAL PRIMARY KEY,
    pharmacy_id BIGINT NOT NULL REFERENCES pharmacies(id),
    prescription_type prescription_t NOT NULL,
    stars INT NOT NULL,
    hrt_kind BPCHAR(1) NOT NULL CONSTRAINT hrt_kind_check CHECK (hrt_kind = 'e' OR hrt_kind = 't'),
    nationality BPCHAR(2),
    review TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    modification_code VARCHAR(64) NOT NULL
);

CREATE INDEX idx_pharmacy_review_pharmacy_id ON pharmacy_reviews (pharmacy_id);
CREATE INDEX idx_pharmacy_review_updated_at ON pharmacy_reviews (updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_pharmacy_review_updated_at;
DROP INDEX idx_pharmacy_review_pharmacy_id;
DROP TABLE pharmacy_reviews;
DROP TYPE prescription_t;
-- +goose StatementEnd
