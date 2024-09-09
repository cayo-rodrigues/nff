-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN sender_ie TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices DROP COLUMN sender_ie;
-- +goose StatementEnd
