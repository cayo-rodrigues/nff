-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN is_interstate TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices DROP COLUMN is_interstate;
-- +goose StatementEnd
