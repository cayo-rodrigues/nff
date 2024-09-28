-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_history DROP COLUMN total_income;
ALTER TABLE metrics_history DROP COLUMN total_expenses;
ALTER TABLE metrics_history DROP COLUMN avg_income;
ALTER TABLE metrics_history DROP COLUMN avg_expenses;
ALTER TABLE metrics_history DROP COLUMN diff;
ALTER TABLE metrics_history DROP COLUMN is_positive;
ALTER TABLE metrics_history DROP COLUMN total_records;
ALTER TABLE metrics_history DROP COLUMN positive_records;
ALTER TABLE metrics_history DROP COLUMN negative_records;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_history ADD COLUMN total_income TEXT DEFAULT '0,00';
ALTER TABLE metrics_history ADD COLUMN total_expenses TEXT DEFAULT '0,00';
ALTER TABLE metrics_history ADD COLUMN avg_income TEXT DEFAULT '0,00';
ALTER TABLE metrics_history ADD COLUMN avg_expenses TEXT DEFAULT '0,00';
ALTER TABLE metrics_history ADD COLUMN diff TEXT DEFAULT '0,00';
ALTER TABLE metrics_history ADD COLUMN is_positive BOOLEAN DEFAULT false;
ALTER TABLE metrics_history ADD COLUMN total_records INT DEFAULT 0;
ALTER TABLE metrics_history ADD COLUMN positive_records INT DEFAULT 0;
ALTER TABLE metrics_history ADD COLUMN negative_records INT DEFAULT 0;
-- +goose StatementEnd
