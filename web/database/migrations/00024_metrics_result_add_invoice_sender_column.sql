-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_results ADD COLUMN invoice_sender VARCHAR(128) NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_results DROP COLUMN invoice_sender;
-- +goose StatementEnd
