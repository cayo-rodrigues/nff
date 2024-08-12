-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN recipient_ie VARCHAR(13) NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices DROP COLUMN recipient_ie;
-- +goose StatementEnd
