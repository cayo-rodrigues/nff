-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_results ADD COLUMN issue_date TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_results DROP COLUMN issue_date;
-- +goose StatementEnd
