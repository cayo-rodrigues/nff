-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ADD COLUMN file_name VARCHAR(128) DEFAULT '';
ALTER TABLE invoices RENAME COLUMN custom_file_name TO custom_file_name_prefix;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices RENAME COLUMN custom_file_name_prefix TO custom_file_name;
ALTER TABLE invoices DROP COLUMN file_name;
-- +goose StatementEnd
