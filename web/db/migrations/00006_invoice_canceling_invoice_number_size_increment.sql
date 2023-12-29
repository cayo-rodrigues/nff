-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices_cancelings ALTER COLUMN invoice_number TYPE VARCHAR(11);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices_cancelings ALTER COLUMN invoice_number TYPE VARCHAR(9);
-- +goose StatementEnd
