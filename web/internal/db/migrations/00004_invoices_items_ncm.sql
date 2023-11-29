-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices_items ADD COLUMN ncm VARCHAR(16) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices_items DROP COLUMN ncm;
-- +goose StatementEnd
