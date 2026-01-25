-- +goose Up
-- +goose StatementBegin
ALTER TYPE chain_t ADD VALUE 'Kalamaja';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
