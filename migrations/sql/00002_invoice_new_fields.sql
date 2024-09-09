-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN extra_notes TEXT DEFAULT '';
ALTER TABLE invoices ADD COLUMN custom_file_name TEXT DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices DROP COLUMN custom_file_name;
ALTER TABLE invoices DROP COLUMN extra_notes;
-- +goose StatementEnd
