-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices_printings ADD COLUMN custom_file_name TEXT DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices_printings DROP COLUMN custom_file_name;
-- +goose StatementEnd
