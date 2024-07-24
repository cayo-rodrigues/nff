-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices_printings ADD COLUMN custom_file_name VARCHAR(64) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices_printings DROP COLUMN custom_file_name;
-- +goose StatementEnd
