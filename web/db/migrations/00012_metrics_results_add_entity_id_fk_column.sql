-- +goose Up
-- +goose StatementBegin
ALTER TABLE metrics_results ADD COLUMN entity_id BIGINT DEFAULT 1;

ALTER TABLE metrics_results
ADD CONSTRAINT fk_metrics_results_entity
FOREIGN KEY (entity_id) REFERENCES entities(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metrics_results DROP CONSTRAINT fk_metrics_results_entity;
ALTER TABLE metrics_results DROP COLUMN entity_id;
-- +goose StatementEnd
