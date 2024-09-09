-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT, -- total, month, record
    month_name TEXT,
    total_income TEXT DEFAULT '0,00',
    total_expenses TEXT DEFAULT '0,00',
    avg_income TEXT DEFAULT '0,00',
    avg_expenses TEXT DEFAULT '0,00',
    diff TEXT DEFAULT '0,00',
    is_positive BOOLEAN DEFAULT false,
    total_records INT DEFAULT 0,
    positive_records INT DEFAULT 0,
    negative_records INT DEFAULT 0,
    metrics_id INT NOT NULL,
    created_by INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    issue_date DATETIME,
    invoice_id TEXT DEFAULT '',   
    entity_id INT,
    CONSTRAINT fk_results_metrics FOREIGN KEY (metrics_id) REFERENCES metrics_history(id) ON DELETE CASCADE,
    CONSTRAINT fk_results_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_metrics_results_entity FOREIGN KEY (entity_id) REFERENCES entities(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics_results;
-- +goose StatementEnd
