-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_results ADD COLUMN printing_id INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_results DROP COLUMN printing_id;
-- +goose StatementEnd
