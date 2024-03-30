-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN gta TYPE VARCHAR(32);
ALTER TABLE invoices ALTER COLUMN gta SET DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN gta DROP DEFAULT;
ALTER TABLE invoices ALTER COLUMN gta TYPE VARCHAR(16);
-- +goose StatementEnd
