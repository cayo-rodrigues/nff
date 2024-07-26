-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_results ADD COLUMN invoice_pdf VARCHAR(128) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_results DROP COLUMN invoice_pdf;
-- +goose StatementEnd
