-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_results ALTER COLUMN type TYPE VARCHAR(6);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_results ALTER COLUMN type TYPE VARCHAR(5);
-- +goose StatementEnd
