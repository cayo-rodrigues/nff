-- +goose Up
-- +goose StatementBegin
ALTER TABLE entities ALTER COLUMN ie TYPE VARCHAR(16);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE entities ALTER COLUMN ie TYPE VARCHAR(13);
-- +goose StatementEnd

