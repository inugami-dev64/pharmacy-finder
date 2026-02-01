-- +goose Up
-- +goose StatementBegin
CREATE TYPE comment_review_result_t AS ENUM ('APPROVED', 'PERSONAL_ATTACK', 'OFFENSIVE', 'OTHER', 'NONE');
CREATE TABLE comment_reviews (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    comment_id BIGINT NOT NULL REFERENCES pharmacy_reviews(id),
    moderator_id UUID NOT NULL REFERENCES moderator_users(id),
    result comment_review_result_t NOT NULL,
    mod_comment VARCHAR(1024) NULL,
    marked_for_deletion BOOLEAN NOT NULL DEFAULT FALSE,
    reviewed_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT comment_reviews_approval CHECK (result = 'APPROVED' AND NOT marked_for_deletion OR result != 'APPROVED'),
    CONSTRAINT comment_reviews_result_not_none CHECK (result != 'NONE')
);

CREATE INDEX idx_comment_reviews_marked_for_deletion_reviewed_at ON comment_reviews(marked_for_deletion, reviewed_at);
CREATE UNIQUE INDEX idx_comment_reviews_comment_id_moderator_id ON comment_reviews(comment_id, moderator_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_comment_reviews_comment_id_moderator_id;
DROP INDEX idx_comment_reviews_marked_for_deletion_reviewed_at;
DROP TABLE comment_reviews;
DROP TYPE comment_review_result_t;
-- +goose StatementEnd
