-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(128) UNIQUE NOT NULL,
    password VARCHAR(256) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS entities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    user_type VARCHAR(64),
    cpf_cnpj VARCHAR(14),
    ie VARCHAR(13),
    email VARCHAR(128),
    password TEXT,
    postal_code VARCHAR(8),
    neighborhood VARCHAR(64),
    street_type VARCHAR(64),
    street_name VARCHAR(64),
    number VARCHAR(6),
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_entity_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(9),
    protocol VARCHAR(13),
    operation VARCHAR(7) NOT NULL,
    cfop INT NOT NULL,
    is_final_customer VARCHAR(4) NOT NULL,
    is_icms_contributor VARCHAR(6) NOT NULL,
    shipping DOUBLE PRECISION NOT NULL,
    add_shipping_to_total VARCHAR(4) NOT NULL,
    gta VARCHAR(16),
    invoice_pdf VARCHAR(128) DEFAULT '',
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(512) DEFAULT 'Em andamento...',
    sender_id BIGINT NOT NULL,
    recipient_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES entities(id) ON DELETE CASCADE,
    CONSTRAINT fk_recipient FOREIGN KEY (recipient_id) REFERENCES entities(id) ON DELETE CASCADE,
    CONSTRAINT fk_invoice_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_items (
    id BIGSERIAL PRIMARY KEY,
    item_group VARCHAR(64) NOT NULL,
    description VARCHAR(128) NOT NULL,
    origin VARCHAR(64) NOT NULL,
    unity_of_measurement VARCHAR(8) NOT NULL,
    quantity DOUBLE PRECISION NOT NULL,
    value_per_unity DOUBLE PRECISION NOT NULL,
    invoice_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_invoice FOREIGN KEY(invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    CONSTRAINT fk_item_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_cancelings (
    id BIGSERIAL PRIMARY KEY,
    invoice_number VARCHAR(9) NOT NULL,
    year INT NOT NULL,
    justification VARCHAR(128),
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(512) DEFAULT 'Em andamento...',
    entity_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_canceling_entity FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    CONSTRAINT fk_canceling_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invoices_printings (
    id BIGSERIAL PRIMARY KEY,
    invoice_id VARCHAR(13), -- number or protocol
    invoice_id_type VARCHAR(13),
    invoice_pdf VARCHAR(128) DEFAULT '',
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(512) DEFAULT 'Em andamento...',
    entity_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_printing_entity FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    CONSTRAINT fk_printing_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS metrics_history (
    id BIGSERIAL PRIMARY KEY,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    total_income VARCHAR(16) DEFAULT '0,00',
    total_expenses VARCHAR(16) DEFAULT '0,00',
    avg_income VARCHAR(16) DEFAULT '0,00',
    avg_expenses VARCHAR(16) DEFAULT '0,00',
    diff VARCHAR(16) DEFAULT '0,00',
    is_positive BOOLEAN DEFAULT false,
    total_records INT DEFAULT 0,
    positive_records INT DEFAULT 0,
    negative_records INT DEFAULT 0,
    req_status VARCHAR(7) DEFAULT 'pending', -- success, warning, error, pending
    req_msg VARCHAR(512) DEFAULT 'Em andamento...',
    entity_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_metrics_entity FOREIGN KEY (entity_id) REFERENCES entities(id) ON DELETE CASCADE,
    CONSTRAINT fk_metrics_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics_history;
DROP TABLE invoices_printings;
DROP TABLE invoices_cancelings;
DROP TABLE invoices_items;
DROP TABLE invoices;
DROP TABLE entities;
DROP TABLE users;
-- +goose StatementEnd
