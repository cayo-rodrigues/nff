-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices_printings ADD COLUMN file_name TEXT DEFAULT '';
ALTER TABLE invoices_printings RENAME COLUMN custom_file_name TO custom_file_name_prefix;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices_printings RENAME COLUMN custom_file_name_prefix TO custom_file_name;
ALTER TABLE invoices_printings DROP COLUMN file_name;
-- +goose StatementEnd
