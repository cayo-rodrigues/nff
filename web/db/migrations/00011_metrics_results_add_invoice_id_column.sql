-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_results ADD COLUMN invoice_id VARCHAR(11) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_results DROP COLUMN invoice_id;
-- +goose StatementEnd
