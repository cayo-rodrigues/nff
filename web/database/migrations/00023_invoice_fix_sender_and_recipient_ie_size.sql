-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN sender_ie TYPE VARCHAR(16);
ALTER TABLE invoices ALTER COLUMN recipient_ie TYPE VARCHAR(16);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN sender_ie TYPE VARCHAR(13);
ALTER TABLE invoices ALTER COLUMN recipient_ie TYPE VARCHAR(13);
-- +goose StatementEnd
