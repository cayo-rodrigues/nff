-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics_results (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(5), -- total, month, record (size increases to 6 in the next migration)
    month_name VARCHAR(9),
    total_income VARCHAR(16) DEFAULT '0,00',
    total_expenses VARCHAR(16) DEFAULT '0,00',
    avg_income VARCHAR(16) DEFAULT '0,00',
    avg_expenses VARCHAR(16) DEFAULT '0,00',
    diff VARCHAR(16) DEFAULT '0,00',
    is_positive BOOLEAN DEFAULT false,
    total_records INT DEFAULT 0,
    positive_records INT DEFAULT 0,
    negative_records INT DEFAULT 0,
    metrics_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_results_metrics FOREIGN KEY (metrics_id) REFERENCES metrics_history(id) ON DELETE CASCADE,
    CONSTRAINT fk_results_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics_results;
-- +goose StatementEnd
