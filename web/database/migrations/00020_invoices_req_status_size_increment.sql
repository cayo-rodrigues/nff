-- +goose Up
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN req_status TYPE VARCHAR(8);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE invoices ALTER COLUMN req_status TYPE VARCHAR(7);
-- +goose StatementEnd
