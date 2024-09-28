-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices_items ADD COLUMN ncm TEXT DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices_items DROP COLUMN ncm;
-- +goose StatementEnd
