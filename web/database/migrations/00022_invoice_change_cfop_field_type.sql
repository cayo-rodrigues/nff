-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN cfop TYPE VARCHAR(197);
ALTER TABLE invoices ALTER COLUMN cfop SET DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN cfop DROP DEFAULT;
ALTER TABLE invoices ALTER COLUMN cfop TYPE INT;
-- +goose StatementEnd
