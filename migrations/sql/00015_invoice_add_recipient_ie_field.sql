-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN recipient_ie TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices DROP COLUMN recipient_ie;
-- +goose StatementEnd
